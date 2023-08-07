package parser

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/pyroscope-io/jfr-parser/reader"
)

var testfiles = []string{
	"example",
	"async-profiler", // -e cpu -i 10ms --alloc 512k --wall 200ms --lock 10ms -d 60 (async-profiler 2.10)
}

func TestParse(t *testing.T) {
	for _, testfile := range testfiles {
		jfrfile := testfile + ".jfr.gz"
		jsonfile := testfile + "_parsed.json.gz"
		jfr, err := readGzipFile("./testdata/" + jfrfile)
		if err != nil {
			t.Fatalf("Unable to read JFR file: %s", err)
		}
		expectedJson, err := readGzipFile("./testdata/" + jsonfile)
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
}

func TestParseBaseTypeAndDrop(t *testing.T) {
	r := reader.NewReader([]byte{1}, false, false)
	err := parseFields(
		r,
		map[int]*ClassMetadata{}, map[int]*CPool{},
		&ClassMetadata{
			Fields: []FieldMetadata{
				{
					Name:                 "boolean",
					isBaseType:           true,
					parseBaseTypeAndDrop: parseBaseTypeAndDrops["boolean"],
				},
			},
		},
		nil, false,
		func(reader reader.Reader, s string, resolvable ParseResolvable) error {
			return nil
		})
	if err != nil || r.Offset() != 1 {
		t.Fatalf("failed to parse and drop base type: %s", err)
	}
}

func BenchmarkParse(b *testing.B) {
	for _, testfile := range testfiles {
		b.Run(testfile, func(b *testing.B) {
			jfr, err := readGzipFile("./testdata/" + testfile + ".jfr.gz")
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
		})
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
