package pprof

import (
	"time"

	profilev1 "github.com/grafana/pyroscope/api/gen/proto/go/google/v1"
	typesv1 "github.com/grafana/pyroscope/api/gen/proto/go/types/v1"
)

type Labels []*typesv1.LabelPair

type ParseInput struct {
	StartTime  time.Time
	EndTime    time.Time
	SampleRate int64
}

type Profiles struct {
	Profiles []Profile
	JFREvent string
}

type Profile struct {
	Profile *profilev1.Profile
	Metric  string
}
