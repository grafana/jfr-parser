package format

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/grafana/jfr-parser/pprof"
)

type formatterPprof struct{}

func NewFormatterPprof() *formatterPprof {
	return &formatterPprof{}
}

func (f *formatterPprof) Format(buf []byte, dest string) ([]string, [][]byte, error) {
	pi := &pprof.ParseInput{
		StartTime:  time.Now(),
		EndTime:    time.Now(),
		SampleRate: 100,
	}
	profiles, err := pprof.ParseJFR(buf, pi, nil)
	if err != nil {
		return nil, nil, err
	}

	data := make([][]byte, 0)
	dests := make([]string, 0)
	destDir := filepath.Dir(dest)
	destBase := filepath.Base(dest)
	for i := 0; i < len(profiles.Profiles); i++ {
		filename := fmt.Sprintf("%s.%s", profiles.Profiles[i].Metric, destBase)
		dests = append(dests, filepath.Join(destDir, filename))

		bs, err := profiles.Profiles[i].Profile.MarshalVT()
		if err != nil {
			return nil, nil, err
		}
		data = append(data, bs)
	}
	return dests, data, nil
}
