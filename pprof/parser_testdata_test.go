package pprof

import (
	"io"
	"os"
	"syscall"
	"testing"

	"compress/gzip"
	"github.com/stretchr/testify/require"
)

type testdataReader func(t testing.TB, fname string) ([]byte, func())

func readLabels(t testing.TB, td testdata, r testdataReader) (*LabelsSnapshot, func()) {
	ls := new(LabelsSnapshot)
	if td.labels != "" {
		labelsBytes, cleanup := r(t, testdataDir+td.labels)
		err := ls.UnmarshalVT(labelsBytes)
		require.NoError(t, err)
		return ls, cleanup
	}
	return ls, func() {

	}
}

func testDataReaders() []testdataReader {
	readers := []testdataReader{heapReader(), poorManSanitizerReader()}
	return readers
}

func heapReader() func(t testing.TB, fname string) ([]byte, func()) {
	return func(t testing.TB, fname string) ([]byte, func()) {
		return readGzipFile(t, fname), func() {
		}
	}
}

// todo guard with buildtag
func poorManSanitizerReader() func(t testing.TB, fname string) ([]byte, func()) {
	return func(t testing.TB, fname string) ([]byte, func()) {
		data, _ := heapReader()(t, fname)
		mmapedData, err := syscall.Mmap(-1, 0, len(data), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_ANON|syscall.MAP_PRIVATE)
		require.NoError(t, err)
		copy(mmapedData, data)
		require.Len(t, mmapedData, len(data))
		return mmapedData, func() {
			err = syscall.Mprotect(mmapedData, syscall.PROT_NONE)
			require.NoError(t, err)

		}
	}
}

func readGzipFile(t testing.TB, fname string) []byte {
	f, err := os.Open(fname)
	require.NoError(t, err)
	defer f.Close()
	r, err := gzip.NewReader(f)
	require.NoError(t, err)
	defer r.Close()
	bs, err := io.ReadAll(r)
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
