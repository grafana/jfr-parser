package parser

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/pyroscope-io/jfr-parser/reader"
)

var testfiles = []string{
	"example",
	"async-profiler", // -e cpu -i 10ms --alloc 512k --wall 200ms --lock 10ms -d 60 (async-profiler 2.10)
	"cortex-dev-01__kafka-0__cpu_lock0_alloc0__0",
	"goland",
	"goland-multichunk",
}

type ExpectedActiveSetting struct {
	key, value string
}
type ExpectedStacktrace struct {
	Frames    []string
	ContextID int64
}
type Expected struct {
	ExecutionSample  []ExpectedStacktrace
	AllocInTLAB      []ExpectedStacktrace
	AllocOutsideTLAB []ExpectedStacktrace
	MonitorEnter     []ExpectedStacktrace
	ThreadPark       []ExpectedStacktrace
	LiveObject       []ExpectedStacktrace
	ActiveSetting    []ExpectedActiveSetting
}

func TestParse(t *testing.T) {
	for _, testfile := range testfiles {
		t.Run(testfile, func(t *testing.T) {
			jfrfile := "./testdata/" + testfile + ".jfr.gz"
			jsonfile := "./testdata/" + testfile + "_parsed.json.gz"
			jfr, err := readGzipFile(jfrfile)
			if err != nil {
				t.Fatalf("Unable to read JFR file: %s", err)
			}
			expectedJson, err := readGzipFile(jsonfile)
			if err != nil {
				t.Fatalf("Unable to read example_parsd.json")
			}
			chunks, err := Parse(bytes.NewReader(jfr))
			if err != nil {
				t.Fatalf("Failed to parse JFR: %s", err)
				return
			}

			e := new(Expected)
			for _, chunk := range chunks {
				for chunk.Next() {
					event := chunk.Event
					switch event.(type) {
					case *ExecutionSample:
						ex := event.(*ExecutionSample)
						e.ExecutionSample = append(e.ExecutionSample, ExpectedStacktrace{
							Frames:    stacktraceToFrames(ex.StackTrace),
							ContextID: ex.ContextId,
						})
					case *ObjectAllocationInNewTLAB:
						ex := event.(*ObjectAllocationInNewTLAB)
						e.AllocInTLAB = append(e.AllocInTLAB, ExpectedStacktrace{
							Frames:    stacktraceToFrames(ex.StackTrace),
							ContextID: ex.ContextId,
						})
					case *ObjectAllocationOutsideTLAB:
						ex := event.(*ObjectAllocationOutsideTLAB)
						e.AllocOutsideTLAB = append(e.AllocOutsideTLAB, ExpectedStacktrace{
							Frames:    stacktraceToFrames(ex.StackTrace),
							ContextID: ex.ContextId,
						})
					case *JavaMonitorEnter:
						ex := event.(*JavaMonitorEnter)
						e.MonitorEnter = append(e.MonitorEnter, ExpectedStacktrace{
							Frames:    stacktraceToFrames(ex.StackTrace),
							ContextID: ex.ContextId,
						})
					case *ThreadPark:
						ex := event.(*ThreadPark)
						e.ThreadPark = append(e.ThreadPark, ExpectedStacktrace{
							Frames:    stacktraceToFrames(ex.StackTrace),
							ContextID: ex.ContextId,
						})
					case *LiveObject:
						ex := event.(*LiveObject)
						e.LiveObject = append(e.LiveObject, ExpectedStacktrace{
							Frames:    stacktraceToFrames(ex.StackTrace),
							ContextID: 0,
						})

					}
				}
				err = chunk.Err()
				if err != nil {
					t.Fatal(err)
				}
			}
			actualJson, _ := json.Marshal(e)
			//os.WriteFile("./testdata/"+testfile+"_parsed.json", actualJson, 0644)
			if !bytes.Equal(expectedJson, actualJson) {
				t.Fatalf("Failed to parse JFR: %s", err)
				return
			}
		})
	}
}

func stacktraceToFrames(trace *StackTrace) []string {
	frames := make([]string, len(trace.Frames))
	for i, frame := range trace.Frames {
		frames[i] = frame.Method.Type.Name.String + "." + frame.Method.Name.String
	}
	return frames
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
			b.ResetTimer()
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
