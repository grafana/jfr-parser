//go:build libfuzzer

package main

// #include <stdint.h>
import "C"
import (
	"encoding/binary"
	"github.com/grafana/jfr-parser/pprof"
	"time"
	"unsafe"
)

//export LLVMFuzzerInitialize
func LLVMFuzzerInitialize(argc *C.int, argv ***C.char) C.int {
	return 0
}

type fuzzdata struct {
	data []byte
}

func (f *fuzzdata) u8() uint8 {
	if len(f.data) == 0 {
		return 0
	}
	b := f.data[0]
	f.data = f.data[1:]
	return b
}
func (f *fuzzdata) bytes(sz int) []byte {
	if sz == 0 {
		return nil
	}
	if len(f.data) < sz {
		res := f.data
		f.data = nil
		return res
	}
	res := f.data[:sz]
	f.data = f.data[sz:]
	return res
}
func (f *fuzzdata) u64() uint64 {
	if len(f.data) < 8 {
		return 0
	}
	v := binary.LittleEndian.Uint64(f.data[0:8])
	f.data = f.data[8:]
	return v
}

//export LLVMFuzzerTestOneInput
func LLVMFuzzerTestOneInput(data *C.char, size C.size_t) C.int {
	gdata := unsafe.Slice((*byte)(unsafe.Pointer(data)), size)
	if len(gdata) == 0 {
		return 0
	}
	fd := fuzzdata{gdata}
	flags := fd.u8()
	withLabels := flags&1 == 1
	truncatedFrame := (flags>>1)&1 == 1
	var ls *pprof.LabelsSnapshot
	if withLabels {
		lsb := fd.bytes(int(fd.u8()))
		ls = &pprof.LabelsSnapshot{}
		_ = ls.UnmarshalVT(lsb)
	}
	pi := &pprof.ParseInput{
		StartTime:  time.UnixMilli(int64(fd.u64())),
		EndTime:    time.UnixMilli(int64(fd.u64())),
		SampleRate: int64(fd.u64()),
	}

	_, _ = pprof.ParseJFR(fd.bytes(len(gdata)), pi, ls, pprof.WithTruncatedFrame(truncatedFrame), pprof.WithDisablePanicRecovery(true))
	return 0
}

func main() {

}
