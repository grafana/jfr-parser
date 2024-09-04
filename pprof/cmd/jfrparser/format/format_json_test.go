package format

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
	testDataDir := filepath.Join("..", "..", "..", "parser", "testdata")
	fnamePrefix := "cortex-dev-01__kafka-0__cpu__0"
	dest := "example"

	fmtr := NewFormatterJson()
	tests := []struct {
		name     string
		pathJfr  string
		pathJson string
	}{{
		"example",
		filepath.Join(testDataDir, fmt.Sprintf("%s.jfr.gz", fnamePrefix)),
		filepath.Join(testDataDir, fmt.Sprintf("%s.json.gz", fnamePrefix)),
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := loadTestDataGzip(t, tt.pathJfr)
			expected := loadTestDataGzip(t, tt.pathJson)
			dests, data, err := fmtr.Format(in, dest)
			assert.NoError(t, err)
			assert.Equal(t, 1, len(dests))
			assert.True(t, dest == dests[0])
			assert.Equal(t, 1, len(data))
			assert.Equal(t, 0, bytes.Compare(expected, data[0]))
		})
	}
}
