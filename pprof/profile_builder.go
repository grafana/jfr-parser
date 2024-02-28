package pprof

import (
	profilev1 "github.com/grafana/pyroscope/api/gen/proto/go/google/v1"
)

type ProfileBuilder struct {
	*profilev1.Profile
	strings                       map[string]int
	externalLocationID2LocationID map[ExternalLocationID]PPROFLocationID
	externalFunctionID2FunctionID map[ExternalFunctionID]PPROFFunctionID
	externalSampleID2SampleIndex  map[sampleID]uint32
	metricName                    string
}

type sampleID struct {
	locationsID uint64
	labelsID    uint64
}

// NewProfileBuilderWithLabels creates a new ProfileBuilder with the given nanoseconds timestamp and labels.
func NewProfileBuilderWithLabels(ts int64) *ProfileBuilder {
	profile := new(profilev1.Profile)
	profile.TimeNanos = ts
	profile.Mapping = append(profile.Mapping, &profilev1.Mapping{
		Id: 1, HasFunctions: true,
	})
	p := &ProfileBuilder{
		Profile:                       profile,
		strings:                       map[string]int{},
		externalFunctionID2FunctionID: map[ExternalFunctionID]PPROFFunctionID{},
		externalLocationID2LocationID: map[ExternalLocationID]PPROFLocationID{},
	}
	p.addString("")
	return p
}

type ExternalFunctionID uint32
type ExternalLocationID struct {
	ExternalFunctionID ExternalFunctionID
	Line               uint32
}
type PPROFFunctionID uint64
type PPROFLocationID uint64

func (m *ProfileBuilder) AddSampleType(typ, unit string) {
	m.Profile.SampleType = append(m.Profile.SampleType, &profilev1.ValueType{
		Type: m.addString(typ),
		Unit: m.addString(unit),
	})
}

func (m *ProfileBuilder) MetricName(name string) {
	m.metricName = name
}

func (m *ProfileBuilder) PeriodType(periodType string, periodUnit string) {
	m.Profile.PeriodType = &profilev1.ValueType{
		Type: m.addString(periodType),
		Unit: m.addString(periodUnit),
	}
}

func (m *ProfileBuilder) addString(s string) int64 {
	i, ok := m.strings[s]
	if !ok {
		i = len(m.strings)
		m.strings[s] = i
		m.StringTable = append(m.StringTable, s)
	}
	return int64(i)
}

func (m *ProfileBuilder) FindLocationByExternalID(externalLocationID ExternalLocationID) (PPROFLocationID, bool) {
	loc, ok := m.externalLocationID2LocationID[externalLocationID]
	return loc, ok
}

func (m *ProfileBuilder) FindFunctionByExternalID(externalFunctionID ExternalFunctionID) (PPROFFunctionID, bool) {
	loc, ok := m.externalFunctionID2FunctionID[externalFunctionID]
	return loc, ok
}

func (m *ProfileBuilder) AddExternalFunction(frame string, id ExternalFunctionID) PPROFFunctionID {
	fname := m.addString(frame)
	funcID := uint64(len(m.Function)) + 1
	m.Function = append(m.Function, &profilev1.Function{
		Id:   funcID,
		Name: fname,
	})
	ret := PPROFFunctionID(funcID)
	m.externalFunctionID2FunctionID[id] = ret
	return ret
}

func (m *ProfileBuilder) AddExternalLocation(id ExternalLocationID, pprofFunctionID PPROFFunctionID) PPROFLocationID {
	locID := uint64(len(m.Location)) + 1
	m.Location = append(m.Location, &profilev1.Location{
		Id:        locID,
		MappingId: uint64(1),
		Line:      []*profilev1.Line{{FunctionId: uint64(pprofFunctionID), Line: int64(id.Line)}},
	})
	ret := PPROFLocationID(locID)
	m.externalLocationID2LocationID[id] = ret
	return ret

}

func (m *ProfileBuilder) AddExternalSample(locs []uint64, values []int64, externalSampleID uint32) {
	m.AddExternalSampleWithLabels(locs, values, nil, nil, uint64(externalSampleID), 0)
}

func (m *ProfileBuilder) FindExternalSample(externalSampleID uint32) *profilev1.Sample {
	return m.FindExternalSampleWithLabels(uint64(externalSampleID), 0)
}

func (m *ProfileBuilder) AddExternalSampleWithLabels(locs []uint64, values []int64, labelsCtx *Context, labelsSnapshot *LabelsSnapshot, locationsID, labelsID uint64) {
	sample := &profilev1.Sample{
		LocationId: locs,
		Value:      values,
	}
	if m.externalSampleID2SampleIndex == nil {
		m.externalSampleID2SampleIndex = map[sampleID]uint32{}
	}
	m.externalSampleID2SampleIndex[sampleID{locationsID: locationsID, labelsID: labelsID}] = uint32(len(m.Profile.Sample))
	m.Profile.Sample = append(m.Profile.Sample, sample)
	if labelsCtx != nil && labelsSnapshot != nil {
		sample.Label = make([]*profilev1.Label, 0, len(labelsCtx.Labels))
		for k, v := range labelsCtx.Labels { //todo iterating over map is not deterministic, this can break tests and maybe even affect performance
			sample.Label = append(sample.Label, &profilev1.Label{
				Key: m.addString(labelsSnapshot.Strings[k]),
				Str: m.addString(labelsSnapshot.Strings[v]),
			})
		}
	}
}

func (m *ProfileBuilder) FindExternalSampleWithLabels(locationsID, labelsID uint64) *profilev1.Sample {
	sampleIndex, ok := m.externalSampleID2SampleIndex[sampleID{locationsID: locationsID, labelsID: labelsID}]
	if !ok {
		return nil
	}
	sample := m.Profile.Sample[sampleIndex]
	return sample
}
