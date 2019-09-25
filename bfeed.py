#!/env/python
# Burp URL Feeder Non-Threaded
# ZephrFish 2018
import urllib3
import sys
import re
import requests
import argparse

# urllib3.exceptions.InsecureRequestWarning


def burpFeed(urls):
        proxy = {

                "http": "http://127.0.0.1:8080",
                "https": "https://127.0.0.1:8080",
        }

        with open(urls.rstrip(), 'r') as f:
                for url in f:
                        regex=re.compile('^http://|^https://')
                        if re.match(regex, url):
                                try:
                                        normalresponse = requests.get(url.rstrip(), proxies=proxy, verify=False, timeout=8)
                                        print(url, normalresponse.status_code)
                                except: 
                                        pass
                        else:
                                HTTPSecure = "https://"+url.rstrip()
                                HTTPNot = "http://"+url.rstrip()
                                try:
                                        httpsresponse = requests.get(HTTPSecure, proxies=proxy, verify=False, timeout=8)
                                        httpresponse = requests.get(HTTPNot, proxies=proxy, verify=False, timeout=8)
                                        print(url.rstrip(), httpsresponse.status_code, httpresponse.status_code)
                                except:
                                        pass


if __name__ == '__main__':
    burpFeed(sys.argv[1])
