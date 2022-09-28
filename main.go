package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

var headers []string

type ProgArgs struct {
	proxyURL    string
	debugMode   bool
	threads     int
	HTTPTimeout int
	jitter      int
	sleep       int
	randomAgent bool
	userAgent   string
}

var wg sync.WaitGroup
var args ProgArgs

func getRandomUserAgent() (string, error) {
	// Read user-agents.txt and select a random user agent
	data, err := ioutil.ReadFile("user-agents.txt")
	if err != nil {
		return "", err
	}

	agents := strings.Split(string(data), "\n")
	return strings.Trim(agents[rand.Intn(len(agents))], "\r\n"), nil
}

func makeRequest(u string) {

	proxyURL, err := url.Parse(args.proxyURL)
	if err != nil {
		log.Panic(err)
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			Dial: (&net.Dialer{
				Timeout:   time.Duration(args.HTTPTimeout) * time.Second,
				KeepAlive: time.Duration(args.HTTPTimeout) * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			Proxy:                 http.ProxyURL(proxyURL),
		},
	}

	dest, err := url.Parse(u)

	if err == nil {

		req, err := http.NewRequest("GET", dest.String(), nil)

		if err != nil {
			log.Panic(err)
		}

		req.Header.Set("User-Agent", args.userAgent)

		for _, header := range headers {
			headerSplit := strings.Split(header, ":")
			req.Header.Add(headerSplit[0], headerSplit[1])
		}

		resp, err := client.Do(req)
		if err != nil {
			log.Panic(err)
		}

		if args.debugMode {
			fmt.Println(resp)
		}

		fmt.Println(u)
	}
}

func processJob(c <-chan string, sleep time.Duration, wg *sync.WaitGroup) {
	defer wg.Done()

	for in := range c {
		makeRequest(in)
		time.Sleep(sleep)
	}
}

func main() {

	jobs := make(chan string, 1000)

	var urlsFile string

	// To include headers in requests
	var headersFile string

	flag.StringVar(&urlsFile, "i", "", "./path/to/urls.txt (required)")
	flag.StringVar(&headersFile, "headers", "", "./path/to/headers.txt - Headers file should be in the format 'Header: Value'")
	flag.BoolVar(&args.debugMode, "debug", false, "Turn on debug mode")
	flag.IntVar(&args.threads, "threads", 10, "Number of concurrent jobs to run")
	flag.IntVar(&args.HTTPTimeout, "timeout", 10, "HTTP Timeout time in seconds")
	flag.StringVar(&args.proxyURL, "proxy", "http://127.0.0.1:8080", "The HTTP proxy you want to feed it through")
	flag.IntVar(&args.jitter, "jitter", 5, "A jitter amount to add to the sleep time")
	flag.IntVar(&args.sleep, "sleep", 0, "The number of milliseconds to sleep per request")
	flag.BoolVar(&args.randomAgent, "random-agent", false, "Use a random user-agent string")
	flag.StringVar(&args.userAgent, "user-agent", "", "Use a specified user-agent string")
	flag.Parse()

	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	var sleep time.Duration
	if args.sleep > 0 {

		// Create a random number to sleep by

		randSleep := time.Duration(args.sleep + rand.Intn(args.jitter))
		sleep = time.Duration(time.Millisecond * time.Duration(randSleep))
		fmt.Println("Sleeping for", sleep, "milliseconds per request")
	}

	// Jobs queue
	for j := 0; j < args.threads; j++ {
		wg.Add(1)
		go processJob(jobs, sleep, &wg)
	}

	// Read the headers file
	if headersFile != "" {
		file, err := os.Open(headersFile)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		// Add the lines to the headers array
		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			headers = append(headers, scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}

	// Set the user agent
	if args.userAgent == "" {
		if args.randomAgent {
			var err error
			args.userAgent, err = getRandomUserAgent()
			if err != nil {
				log.Panic(err)
			}
		}
	}

	// Read the URLs
	readFile, err := os.Open(urlsFile)
	if err != nil {
		log.Fatal(err)
	}

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		text := fileScanner.Text()
		if strings.Trim(text, "\r\n\t ") != "" {
			if !strings.HasPrefix(text, "http") {
				jobs <- "https://" + text
				jobs <- "http://" + text
			} else {
				jobs <- text
			}
		}
	}

	close(jobs)

	wg.Wait()
}
