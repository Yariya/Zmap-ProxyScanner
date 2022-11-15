/*
	(c) Yariya
*/

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func Scanner() {
	if *fetch != "" {
		log.Printf("Detected URL Mode.\n")
		res, err := http.Get(*fetch)
		if err != nil {
			log.Fatalln("fetch error")
		}
		body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatalln("fetch body error")
		}
		res.Body.Close()

		scanner := bufio.NewScanner(bytes.NewReader(body))
		for scanner.Scan() {
			ip := scanner.Text()
			queueChan <- ip
		}
	} else if *input != "" {
		fmt.Printf("Detected FILE Mode.\n")
		b, err := os.ReadFile(*input)
		if err != nil {
			log.Fatalln("open file err")
		}
		lines := strings.Split(string(b), "\n")
		for _, line := range lines {
			queueChan <- line
		}
	} else {
		fmt.Printf("Detected ZMAP Mode.\n")
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			ip := scanner.Text()
			queueChan <- ip
		}
	}
}
