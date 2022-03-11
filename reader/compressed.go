package reader

import (
	"encoding/binary"
	"fmt"
	"io"
)

type compressed struct {
	io.ByteReader
}

func newCompressed(r io.ByteReader) VarReader {
	return compressed{ByteReader: r}
}

func (c compressed) VarShort() (int16, error) {
	n, err := binary.ReadUvarint(c)
	if err != nil {
		return 0, err
	}
	if (n >> 48) > 0 {
		// TODO
		return 0, fmt.Errorf("overflow: %d bigger than 32 bits", n)
	}
	return int16(n), nil
	/*
		   FIXME: Is it unsigned LEB128?
		x := int16(n >> 1)
		if n&1 == 1 {
			x = ^x
		}
		return x, nil
	*/
}

func (c compressed) VarInt() (int32, error) {
	n, err := binary.ReadUvarint(c)
	if err != nil {
		return 0, err
	}
	if (n >> 32) > 0 {
		// TODO
		return 0, fmt.Errorf("overflow: %d bigger than 32 bits", n)
	}
	return int32(n), nil
	/*
		   FIXME: Is it unsigned LEB128?
		x := int32(n >> 1)
		if n&1 == 1 {
			x = ^x
		}
		return x, nil
	*/
}

func (c compressed) VarLong() (int64, error) {
	/*
			   FIXME: Is it unsigned LEB128?
		return binary.ReadVarint(c)
	*/
	n, err := binary.ReadUvarint(c)
	return int64(n), err
}
