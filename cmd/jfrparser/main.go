package main

import (
	"flag"
	"os"

	"github.com/grafana/jfr-parser/cmd/jfrparser/format"
	"github.com/grafana/jfr-parser/parser"
)

const (
	formatJson = "Json"
)

type command struct {
	// Opts
	format string

	// Args
	src    string
	target string
}

func parseCommand(c *command) {
	format := flag.String("format", formatJson, "output format")
	flag.Parse()
	c.format = *format

	args := flag.Args()
	c.src = args[0]
	c.target = args[1]
}

type formatter interface {
	// Formats the given JFR
	Format(*parser.Parser) ([]byte, error)
}

// TODO: mov jfr2pprof here

// Usage: ./jfrparser -format=Json /path/to/jfr /path/to/json
func main() {
	c := new(command)
	parseCommand(c)

	buf, err := os.ReadFile(c.src)
	if err != nil {
		panic(err)
	}
	p := parser.NewParser(buf, parser.Options{})

	var fmtr formatter = nil
	switch c.format {
	case formatJson:
		fmtr = format.NewFormatterJson()
	default:
		panic("unsupported format")
	}

	out, err := fmtr.Format(p)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(c.target, out, 0644)
	if err != nil {
		panic(err)
	}
}
