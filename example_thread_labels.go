package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/grafana/jfr-parser/pprof"
)

func main() {
	// Read JFR file
	jfrData, err := os.ReadFile("your-jfr-file.jfr")
	if err != nil {
		log.Fatal(err)
	}

	// Set up parse input
	parseInput := &pprof.ParseInput{
		StartTime:  time.Now(),
		EndTime:    time.Now(),
		SampleRate: 100,
	}

	// Parse JFR with thread name labels enabled
	profiles, err := pprof.ParseJFR(jfrData, parseInput, nil, pprof.WithThreadNameLabels(true))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Parsed %d profiles with thread name labels\n", len(profiles.Profiles))
	
	// You can also enable both thread root frames and thread labels:
	// profiles, err := pprof.ParseJFR(jfrData, parseInput, nil, 
	//     pprof.WithThreadRootFrame(true), 
	//     pprof.WithThreadNameLabels(true))
}