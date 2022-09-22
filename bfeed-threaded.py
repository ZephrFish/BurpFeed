#!/env/python3
# Burp URL Feeder Threaded
# ZephrFish & Mantis 2022
# Python 3 Conversion

import urllib3
import sys
import re
import requests
import argparse
from requests.packages.urllib3.exceptions import InsecureRequestWarning
from multiprocessing import Pool

urllib3.exceptions.InsecureRequestWarning

version = "2.2"

def printBanner():
	print(""""
    __________                   ___________               .___
    \______   \__ _______________\_   _____/___   ____   __| _/
     |    |  _/  |  \_  __ \____ \|    __)/ __ \_/ __ \ / __ | 
     |    |   \  |  /|  | \/  |_> >     \\  ___/\  ___// /_/ | 
     |______  /____/ |__|  |   __/\___  / \___  >\___  >____ | 
            \/             |__|       \/      \/     \/     \/ 
      Version {0}""".format(version)
      sleep(1)

def fetchUrl(url):
        requests.packages.urllib3.disable_warnings(InsecureRequestWarning)
        
        # Change me to whatever you want, can be IP of host or localhost wherever Burp is listening
        proxy = {
                "http": "http://localhost:8080",
                "https": "https://localhost:8080",
        }

        headers = sys.argv[2]

        regex=re.compile('^http://|^https://')
        if re.match(regex, url):
                try:
                        normalresponse = requests.get(url.rstrip(), proxies=proxy, verify=False, timeout=8, headers=headers)
                        print("URL: {0} | Status: {1}".format(url.rstrip(), normalresponse.status_code))
                except: 
                        pass
        else:
                HTTPSecure = "https://"+url.rstrip()
                HTTPNot = "http://"+url.rstrip()
                try:
                        httpsresponse = requests.get(HTTPSecure, proxies=proxy, verify=False, timeout=8, headers=headers)
                        httpresponse = requests.get(HTTPNot, proxies=proxy, verify=False, timeout=8, headers=headers)
                        print("URL: {0} | Status: {1}".format(HTTPNot, httpresponse.status_code))
                        print("URL: {0} | Status: {1}".format(HTTPSecure, httpsresponse.status_code))

                except:
                        pass

def burpFeed(urls, threads):
        pool = Pool(int(threads))
        with open(urls, encoding="utf8") as source_file:
                results = pool.map(fetchUrl, source_file, int(threads))
                print(results)
                                
if __name__ == '__main__':
    try:
        burpFeed(sys.argv[1], sys.argv[2])
    except:
        print("Not enough arguments! %s <hosts file> <headers>" % sys.argv[0])
        sys.exit()
    
