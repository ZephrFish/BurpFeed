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
var wg sync.WaitGroup

func makeRequest(u string) {

	defer wg.Done()

	proxyUrl, err := url.Parse(proxyUrl)
	if err != nil {
		log.Panic(err)
	}

	client := &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   15 * time.Second,
				KeepAlive: 15 * time.Second,
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
	}
}

func processJob(c chan string) {
	for {
		select {
		case url := <-c:
			fmt.Println(url)
			wg.Add(1)
			go func() {
				defer wg.Done()
				makeRequest(url)
			}()
		}
	}
}

func main() {

	c := make(chan string, 1000)
	defer close(c)

	var filename string
	var threads int
	flag.StringVar(&filename, "filename", "", "./path/to/urls.txt")
	flag.BoolVar(&debugMode, "debug", false, "Turn on debug mode")
	flag.IntVar(&threads, "threads", 10, "Number of concurrent jobs to run")
	flag.StringVar(&proxyUrl, "proxy", "http://127.0.0.1:8080", "The HTTP proxy you want to feed it through")
	flag.Parse()

	// Jobs queue
	for j := 0; j < threads; j++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			processJob(c)
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
			c <- fileScanner.Text()
		}
	}

	wg.Wait()
}
