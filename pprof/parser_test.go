package pprof

import (
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"slices"
	"strings"
	"testing"
	"time"

	gpprof "github.com/google/pprof/profile"
	profilev1 "github.com/grafana/pyroscope/api/gen/proto/go/google/v1"
	"github.com/k0kubun/pp/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testdataDir = "../parser/testdata/"
const doDump = false

type testdata struct {
	jfr, labels   string
	expectedCount int
}

var testfiles = []testdata{
	{"example", "", 4},
	{"async-profiler", "", 3}, // -e cpu -i 10ms --alloc 512k --wall 200ms --lock 10ms -d 60 (async-profiler 2.10)
	{"goland", "", 5},
	{"goland-multichunk", "", 5},
	{"FastSlow_2024_01_16_180855", "", 2}, // from IJ Ultimate, multichunk, chunked CP
	{"cortex-dev-01__kafka-0__cpu__0", "", 1},
	{"cortex-dev-01__kafka-0__cpu__1", "", 1},
	{"cortex-dev-01__kafka-0__cpu__2", "", 1},
	{"cortex-dev-01__kafka-0__cpu__3", "", 1},
	{"cortex-dev-01__kafka-0__cpu_lock0_alloc0__0", "", 5},
	{"cortex-dev-01__kafka-0__cpu_lock_alloc__0", "", 2},
	{"cortex-dev-01__kafka-0__cpu_lock_alloc__1", "", 2},
	{"cortex-dev-01__kafka-0__cpu_lock_alloc__2", "", 2},
	{"cortex-dev-01__kafka-0__cpu_lock_alloc__3", "", 2},
	{"dump1", "dump1.labels.pb.gz", 1},
	{"dump2", "dump2.labels.pb.gz", 4},
}

type gprofile struct {
	profile *gpprof.Profile
	proto   *profilev1.Profile
	metric  string
}

func TestDoDump(t *testing.T) {
	assert.False(t, doDump)
}

var parseInput = &ParseInput{
	StartTime:  time.Unix(1706241880, 0),
	EndTime:    time.Unix(1706241890, 0),
	SampleRate: 100,
}

func TestParse(t *testing.T) {
	for _, testfile := range testfiles {
		mypp := pp.New()
		mypp.SetColoringEnabled(false)
		mypp.SetExportedOnly(true)
		t.Run(testfile.jfr, func(t *testing.T) {
			jfrFile := testdataDir + testfile.jfr + ".jfr.gz"
			jfr := readGzipFile(t, jfrFile)
			ls := readLabels(t, testfile)

			profiles, err := ParseJFR(jfr, parseInput, ls)
			require.NoError(t, err)

			gprofiles := toGoogleProfiles(t, profiles.Profiles)
			profiles = nil

			slices.SortFunc(gprofiles, func(i, j gprofile) int {
				return strings.Compare(i.metric, j.metric)
			})
			assert.Equal(t, testfile.expectedCount, len(gprofiles))

			for i, profile := range gprofiles {
				actual := profileToString(t, profile)
				actualCollapsed := stackCollapseProto(profile.proto, true)
				expectedFile := fmt.Sprintf("%s%s_%d_%s_expected.txt.gz", testdataDir, testfile.jfr, i, profile.metric)
				expectedCollapsedFile := fmt.Sprintf("%s%s_%d_%s_expected_collapsed.txt.gz", testdataDir, testfile.jfr, i, profile.metric)
				assert.NotEmpty(t, actual)
				assert.NotEmpty(t, actualCollapsed)
				if doDump {
					writeGzipFile(t, expectedFile, []byte(actual))
					writeGzipFile(t, expectedCollapsedFile, []byte(actualCollapsed))
				} else {
					expected := readGzipFile(t, expectedFile)
					require.NoError(t, err)
					expectedCollapsed := readGzipFile(t, expectedCollapsedFile)
					require.NoError(t, err)

					assert.Equal(t, string(expected), actual)
					assert.Equal(t, string(expectedCollapsed), actualCollapsed)

					if string(expected) != actual {
						os.WriteFile("actual.txt", []byte(actual), 0644)
						os.WriteFile("expected.txt", expected, 0644)
					}

					if string(expectedCollapsed) != actualCollapsed {
						os.WriteFile("actual_collapsed.txt", []byte(actualCollapsed), 0644)
						os.WriteFile("expected_collapsed.txt", expectedCollapsed, 0644)
					}
				}
			}
		})
	}
}

func profileToString(t *testing.T, profile gprofile) string {
	res := profile.profile.String()
	re := regexp.MustCompile("\nTime: ([^\n]+)\n")
	matches := re.FindAllStringSubmatch(res, -1)
	require.Equal(t, 1, len(matches))
	t2, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", matches[0][1])
	assert.NoError(t, err)
	res = re.ReplaceAllString(res, fmt.Sprintf("\nTime: %d\n", t2.UnixNano()))
	return res
}

func readLabels(t testing.TB, td testdata) *LabelsSnapshot {
	ls := new(LabelsSnapshot)
	if td.labels != "" {
		labelsBytes := readGzipFile(t, testdataDir+td.labels)
		err := ls.UnmarshalVT(labelsBytes)
		require.NoError(t, err)
	}
	return ls
}

func BenchmarkParse(b *testing.B) {
	for _, testfile := range testfiles {
		b.Run(testfile.jfr, func(b *testing.B) {
			jfr := readGzipFile(b, testdataDir+testfile.jfr+".jfr.gz")
			ls := readLabels(b, testfile)
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				profiles, err := ParseJFR(jfr, parseInput, ls)
				if err != nil {
					b.Fatalf("Unable to parse JFR file: %s", err)
				}
				if len(profiles.Profiles) == 0 {
					b.Fatalf("No profiles found")
				}
			}
		})
	}
}

