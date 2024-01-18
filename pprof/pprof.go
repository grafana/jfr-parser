package pprof

import (
	"github.com/grafana/jfr-parser/parser"
	"github.com/grafana/jfr-parser/parser/types"
)

const (
	sampleTypeCPU        = 0
	sampleTypeWall       = 1
	sampleTypeInTLAB     = 2
	sampleTypeOutTLAB    = 3
	sampleTypeLock       = 4
	sampleTypeThreadPark = 5
	sampleTypeLiveObject = 6
)

func newJfrPprofBuilders(p *parser.Parser, jfrLabels *LabelsSnapshot, piOriginal *ParseInput) *jfrPprofBuilders {
	st := piOriginal.StartTime.UnixNano()
	et := piOriginal.EndTime.UnixNano()
	var period int64
	if piOriginal.SampleRate == 0 {
		period = 0
	} else {
		period = 1e9 / int64(piOriginal.SampleRate)
	}

	res := &jfrPprofBuilders{
		parser:        p,
		builders:      make(map[int64]*ProfileBuilder),
		jfrLabels:     jfrLabels,
		timeNanos:     st,
		durationNanos: et - st,
		period:        period,
	}
	return res
}

type jfrPprofBuilders struct {
	parser        *parser.Parser
	builders      map[int64]*ProfileBuilder
	jfrLabels     *LabelsSnapshot
	timeNanos     int64
	durationNanos int64
	period        int64
}

func (b *jfrPprofBuilders) addStacktrace(sampleType int64, contextID uint64, ref types.StackTraceRef, values []int64) {
	p := b.profileBuilderForSampleType(sampleType)
	st := b.parser.GetStacktrace(ref)
	if st == nil {
		return
	}

	addValues := func(dst []int64) {
		mul := 1
		if sampleType == sampleTypeCPU || sampleType == sampleTypeWall {
			mul = int(b.period)
		}
		for i, value := range values {
			dst[i] += value * int64(mul)
		}
	}

	sample := p.FindExternalSampleWithLabels(uint64(ref), contextID)
	if sample != nil {
		addValues(sample.Value)
		return
	}

	locations := make([]uint64, 0, len(st.Frames))
	for i := 0; i < len(st.Frames); i++ {
		f := st.Frames[i]
		extLocID := ExternalLocationID{
			ExternalFunctionID: ExternalFunctionID(f.Method),
			Line:               f.LineNumber,
		}
		loc, found := p.FindLocationByExternalID(extLocID)
		if found {
			locations = append(locations, uint64(loc))
			continue
		}
		m := b.parser.GetMethod(f.Method)
		if m != nil {

			pprofFuncID, found := p.FindFunctionByExternalID(extLocID.ExternalFunctionID)
			if found {
				// add new location with old function
			} else {
				cls := b.parser.GetClass(m.Type)
				if cls == nil {
					continue
				}
				clsName := b.parser.GetSymbolString(cls.Name)
				methodName := b.parser.GetSymbolString(m.Name)
				frame := clsName + "." + methodName
				pprofFuncID = p.AddExternalFunction(frame, extLocID.ExternalFunctionID)
			}
			loc = p.AddExternalLocation(extLocID, pprofFuncID)
			locations = append(locations, uint64(loc))

			//todo remove Scratch field from the Method
		}
	}
	vs := make([]int64, len(values))
	addValues(vs)
	p.AddExternalSampleWithLabels(locations, vs, b.contextLabels(contextID), b.jfrLabels, uint64(ref), contextID)
}

func (b *jfrPprofBuilders) profileBuilderForSampleType(sampleType int64) *ProfileBuilder {
	if builder, ok := b.builders[sampleType]; ok {
		return builder
	}
	builder := NewProfileBuilderWithLabels(b.timeNanos)
	builder.DurationNanos = b.durationNanos
	var metric string
	switch sampleType {
	case sampleTypeCPU:
		builder.AddSampleType("cpu", "nanoseconds")
		builder.PeriodType("cpu", "nanoseconds")
		metric = "process_cpu"
	case sampleTypeWall:
		builder.AddSampleType("wall", "nanoseconds")
		builder.PeriodType("wall", "nanoseconds")
		metric = "wall"
	case sampleTypeInTLAB:
		builder.AddSampleType("alloc_in_new_tlab_objects", "count")
		builder.AddSampleType("alloc_in_new_tlab_bytes", "bytes")
		builder.PeriodType("space", "bytes")
		metric = "memory"
	case sampleTypeOutTLAB:
		builder.AddSampleType("alloc_outside_tlab_objects", "count")
		builder.AddSampleType("alloc_outside_tlab_bytes", "bytes")
		builder.PeriodType("space", "bytes")
		metric = "memory"
	case sampleTypeLock:
		builder.AddSampleType("contentions", "count")
		builder.AddSampleType("delay", "nanoseconds")
		builder.PeriodType("mutex", "count")
		metric = "mutex"
	case sampleTypeThreadPark:
		builder.AddSampleType("contentions", "count")
		builder.AddSampleType("delay", "nanoseconds")
		builder.PeriodType("block", "count")
		metric = "block"
	case sampleTypeLiveObject:
		builder.AddSampleType("live", "count")
		builder.PeriodType("objects", "count")
		metric = "memory"
	}
	builder.MetricName(metric)
	b.builders[sampleType] = builder
	return builder
}

func (b *jfrPprofBuilders) contextLabels(contextID uint64) *Context {
	if b.jfrLabels == nil {
		return nil
	}
	return b.jfrLabels.Contexts[int64(contextID)]
}

func (b *jfrPprofBuilders) build(jfrEvent string) *Profiles {
	profiles := make([]Profile, 0, len(b.builders))
	for _, builder := range b.builders {
		profiles = append(profiles, Profile{
			Profile: builder.Profile,
			Metric:  builder.metricName,
		})
	}
	return &Profiles{
		Profiles: profiles,
		JFREvent: jfrEvent,
	}
}
