package format

import (
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/grafana/jfr-parser/parser"
	"github.com/stretchr/testify/assert"
)

type test struct {
	name     string
	pathJfr  string
	pathJson string
}

func loadTestDataGzip(t *testing.T, filename string) []byte {
	f, err := os.Open(filename)
	assert.NoError(t, err)
	defer f.Close()
	r, err := gzip.NewReader(f)
	assert.NoError(t, err)
	defer r.Close()
	b, err := io.ReadAll(r)
	assert.NoError(t, err)
	return b
}

func TestFormatterJson(t *testing.T) {
	fmtr := NewFormatterJson()
	tests := []test{{
		"example",
		filepath.Join("..", "..", "..", "parser", "testdata", "cortex-dev-01__kafka-0__cpu__0.jfr.gz"),
		filepath.Join("testdata", "cortex-dev-01__kafka-0__cpu__0.json.gz"),
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := loadTestDataGzip(t, tt.pathJfr)
			expected := loadTestDataGzip(t, tt.pathJson)
			actual, err := fmtr.Format(parser.NewParser(in, parser.Options{}))
			assert.NoError(t, err)
			assert.Equal(t, 0, bytes.Compare(expected, actual))
		})
	}
}
