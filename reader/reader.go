package reader

import (
	"encoding/binary"
	"fmt"
	"io"
)

type VarReader interface {
	VarShort() (int16, error)
	VarInt() (int32, error)
	VarLong() (int64, error)
}

type Reader interface {
	Boolean() (bool, error)
	Byte() (int8, error)
	Short() (int16, error)
	Char() (uint16, error)
	Int() (int32, error)
	Long() (int64, error)
	Float() (float32, error)
	Double() (float64, error)
	Bytes() ([]byte, error)

	VarReader

	// TODO: Support arrays
}

type InputReader interface {
	io.Reader
	io.ByteReader
}

type reader struct {
	InputReader
	varR  VarReader
	bytes []byte
}

func NewReader(r InputReader, compressed bool) Reader {
	var varR VarReader
	if compressed {
		varR = newCompressed(r)
	} else {
		varR = newUncompressed(r)
	}
	return reader{
		InputReader: r,
		varR:        varR,
		bytes:       make([]byte, 0),
	}
}

func (r reader) Boolean() (bool, error) {
	var n int8
	err := binary.Read(r, binary.BigEndian, &n)
	if n == 0 {
		return false, err
	}
	return true, err
}

func (r reader) Byte() (int8, error) {
	var n int8
	err := binary.Read(r, binary.BigEndian, &n)
	return n, err
}

func (r reader) Short() (int16, error) {
	return Short(r)
}

func (r reader) Char() (uint16, error) {
	var n uint16
	err := binary.Read(r, binary.BigEndian, &n)
	return n, err
}

func (r reader) Int() (int32, error) {
	return Int(r)
}

func (r reader) Long() (int64, error) {
	return Long(r)
}

func (r reader) Float() (float32, error) {
	var n float32
	err := binary.Read(r, binary.BigEndian, &n)
	return n, err
}

func (r reader) Double() (float64, error) {
	var n float64
	err := binary.Read(r, binary.BigEndian, &n)
	return n, err
}

// TODO: Should we differentiate between null and empty?
func (r reader) Bytes() ([]byte, error) {
	enc, err := r.Byte()
	if err != nil {
		return r.bytes[:0], err
	}
	switch enc {
	case 0:
		return r.bytes[:0], nil
	case 1:
		return r.bytes[:0], nil
	case 3, 4, 5:
		return r.utf8()
	default:
		// TODO
		return r.bytes[:0], fmt.Errorf("Unsupported string type :%d", enc)
	}
}

func (r reader) VarShort() (int16, error) {
	return r.varR.VarShort()
}

func (r reader) VarInt() (int32, error) {
	return r.varR.VarInt()
}

func (r reader) VarLong() (int64, error) {
	return r.varR.VarLong()
}

func (r reader) utf8() ([]byte, error) {
	n, err := r.varR.VarInt()
	if err != nil {
		return r.bytes, nil
	}

	if cap(r.bytes) < int(n) {

		r.bytes = make([]byte, int(n))
	} else {
		r.bytes = r.bytes[:n]
	}

	_, err = io.ReadFull(r, r.bytes)
	return r.bytes, err
}
