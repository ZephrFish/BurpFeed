#!/env/python
# Burp URL Feeder Threaded
# ZephrFish 2019
# Python2.7 based added in threading and some other jazz
# DO NOT USE
import urllib3
import sys
import re
import requests
import argparse
import threading

urllib3.exceptions.InsecureRequestWarning

def fetchUrl(url):
        proxy = {
                "http": "http://localhost:8080",
                "https": "https://localhost:8080",
        }
        regex=re.compile('^http://|^https://')
        if re.match(regex, url):
                normalresponse = requests.get(url.rstrip(), proxies=proxy, verify=False)
                print(url, normalresponse.status_code)
        else:
                HTTPSecure = "https://"+url.rstrip()
                HTTPNot = "http://"+url.rstrip()
                httpsresponse = requests.get(HTTPSecure, proxies=proxy, verify=False)
                httpresponse = requests.get(HTTPNot, proxies=proxy, verify=False)
                print(url.rstrip(), httpsresponse.status_code, httpresponse.status_code)


def burpFeed(urls):
        tCount = 60
        threads = []
        with open(urls.rstrip(), 'r') as f:
                for url in f:
                        while(len(threads) <= tCount):
                                for t in range(tCount):
                                        threading.Thread(target=fetchUrl, args=(url,)).start()
        
                                
                                
if __name__ == '__main__':
    burpFeed(sys.argv[1])
