#!/env/python
# Burp URL Feeder Non-Threaded
# ZephrFish 2.1 2021
import urllib3
import sys
import re
import requests
import argparse
from argparse import ArgumentParser, FileType

urllib3.exceptions.InsecureRequestWarning

version = "2.1"

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

def burpFeed(urls):
        proxy = {

                "http": "http://127.0.0.1:8080",
                "https": "https://127.0.0.1:8080",
        }

        headers = sys.argv[2]

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
        try:
                burpFeed(sys.argv[1])
        except:
                print("File argument needed! %s <hosts file>" % sys.argv[0])
                sys.exit()
    

if __name__ == "__main__":
	parser = ArgumentParser(prog="befeed.py", description="BurpFeed - Feed URLs or Domains into Burp or a Proxy")
	parser.add_argument("-f", "--urls", action="store", dest="urls", help="domain to search")
	parser.add_argument("-p", "--proxy", action="store", dest="proxy", help="Proxy to pass data to", default="http://127.0.0.1:8080")
	parser.add_argument("-rh", "--headers", action="store_true", dest="bfheaders", help="Add additional headers", default=False)
	parser.add_argument("-v", "--version", action="version", version="BurpFeed v{0}".format(version))
	args = parser.parse_args()

    printBanner()
    burpFeed(urls