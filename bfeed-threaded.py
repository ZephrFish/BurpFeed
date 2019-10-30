#!/env/python
# Burp URL Feeder Threaded
# ZephrFish 2019
# Modified by Mantis 2019
# Python2.7 based added in threading and some other jazz
import urllib3
import sys
import re
import requests
import argparse
import urllib3
from multiprocessing import Pool

urllib3.exceptions.InsecureRequestWarning

def fetchUrl(url):
        proxy = {
                "http": "http://localhost:8080",
                "https": "https://localhost:8080",
        }
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

def burpFeed(urls, threads):
        pool = Pool(int(threads))
        with open(urls, encoding="utf8") as source_file:
                results = pool.map(fetchUrl, source_file, int(threads))
                print(results)
                                
if __name__ == '__main__':
    urllib3.disable_warnings()
    burpFeed(sys.argv[1], sys.argv[2])
