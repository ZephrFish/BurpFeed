# BurpFeed - Go Edition! Still a WIP
![](https://github.com/ZephrFish/BurpFeed/blob/master/LogoBurpFeed.png)
![](https://raw.githubusercontent.com/egonelbre/gophers/ac77b513f41f44a7805694063aaef16ccd95a9b3/vector/projects/network.svg)

A tool for passing and adding a list of URLs to Burp's sitemap/target tab, really useful for populating the targets tab with a big list of URLs. 

## Setup
To set this up, you'll need the following:
- Burp Suite
- Go 
  - Only tested on go1.14.2
  
Chuck your target URLs or IPs in a file, can be named whatever but must have http/https prefixed at the start of line for this to work.


When you've got all of this setup you can refer to usage.


## Usage:
```
go build .
./BurpFeed -h
Usage of BurpFeed:               
-debug         
  Turn on debug mode                            
-filename string
  ./path/to/urls.txt
-proxy string
  The HTTP proxy you want to feed it through (default "http://127.0.0.1:8080")     
-threads int
  Number of concurrent jobs to run (default 10) 
```


## Gophers
The gopher was found in this repository: https://github.com/egonelbre/gophers/
