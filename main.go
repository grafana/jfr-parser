package main

import (
	"fmt"
	"log"
	"os"

	"github.com/pyroscope-io/jfr-parser/parser"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please indicate the JFR file to parse")
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalf("Unable to open file: %s", err)
	}
	chunks, err := parser.Parse(f)
	if err != nil {
		log.Fatalf("Unable to parse: %s", err)
	}
	log.Printf("Parsed %d chunks", len(chunks))
	events := make(map[string]int)
	for _, c := range chunks {
		for _, e := range c.Events {
			events[fmt.Sprintf("%T", e)]++
		}
	}
	for k, v := range events {
		log.Printf("%s: %d\n", k, v)
	}
}
