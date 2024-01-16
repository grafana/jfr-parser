package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/grafana/jfr-parser/pprof"
)

func main() {
	jfrFile, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	dstFile := os.Args[2]
	pi := &pprof.ParseInput{
		StartTime:  time.Now(),
		EndTime:    time.Now(),
		SampleRate: 100,
	}
	profiles, err := pprof.ParseJFR(jfrFile, pi, nil)
	if err != nil {
		panic(err)
	}
	for _, profile := range profiles.Profiles {
		bs, err := profile.Profile.MarshalVT()
		if err != nil {
			panic(err)
		}
		dir := filepath.Dir(dstFile)
		filename := fmt.Sprintf("%s.%s", profile.Metric, filepath.Base(dstFile))
		err = os.WriteFile(filepath.Join(dir, filename), bs, 0644)
		if err != nil {
			panic(err)
		}
	}
}
