/*
	(c) Yariya
*/

package main

import (
	"log"
	"os"
)

type Exporter struct {
	f   *os.File
	out string
}

var exporter *Exporter

func (e *Exporter) create() {
	var err error
	e.f, err = os.OpenFile(e.out, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalln(err)
	}
}

func (e *Exporter) Close() {
	e.f.WriteString("\n")
	e.f.Close()
}

func (e *Exporter) Add(s string) {
	_, err := e.f.WriteString(s + "\n")
	if err != nil {
		log.Println(err)
		return
	}
}
