package parser

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/grafana/jfr-parser/parser/types"
)

var testfiles = []string{
	"example",
	"async-profiler", // -e cpu -i 10ms --alloc 512k --wall 200ms --lock 10ms -d 60 (async-profiler 2.10)
	"cortex-dev-01__kafka-0__cpu_lock0_alloc0__0",
	"goland",
	"goland-multichunk",
}

type ExpectedActiveSetting struct {
	Key, Value string
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
			t1 := time.Now()
			parser, err := NewParser(jfr, Options{})
			if err != nil {
				t.Fatalf("Failed to parse JFR: %s", err)
				return
			}

			e := new(Expected)
			stacktraceToFrames := func(stacktrace types.StackTraceRef) []string {
				st := parser.GetStacktrace(stacktrace)
				if st == nil {
					t.Fatalf("stacktrace not found: %d\n", stacktrace)
				}
				frames := make([]string, len(st.Frames))
				for i, frame := range st.Frames {
					m := parser.GetMethod(frame.Method)
					if m == nil {
						t.Fatalf("method not found: %d\n", frame.Method)
					}
					if m.Scratch == "" {
						cls := parser.GetClass(m.Type)
						if cls == nil {
							t.Fatalf("class not found: %d\n", m.Type)
						}
						symbolString := parser.GetSymbolString(cls.Name)
						getSymbolString := parser.GetSymbolString(m.Name)
						m.Scratch = symbolString + "." + getSymbolString
					}
					frames[i] = m.Scratch
				}
				return frames
			}

			for {
				typ, err := parser.ParseEvent()
				if err != nil {
					if errors.Is(err, io.EOF) {
						break
					}
					t.Fatalf("Failed to parse JFR: %s", err)
				}

				switch typ {
				case parser.TypeMap.T_EXECUTION_SAMPLE:
					e.ExecutionSample = append(e.ExecutionSample, ExpectedStacktrace{
						Frames:    stacktraceToFrames(parser.ExecutionSample.StackTrace),
						ContextID: int64(parser.ExecutionSample.ContextId),
					})
				case parser.TypeMap.T_ALLOC_IN_NEW_TLAB:
					e.AllocInTLAB = append(e.AllocInTLAB, ExpectedStacktrace{
						Frames:    stacktraceToFrames(parser.ObjectAllocationInNewTLAB.StackTrace),
						ContextID: int64(parser.ObjectAllocationInNewTLAB.ContextId),
					})
				case parser.TypeMap.T_ALLOC_OUTSIDE_TLAB:
					e.AllocOutsideTLAB = append(e.AllocOutsideTLAB, ExpectedStacktrace{
						Frames:    stacktraceToFrames(parser.ObjectAllocationOutsideTLAB.StackTrace),
						ContextID: int64(parser.ObjectAllocationOutsideTLAB.ContextId),
					})
				case parser.TypeMap.T_MONITOR_ENTER:
					e.MonitorEnter = append(e.MonitorEnter, ExpectedStacktrace{
						Frames:    stacktraceToFrames(parser.JavaMonitorEnter.StackTrace),
						ContextID: int64(parser.JavaMonitorEnter.ContextId),
					})
				case parser.TypeMap.T_THREAD_PARK:
					e.ThreadPark = append(e.ThreadPark, ExpectedStacktrace{
						Frames:    stacktraceToFrames(parser.ThreadPark.StackTrace),
						ContextID: int64(parser.ThreadPark.ContextId),
					})
				case parser.TypeMap.T_LIVE_OBJECT:
					e.LiveObject = append(e.LiveObject, ExpectedStacktrace{
						Frames:    stacktraceToFrames(parser.LiveObject.StackTrace),
						ContextID: 0,
					})
				case parser.TypeMap.T_ACTIVE_SETTING:
					e.ActiveSetting = append(e.ActiveSetting, ExpectedActiveSetting{
						Key:   parser.ActiveSetting.Name,
						Value: parser.ActiveSetting.Value,
					})

				}
			}
			t2 := time.Now()
			fmt.Println(t2.Sub(t1))
			actualJson, _ := json.Marshal(e)
			//os.WriteFile("./testdata/"+testfile+"_parsed.json", actualJson, 0644)
			if !bytes.Equal(expectedJson, actualJson) {
				os.WriteFile("./"+testfile+"_actual.json", actualJson, 0644)
				os.WriteFile("./"+testfile+"_expected.json", expectedJson, 0644)

				t.Fatalf("Failed to parse JFR: %s", err)
				return
			}
		})
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
				parser, err := NewParser(jfr, Options{})

				stacktraceToFrames := func(stacktrace types.StackTraceRef) []string {
					st := parser.GetStacktrace(stacktrace)
					if st == nil {
						b.Fatalf("stacktrace not found: %d\n", stacktrace)
					}
					//frames := make([]string, len(st.Frames))
					for _, frame := range st.Frames {
						m := parser.GetMethod(frame.Method)
						if m == nil {
							b.Fatalf("method not found: %d\n", frame.Method)
						}
						if m.Scratch == "" {
							cls := parser.GetClass(m.Type)
							if cls == nil {
								b.Fatalf("class not found: %d\n", m.Type)
							}
							symbolString := parser.GetSymbolString(cls.Name)
							getSymbolString := parser.GetSymbolString(m.Name)
							_ = symbolString
							_ = getSymbolString
							//m.Scratch = symbolString + "." + getSymbolString
							m.Scratch = "once"
						}
						//frames[i] = m.Scratch
						//return ni
					}
					return nil
				}

				if err != nil {
					b.Fatalf("Unable to parse JFR file: %s", err)
				}
				for {
					typ, err := parser.ParseEvent()
					if err != nil {
						if errors.Is(err, io.EOF) {
							break
						}
						b.Fatalf("Unable to parse JFR file: %s", err)
					}

					switch typ {
					case parser.TypeMap.T_EXECUTION_SAMPLE:
						_ = stacktraceToFrames(parser.ExecutionSample.StackTrace)
					case parser.TypeMap.T_ALLOC_IN_NEW_TLAB:
						_ = stacktraceToFrames(parser.ObjectAllocationInNewTLAB.StackTrace)
					case parser.TypeMap.T_ALLOC_OUTSIDE_TLAB:
						_ = stacktraceToFrames(parser.ObjectAllocationOutsideTLAB.StackTrace)
					case parser.TypeMap.T_MONITOR_ENTER:
						_ = stacktraceToFrames(parser.JavaMonitorEnter.StackTrace)
					case parser.TypeMap.T_THREAD_PARK:
						_ = stacktraceToFrames(parser.ThreadPark.StackTrace)
					case parser.TypeMap.T_LIVE_OBJECT:
						_ = stacktraceToFrames(parser.LiveObject.StackTrace)

					case parser.TypeMap.T_ACTIVE_SETTING:

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
