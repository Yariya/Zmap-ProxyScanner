package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type IPAPI struct {
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	Isp         string  `json:"isp"`
	Org         string  `json:"org"`
	As          string  `json:"as"`
	Query       string  `json:"query"`
}

func GetISP(proxy string) (isp *IPAPI) {
	res, err := http.Get("http://ip-api.com/json/" + proxy)
	if err != nil {
		log.Println("couldn't fetch isp")
		return
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("body error isp")
		return
	}
	if err := json.Unmarshal(body, &isp); err != nil {
		log.Println("json error isp")
		return
	}
	res.Body.Close()
	return
}

func PrintProxy(proxy string, port int) {
	if config.PrintIps.DisplayIpInfo {
		ipApi := GetISP(proxy)
		if ipApi == nil {
			fmt.Printf("New Proxy \033[32m%s:%d\033[39m Country: \033[34m error\033[39m ISP: \033[34merror\033[39m\n", proxy, port)
		} else {
			fmt.Printf("New Proxy \033[32m%s:%d\033[39m Country: \033[34m %s\033[39m ISP: \033[34m%s\033[39m\n", proxy, port, ipApi.Country, ipApi.Isp)
		}

	} else {
		fmt.Printf("\033[32mNew Proxy %s:%d\033[39m\n", proxy, port)
	}
}
