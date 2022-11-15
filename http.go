/*
	(c) Yariya
*/

package main

import (
	"fmt"
	"h12.io/socks"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Proxy struct {
	ips                  map[string]struct{}
	targetSites          []string
	httpStatusValidation bool
	timeout              time.Duration
	maxHttpThreads       int64

	openHttpThreads int64
	mu              sync.Mutex
}

var Proxies = &Proxy{
	// in work
	targetSites: []string{"https://google.com", "https://cloudflare.com"},

	httpStatusValidation: false,
	// now cfg file
	timeout:        time.Second * 5,
	maxHttpThreads: int64(config.HttpThreads),
	ips:            make(map[string]struct{}),
}

func (p *Proxy) WorkerThread() {
	for {
		for atomic.LoadInt64(&p.openHttpThreads) < int64(config.HttpThreads) {
			p.mu.Lock()
			for proxy, _ := range p.ips {
				if strings.ToLower(config.ProxyType) == "http" {
					go p.CheckProxyHTTP(proxy)
				} else if strings.ToLower(config.ProxyType) == "socks4" {
					go p.CheckProxySocks4(proxy)
				} else if strings.ToLower(config.ProxyType) == "socks5" {
					go p.CheckProxySocks5(proxy)
				} else {
					log.Fatalln("invalid ProxyType")
				}
				delete(p.ips, proxy)
				break
			}
			p.mu.Unlock()

		}
		time.Sleep(time.Millisecond * 100)
	}
}

func (p *Proxy) CheckProxyHTTP(proxy string) {
	atomic.AddInt64(&p.openHttpThreads, 1)
	defer func() {
		atomic.AddInt64(&p.openHttpThreads, -1)
		atomic.AddUint64(&checked, 1)
	}()

	var err error
	var proxyPort = *port
	s := strings.Split(proxy, ":")
	if len(s) > 1 {
		proxyPort, err = strconv.Atoi(strings.TrimSpace(s[1]))
		if err != nil {
			log.Println(err)
			return
		}
	}

	if len(s) > 1 {
		var err error
		proxyPort, err = strconv.Atoi(s[1])
		if err != nil {
			log.Println(err)
			return
		}
	}
	proxyUrl, err := url.Parse(fmt.Sprintf("http://%s:%d", s[0], proxyPort))
	if err != nil {
		log.Println(err)
		return
	}

	tr := &http.Transport{
		Proxy: http.ProxyURL(proxyUrl),
		DialContext: (&net.Dialer{
			Timeout:   time.Second * time.Duration(config.Timeout.HttpTimeout),
			KeepAlive: time.Second,
			DualStack: true,
		}).DialContext,
	}

	client := http.Client{
		Timeout:   time.Second * time.Duration(config.Timeout.HttpTimeout),
		Transport: tr,
	}

	req, err := http.NewRequest("GET", config.CheckSite, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Add("User-Agent", config.Headers.UserAgent)
	req.Header.Add("accept", config.Headers.Accept)

	res, err := client.Do(req)
	if err != nil {
		atomic.AddUint64(&proxyErr, 1)
		if strings.Contains(err.Error(), "timeout") {
			atomic.AddUint64(&timeoutErr, 1)
			return
		}
		return
	}
	res.Body.Close()
	if res.StatusCode != 200 {
		atomic.AddUint64(&statusCodeErr, 1)
	} else {
		if config.PrintIps.Enabled {
			go PrintProxy(s[0], proxyPort)
		}
		atomic.AddUint64(&success, 1)
		exporter.Add(fmt.Sprintf("%s:%d", s[0], proxyPort))
	}
}

func (p *Proxy) CheckProxySocks4(proxy string) {
	atomic.AddInt64(&p.openHttpThreads, 1)
	defer func() {
		atomic.AddInt64(&p.openHttpThreads, -1)
		atomic.AddUint64(&checked, 1)
	}()

	var err error
	var proxyPort = *port
	s := strings.Split(proxy, ":")
	if len(s) > 1 {
		proxyPort, err = strconv.Atoi(strings.TrimSpace(s[1]))
		if err != nil {
			log.Println(err)
			return
		}
	}

	tr := &http.Transport{
		Dial: socks.Dial(fmt.Sprintf("socks4://%s:%d?timeout=%ds", s[0], proxyPort, config.Timeout.Socks4Timeout)),
	}

	client := http.Client{
		Timeout:   time.Second * time.Duration(config.Timeout.HttpTimeout),
		Transport: tr,
	}

	req, err := http.NewRequest("GET", config.CheckSite, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Add("User-Agent", config.Headers.UserAgent)
	req.Header.Add("accept", config.Headers.Accept)

	res, err := client.Do(req)
	if err != nil {
		atomic.AddUint64(&proxyErr, 1)
		if strings.Contains(err.Error(), "timeout") {
			atomic.AddUint64(&timeoutErr, 1)
			return
		}
		return
	}
	res.Body.Close()
	if res.StatusCode != 200 {
		atomic.AddUint64(&statusCodeErr, 1)
	} else {
		if config.PrintIps.Enabled {
			go PrintProxy(s[0], proxyPort)
		}
		atomic.AddUint64(&success, 1)
		exporter.Add(fmt.Sprintf("%s:%d", s[0], proxyPort))
	}
}

func (p *Proxy) CheckProxySocks5(proxy string) {
	atomic.AddInt64(&p.openHttpThreads, 1)
	defer func() {
		atomic.AddInt64(&p.openHttpThreads, -1)
		atomic.AddUint64(&checked, 1)
	}()

	var err error
	var proxyPort = *port
	s := strings.Split(proxy, ":")
	if len(s) > 1 {
		proxyPort, err = strconv.Atoi(strings.TrimSpace(s[1]))
		if err != nil {
			log.Println(err)
			return
		}
	}

	tr := &http.Transport{
		Dial: socks.Dial(fmt.Sprintf("socks5://%s:%d?timeout=%ds", s[0], proxyPort, config.Timeout.Socks4Timeout)),
	}

	client := http.Client{
		Timeout:   time.Second * time.Duration(config.Timeout.HttpTimeout),
		Transport: tr,
	}

	req, err := http.NewRequest("GET", config.CheckSite, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Add("User-Agent", config.Headers.UserAgent)
	req.Header.Add("accept", config.Headers.Accept)

	res, err := client.Do(req)
	if err != nil {
		atomic.AddUint64(&proxyErr, 1)
		if strings.Contains(err.Error(), "timeout") {
			atomic.AddUint64(&timeoutErr, 1)
			return
		}
		return
	}
	res.Body.Close()
	if res.StatusCode != 200 {
		atomic.AddUint64(&statusCodeErr, 1)
	} else {
		if config.PrintIps.Enabled {
			go PrintProxy(s[0], proxyPort)
		}
		atomic.AddUint64(&success, 1)
		exporter.Add(fmt.Sprintf("%s:%d", s[0], proxyPort))
	}
}
