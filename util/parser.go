package util

import (
	"fmt"
	"io"
	"unsafe"

	"github.com/grafana/jfr-parser/parser/types/def"
)

func ParseVarInt(buf []byte, pos *int) (uint32, error) {
	v := uint32(0)
	for shift := uint(0); ; shift += 7 {
		if shift >= 32 {
			return 0, def.ErrIntOverflow
		}
		if *pos >= len(buf) {
			return 0, io.ErrUnexpectedEOF
		}
		b := buf[*pos]
		*pos++
		v |= uint32(b&0x7F) << shift
		if b < 0x80 {
			break
		}
	}
	return v, nil
}

func ParseVarLong(buf []byte, pos *int) (uint64, error) {
	v := uint64(0)
	for shift := uint(0); shift <= 56; shift += 7 {
		if *pos >= len(buf) {
			return 0, io.ErrUnexpectedEOF
		}
		b_ := buf[*pos]
		*pos++
		if shift == 56 {
			v |= uint64(b_&0xFF) << shift
			break
		} else {
			v |= uint64(b_&0x7F) << shift
			if b_ < 0x80 {
				break
			}
		}
	}
	return v, nil
}

func ParseString(buf []byte, pos *int) (string, error) {
	if *pos >= len(buf) {
		return "", io.ErrUnexpectedEOF
	}

	b := buf[*pos]
	*pos++

	switch b { //todo implement 2
	case 0:
		return "", nil //todo this should be nil
	case 1:
		return "", nil
	case 3:
		bs, err := ParseBytes(buf, pos)
		if err != nil {
			return "", err
		}
		str := *(*string)(unsafe.Pointer(&bs))
		return str, nil
	case 4:
		return ParseCharArrayString(buf, pos)
	default:
		return "", fmt.Errorf("unknown string type %d", b)
	}
}

func ParseByte(buf []byte, pos *int) (byte, error) {
	if *pos >= len(buf) {
		return 0, io.ErrUnexpectedEOF
	}
	b := buf[*pos]
	*pos++
	return b, nil

}

func ParseBytes(buf []byte, pos *int) ([]byte, error) {
	l, err := ParseVarInt(buf, pos)
	if err != nil {
		return nil, err
	}
	if *pos+int(l) > len(buf) {
		return nil, io.ErrUnexpectedEOF
	}
	bs := buf[*pos : *pos+int(l)]
	*pos += int(l)
	return bs, nil
}

func ParseCharArrayString(buf []byte, pos *int) (string, error) {
	l, err := ParseVarInt(buf, pos)
	if err != nil {
		return "", err
	}
	chars := make([]rune, int(l))
	for i := 0; i < int(l); i++ {
		c, err := ParseVarInt(buf, pos)
		if err != nil {
			return "", err
		}
		chars[i] = rune(c)
	}

	res := string(chars)
	return res, nil
}