func toGoogleProfiles(t *testing.T, profiles []Profile) []gprofile {
	res := make([]gprofile, 0, len(profiles))
	for _, profile := range profiles {
		bs, err := profile.Profile.MarshalVT()
		require.NoError(t, err)
		p, err := gpprof.ParseData(bs)
		require.NoError(t, err)

		res = append(res, gprofile{p, profile.Profile, fmt.Sprintf("%s_%s", profile.Metric, sampleTypesToString(p))})
	}
	return res
}

func sampleTypesToString(p *gpprof.Profile) string {
	var sh1 string
	for _, s := range p.SampleType {
		dflt := ""
		sh1 = sh1 + fmt.Sprintf("%s__%s%s ", s.Type, s.Unit, dflt)
	}
	return strings.TrimSpace(sh1)
}

func readGzipFile(t testing.TB, fname string) []byte {
	f, err := os.Open(fname)
	require.NoError(t, err)
	defer f.Close()
	r, err := gzip.NewReader(f)
	require.NoError(t, err)
	defer r.Close()
	bs, err := ioutil.ReadAll(r)
	require.NoError(t, err)
	return bs
}

func writeGzipFile(t *testing.T, f string, data []byte) {
	fd, err := os.OpenFile(f, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	require.NoError(t, err)
	defer fd.Close()
	g := gzip.NewWriter(fd)
	_, err = g.Write(data)
	require.NoError(t, err)
	err = g.Close()
	require.NoError(t, err)
}

func stackCollapseProto(p *profilev1.Profile, lineNumbers bool) string {
	allZeros := func(a []int64) bool {
		for _, v := range a {
			if v != 0 {
				return false
			}
		}
		return true
	}
	addValues := func(a, b []int64) {
		for i := range a {
			a[i] += b[i]
		}
	}

	type stack struct {
		funcs string
		value []int64
	}
	locMap := make(map[int64]*profilev1.Location)
	funcMap := make(map[int64]*profilev1.Function)
	for _, l := range p.Location {
		locMap[int64(l.Id)] = l
	}
	for _, f := range p.Function {
		funcMap[int64(f.Id)] = f
	}

	var ret []stack
	for _, s := range p.Sample {
		var funcs []string
		for i := range s.LocationId {
			locID := s.LocationId[len(s.LocationId)-1-i]
			loc := locMap[int64(locID)]
			for _, line := range loc.Line {
				f := funcMap[int64(line.FunctionId)]
				fname := p.StringTable[f.Name]
				if lineNumbers {
					fname = fmt.Sprintf("%s:%d", fname, line.Line)
				}
				funcs = append(funcs, fname)
			}
		}

		vv := make([]int64, len(s.Value))
		copy(vv, s.Value)
		ret = append(ret, stack{
			funcs: strings.Join(funcs, ";"),
			value: vv,
		})
	}
	slices.SortFunc(ret, func(i, j stack) int {
		return strings.Compare(i.funcs, j.funcs)
	})
	var unique []stack
	for _, s := range ret {
		if allZeros(s.value) {
			continue
		}
		if len(unique) == 0 {
			unique = append(unique, s)
			continue
		}

		if unique[len(unique)-1].funcs == s.funcs {
			addValues(unique[len(unique)-1].value, s.value)
			continue
		}
		unique = append(unique, s)

	}

	res := make([]string, 0, len(unique))
	for _, s := range unique {
		res = append(res, fmt.Sprintf("%s %v", s.funcs, s.value))
	}
	return strings.Join(res, "\n")
}
