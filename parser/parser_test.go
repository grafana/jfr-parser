package parser

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	jfr, err := readGzipFile("./testdata/example.jfr.gz")
	if err != nil {
		t.Fatalf("Unable to read JFR file: %s", err)
	}
	expectedJson, err := readGzipFile("./testdata/example_parsed.json.gz")
	if err != nil {
		t.Fatalf("Unable to read example_parsd.json")
	}
	chunks, err := Parse(bytes.NewReader(jfr))
	if err != nil {
		t.Fatalf("Failed to parse JFR: %s", err)
		return
	}
	actualJson, _ := json.Marshal(chunks)
	if !bytes.Equal(expectedJson, actualJson) {
		t.Fatalf("Failed to parse JFR: %s", err)
		return
	}
}

func BenchmarkParse(b *testing.B) {
	jfr, err := readGzipFile("./testdata/example.jfr.gz")
	if err != nil {
		b.Fatalf("Unable to read JFR file: %s", err)
	}
	for i := 0; i < b.N; i++ {
		_, err := Parse(bytes.NewReader(jfr))
		if err != nil {
			b.Fatalf("Unable to parse JFR file: %s", err)
		}
	}
}

func readGzipFile(fname string) ([]byte, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r, err := gzip.NewReader(f)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return ioutil.ReadAll(r)
}
