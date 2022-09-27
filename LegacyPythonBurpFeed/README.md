# BurpFeed
**-- This is the legacy implementaion and is no longer supported therefore it is recommended that you use the go version --**
![](https://github.com/ZephrFish/BurpFeed/blob/master/LogoBurpFeed.png)

A tool for passing and adding a list of URLs to Burp's sitemap/target tab, really useful for populating the targets tab with a big list of URLs. Originally an idea that [@InfoSecPS](https://twitter.com/InfoSecPS) and I threw together, then I tweaked and hacked together this chaos! 

The tool is written in both Python and Go, GoBurpFeed was written by [Mantis](https://github.com/MantisSTS) while the python version has been a collaboration between ZephrFish, InfosecPS and Mantis.


## BurpFeed Setup
To set this up, you'll need the following:
- Burp Suite
  - FLOW Burp Extension (https://github.com/hvqzao/burp-flow), also available on BApps store :-)
- Python3 & Pip
  - `pip install -r requirements.txt`

Chuck your target URLs or IPs in a file, can be named whatever but must have http/https prefixed at the start of line for this to work. Additionally you'll want to edit line 15 (example below), depending on what the IP of your burp proxy is. Either done via localhost or if in a Virtual Machine feed the listening address of burp (you'll need to flip the proxy interface to listening on all interfaces).

```
proxy = {
                "http": "http://localhost:8080",
                "https": "https://localhost:8080",
        }
```

When you've got all of this setup you can refer to usage.


### Usage:
```
python bfeed.py targets.txt
```

To igrnore warnings you can supress them with this:

```
python -W ignore bfeed.py targets.txt
```

This will probably throw errors but alas it's a hacky tool 😉