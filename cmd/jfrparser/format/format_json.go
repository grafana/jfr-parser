package format

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/grafana/jfr-parser/parser"
)

type formatterJson struct {
	outBuf []byte
}

func NewFormatterJson() *formatterJson {
	return &formatterJson{make([]byte, 0)}
}

func (f *formatterJson) append(d any) error {
	var (
		tmp []byte
		err error
	)
	if tmp, err = json.Marshal(d); err != nil {
		return fmt.Errorf("json.Marshal error: %w", err)
	}
	f.outBuf = append(f.outBuf, tmp...)
	return nil
}

func (f *formatterJson) appendInit(p *parser.Parser) error {
	var err error
	f.outBuf = append(f.outBuf, []byte("[{\"Header\":")...)
	if err = f.append(p.ChunkHeader()); err != nil {
		return err
	}
	f.outBuf = append(f.outBuf, []byte(",\"FrameTypes\":")...)
	if err = f.append(p.FrameTypes); err != nil {
		return err
	}
	f.outBuf = append(f.outBuf, []byte(",\"ThreadStates\":")...)
	if err = f.append(p.ThreadStates); err != nil {
		return err
	}
	f.outBuf = append(f.outBuf, []byte(",\"Threads\":")...)
	if err = f.append(p.Threads); err != nil {
		return err
	}
	f.outBuf = append(f.outBuf, []byte(",\"Classes\":")...)
	if err = f.append(p.Classes); err != nil {
		return err
	}
	f.outBuf = append(f.outBuf, []byte(",\"Methods\":")...)
	if err = f.append(p.Methods); err != nil {
		return err
	}
	f.outBuf = append(f.outBuf, []byte(",\"Packages\":")...)
	if err = f.append(p.Packages); err != nil {
		return err
	}
	f.outBuf = append(f.outBuf, []byte(",\"Symbols\":")...)
	if err = f.append(p.Symbols); err != nil {
		return err
	}
	f.outBuf = append(f.outBuf, []byte(",\"LogLevels\":")...)
	if err = f.append(p.LogLevels); err != nil {
		return err
	}
	f.outBuf = append(f.outBuf, []byte(",\"Stacktraces\":")...)
	if err = f.append(p.Stacktrace); err != nil {
		return err
	}
	f.outBuf = append(f.outBuf, []byte(",\"Recordings\":[")...)
	return nil
}

// TODO: support multi-chunk JFR, by exposing new chunk indicator on the parser (a counter), and printing an array of chunks
func (f *formatterJson) Format(buf []byte, dest string) ([]string, [][]byte, error) {
	p := parser.NewParser(buf, parser.Options{})
	first := true
	for {
		typ, err := p.ParseEvent()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, nil, fmt.Errorf("parser.ParseEvent error: %w", err)
		}

		if first {
			if err = f.appendInit(p); err != nil {
				return nil, nil, err
			}
			first = false
		}

		switch typ {
		case p.TypeMap.T_EXECUTION_SAMPLE:
			if err = f.append(p.ExecutionSample); err != nil {
				return nil, nil, err
			}
		case p.TypeMap.T_ALLOC_IN_NEW_TLAB:
			if err = f.append(p.ObjectAllocationInNewTLAB); err != nil {
				return nil, nil, err
			}
		case p.TypeMap.T_ALLOC_OUTSIDE_TLAB:
			if err = f.append(p.ObjectAllocationOutsideTLAB); err != nil {
				return nil, nil, err
			}
		case p.TypeMap.T_MONITOR_ENTER:
			if err = f.append(p.JavaMonitorEnter); err != nil {
				return nil, nil, err
			}
		case p.TypeMap.T_THREAD_PARK:
			if err = f.append(p.ThreadPark); err != nil {
				return nil, nil, err
			}
		case p.TypeMap.T_LIVE_OBJECT:
			if err = f.append(p.LiveObject); err != nil {
				return nil, nil, err
			}
		case p.TypeMap.T_ACTIVE_SETTING:
			if err = f.append(p.ActiveSetting); err != nil {
				return nil, nil, err
			}
		}
		f.outBuf = append(f.outBuf, []byte(",")...)
	}

	f.outBuf = f.outBuf[:len(f.outBuf)-1]
	f.outBuf = append(f.outBuf, []byte("]}]\n")...)
	return []string{dest}, [][]byte{f.outBuf}, nil
}
