package pprof

import (
	"fmt"
	"io"

	"github.com/grafana/jfr-parser/parser"
)

func ParseJFR(body []byte, pi *ParseInput, jfrLabels *LabelsSnapshot) (res *Profiles, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("jfr parser panic: %v", r)
		}
	}()
	p := parser.NewParser(body, parser.Options{
		SymbolProcessor: parser.ProcessSymbols,
	})
	return parse(p, pi, jfrLabels)
}

func parse(parser *parser.Parser, piOriginal *ParseInput, jfrLabels *LabelsSnapshot) (result *Profiles, err error) {
	var event string

	builders := newJfrPprofBuilders(parser, jfrLabels, piOriginal)

	var values = [2]int64{1, 0}

	for {
		typ, err := parser.ParseEvent()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("jfr parser ParseEvent error: %w", err)
		}

		switch typ {
		case parser.TypeMap.T_EXECUTION_SAMPLE:
			ts := parser.GetThreadState(parser.ExecutionSample.State)
			ctx := ContextKey{
				ContextId: parser.ExecutionSample.ContextId,
				SpanId:    parser.ExecutionSample.SpanId,
				SpanName:  parser.ExecutionSample.SpanName,
			}
			if ts != nil && ts.Name != "STATE_SLEEPING" {
				builders.addStacktrace(sampleTypeCPU, ctx, parser.ExecutionSample.StackTrace, values[:1])
			}
			if event == "wall" {
				builders.addStacktrace(sampleTypeWall, ctx, parser.ExecutionSample.StackTrace, values[:1])
			}
		case parser.TypeMap.T_WALL_CLOCK_SAMPLE:
			values[0] = int64(parser.WallClockSample.Samples)
			builders.addStacktrace(sampleTypeWall, ContextKey{}, parser.WallClockSample.StackTrace, values[:1])
		case parser.TypeMap.T_ALLOC_IN_NEW_TLAB:
			values[1] = int64(parser.ObjectAllocationInNewTLAB.TlabSize)
			ctx := ContextKey{
				ContextId: parser.ObjectAllocationInNewTLAB.ContextId,
				SpanId:    parser.ObjectAllocationInNewTLAB.SpanId,
				SpanName:  parser.ObjectAllocationInNewTLAB.SpanName,
			}
			builders.addStacktrace(sampleTypeInTLAB, ctx, parser.ObjectAllocationInNewTLAB.StackTrace, values[:2])
		case parser.TypeMap.T_ALLOC_OUTSIDE_TLAB:
			values[1] = int64(parser.ObjectAllocationOutsideTLAB.AllocationSize)
			ctx := ContextKey{
				ContextId: parser.ObjectAllocationOutsideTLAB.ContextId,
				SpanId:    parser.ObjectAllocationOutsideTLAB.SpanId,
				SpanName:  parser.ObjectAllocationOutsideTLAB.SpanName,
			}
			builders.addStacktrace(sampleTypeOutTLAB, ctx, parser.ObjectAllocationOutsideTLAB.StackTrace, values[:2])
		case parser.TypeMap.T_ALLOC_SAMPLE:
			values[1] = int64(parser.ObjectAllocationSample.Weight)
			builders.addStacktrace(sampleTypeAllocSample, ContextKey{}, parser.ObjectAllocationSample.StackTrace, values[:2])
		case parser.TypeMap.T_MONITOR_ENTER:
			values[1] = int64(parser.JavaMonitorEnter.Duration)
			ctx := ContextKey{
				ContextId: parser.JavaMonitorEnter.ContextId,
				SpanId:    parser.JavaMonitorEnter.SpanId,
				SpanName:  parser.JavaMonitorEnter.SpanName,
			}
			builders.addStacktrace(sampleTypeLock, ctx, parser.JavaMonitorEnter.StackTrace, values[:2])
		case parser.TypeMap.T_THREAD_PARK:
			values[1] = int64(parser.ThreadPark.Duration)
			builders.addStacktrace(sampleTypeThreadPark, ContextKey{}, parser.ThreadPark.StackTrace, values[:2])
		case parser.TypeMap.T_LIVE_OBJECT:
			builders.addStacktrace(sampleTypeLiveObject, ContextKey{}, parser.LiveObject.StackTrace, values[:1])
		case parser.TypeMap.T_MALLOC:
			values[1] = int64(parser.Malloc.Size)
			builders.addStacktrace(sampleTypeMalloc, ContextKey{}, parser.Malloc.StackTrace, values[:2])
		case parser.TypeMap.T_ACTIVE_SETTING:
			if parser.ActiveSetting.Name == "event" {
				event = parser.ActiveSetting.Value
			}

		}
	}

	result = builders.build(event)

	return result, nil
}
