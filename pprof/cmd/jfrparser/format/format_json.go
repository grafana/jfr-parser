package format

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/grafana/jfr-parser/parser"
	"github.com/grafana/jfr-parser/parser/types"
)

type formatterJson struct{}

func NewFormatterJson() *formatterJson {
	return &formatterJson{}
}

type chunk struct {
	Header       parser.ChunkHeader
	FrameTypes   types.FrameTypeList
	ThreadStates types.ThreadStateList
	Threads      types.ThreadList
	Classes      types.ClassList
	Methods      types.MethodList
	Packages     types.PackageList
	Symbols      types.SymbolList
	LogLevels    types.LogLevelList
	Stacktraces  types.StackTraceList
	Recordings   []any
}

func initChunk(c *chunk, p *parser.Parser) {
	c.Header = p.ChunkHeader()
	c.FrameTypes = p.FrameTypes
	c.ThreadStates = p.ThreadStates
	c.Threads = p.Threads
	c.Classes = p.Classes
	c.Methods = p.Methods
	c.Packages = p.Packages
	c.Symbols = p.Symbols
	c.LogLevels = p.LogLevels
	c.Stacktraces = p.Stacktrace
	c.Recordings = make([]any, 0)
}

// TODO: support multi-chunk JFR, by exposing new chunk indicator on the parser (a counter), and printing an array of chunks
func (f *formatterJson) Format(buf []byte, dest string) ([]string, [][]byte, error) {
	p := parser.NewParser(buf, parser.Options{SymbolProcessor: parser.ProcessSymbols})

	ir := make([]chunk, 1)
	chunkIdx := 0
	newChunk := true
	for {
		typ, err := p.ParseEvent()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, nil, fmt.Errorf("parser.ParseEvent error: %w", err)
		}

		if newChunk {
			initChunk(&ir[chunkIdx], p)
			newChunk = false
		}

		switch typ {
		case p.TypeMap.T_EXECUTION_SAMPLE:
			ir[chunkIdx].Recordings = append(ir[chunkIdx].Recordings, p.ExecutionSample)
		case p.TypeMap.T_ALLOC_IN_NEW_TLAB:
			ir[chunkIdx].Recordings = append(ir[chunkIdx].Recordings, p.ObjectAllocationInNewTLAB)
		case p.TypeMap.T_ALLOC_OUTSIDE_TLAB:
			ir[chunkIdx].Recordings = append(ir[chunkIdx].Recordings, p.ObjectAllocationOutsideTLAB)
		case p.TypeMap.T_MONITOR_ENTER:
			ir[chunkIdx].Recordings = append(ir[chunkIdx].Recordings, p.JavaMonitorEnter)
		case p.TypeMap.T_THREAD_PARK:
			ir[chunkIdx].Recordings = append(ir[chunkIdx].Recordings, p.ThreadPark)
		case p.TypeMap.T_LIVE_OBJECT:
			ir[chunkIdx].Recordings = append(ir[chunkIdx].Recordings, p.LiveObject)
		case p.TypeMap.T_ACTIVE_SETTING:
			ir[chunkIdx].Recordings = append(ir[chunkIdx].Recordings, p.ActiveSetting)
		}
	}

	outBuf, err := json.Marshal(ir)
	if err != nil {
		return nil, nil, fmt.Errorf("json.Marshal error: %w", err)
	}
	outBuf = append(outBuf, '\n')
	return []string{dest}, [][]byte{outBuf}, nil
}
