package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/grafana/jfr-parser/cmd/jfrparser/format"
)

const (
	formatJson  = "Json"
	formatPprof = "Pprof"
)

type command struct {
	// Opts
	format string

	// Args
	src  string
	dest string
}

func parseCommand(c *command) {
	format := flag.String("format", formatJson, "output format")
	flag.Parse()
	c.format = *format

	args := flag.Args()
	c.src = args[0]
	c.dest = args[1]
}

type formatter interface {
	// Formats the given JFR
	Format(buf []byte, dest string) ([]string, [][]byte, error)
}

// Usage: ./jfrparser -format=Json /path/to/jfr /path/to/json
func main() {
	c := new(command)
	parseCommand(c)

	buf, err := os.ReadFile(c.src)
	if err != nil {
		panic(err)
	}

	var fmtr formatter = nil
	switch c.format {
	case formatJson:
		fmtr = format.NewFormatterJson()
	case formatPprof:
		fmtr = format.NewFormatterPprof()
	default:
		panic("unsupported format")
	}

	dests, data, err := fmtr.Format(buf, c.dest)
	if err != nil {
		panic(err)
	}

	if len(dests) != len(data) {
		panic(fmt.Errorf("logic error"))
	}

	for i := 0; i < len(dests); i++ {
		err = os.WriteFile(dests[i], data[i], 0644)
		if err != nil {
			panic(err)
		}
	}
}
