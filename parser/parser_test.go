package parser

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
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
	type expectedChunk struct {
		Header      Header
		Metadata    MetadataEvent
		Checkpoints []CheckpointEvent
		Events      []Parseable
	}
	var expectedChunks []expectedChunk
	for _, chunk := range chunks {
		var events []Parseable
		for chunk.Next() {
			// copy event
			elem := reflect.ValueOf(chunk.Event).Elem()
			newEvent := reflect.New(elem.Type())
			newEventElem := newEvent.Elem()
			for i := 0; i < elem.NumField(); i++ {
				newEventElem.Field(i).Set(elem.Field(i))
			}
			events = append(events, newEvent.Interface().(Parseable))
		}
		err = chunk.Err()
		if err != nil {
			t.Fatal(err)
		}
		expectedChunks = append(expectedChunks, expectedChunk{
			Header:      chunk.Header,
			Metadata:    chunk.Metadata,
			Checkpoints: chunk.Checkpoints,
			Events:      events,
		})
	}
	actualJson, _ := json.Marshal(expectedChunks)
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
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		chunks, err := Parse(bytes.NewReader(jfr))
		if err != nil {
			b.Fatalf("Unable to parse JFR file: %s", err)
		}
		for _, chunk := range chunks {
			for chunk.Next() {
			}
			err = chunk.Err()
			if err != nil {
				b.Fatal(err)
			}
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
