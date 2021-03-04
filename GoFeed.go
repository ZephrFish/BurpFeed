package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

var debugMode bool
var proxyUrl string

type ProgArgs struct {
	proxyUrl    string
	debugMode   bool
	threads     int
	httpTimeout int
}

var wg sync.WaitGroup
var args ProgArgs

func makeRequest(u string) {

	proxyUrl, err := url.Parse(args.proxyUrl)
	if err != nil {
		log.Panic(err)
	}

	client := &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   time.Duration(args.httpTimeout) * time.Second,
				KeepAlive: time.Duration(args.httpTimeout) * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			Proxy:                 http.ProxyURL(proxyUrl),
		},
	}

	dest, err := url.Parse(u)
	if err == nil {
		resp, _ := client.Get(dest.String())

		if debugMode {
			fmt.Println(resp)
		}

		fmt.Println(u)
	}
}

func processJob(c <-chan string) {
	for in := range c {
		makeRequest(in)
	}
}

func main() {

	jobs := make(chan string, 1000)

	var filename string

	flag.StringVar(&filename, "filename", "", "./path/to/urls.txt")
	flag.BoolVar(&args.debugMode, "debug", false, "Turn on debug mode")
	flag.IntVar(&args.threads, "threads", 10, "Number of concurrent jobs to run")
	flag.IntVar(&args.httpTimeout, "timeout", 10, "HTTP Timeout time in seconds")
	flag.StringVar(&args.proxyUrl, "proxy", "http://127.0.0.1:8080", "The HTTP proxy you want to feed it through")
	flag.Parse()

	// Jobs queue
	for j := 0; j < args.threads; j++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			processJob(jobs)
		}()
	}

	readFile, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		text := fileScanner.Text()
		if strings.Trim(text, "\r\n\t ") != "" {
			jobs <- fileScanner.Text()
		}
	}

	close(jobs)

	wg.Wait()
}
