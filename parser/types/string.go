// Code generated by gen/main.go. DO NOT EDIT.

package types

import (
	"fmt"
	"github.com/grafana/jfr-parser/parser/types/def"
	"io"
	"unsafe"
)

type BindString struct {
	Temp   String
	Fields []BindFieldString
}

type BindFieldString struct {
	Field *def.Field
}

func NewBindString(typ *def.Class, typeMap *def.TypeMap) *BindString {
	res := new(BindString)
	res.Fields = make([]BindFieldString, 0, len(typ.Fields))
	for i := 0; i < len(typ.Fields); i++ {
		switch typ.Fields[i].Name {
		default:
			res.Fields = append(res.Fields, BindFieldString{Field: &typ.Fields[i]}) // skip unknown new field
		}
	}
	return res
}

type StringRef uint64
type StringList struct {
	IDMap  map[StringRef]uint32
	String []String
}

type String struct {
	String string
}

func (this *StringList) Parse(data []byte, bind *BindString, typeMap *def.TypeMap) (pos int, err error) {
	var (
		v64_  uint64
		v32_  uint32
		s_    string
		b_    byte
		shift = uint(0)
		l     = len(data)
	)
	_ = v64_
	_ = v32_
	_ = s_
	v32_ = uint32(0)
	for shift = uint(0); ; shift += 7 {
		if shift >= 32 {
			return 0, def.ErrIntOverflow
		}
		if pos >= l {
			return 0, io.ErrUnexpectedEOF
		}
		b_ = data[pos]
		pos++
		v32_ |= uint32(b_&0x7F) << shift
		if b_ < 0x80 {
			break
		}
	}
	n := int(v32_)
	this.IDMap = make(map[StringRef]uint32, n)
	this.String = make([]String, n)
	for i := 0; i < n; i++ {
		v64_ = 0
		for shift = uint(0); shift <= 56; shift += 7 {
			if pos >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b_ = data[pos]
			pos++
			if shift == 56 {
				v64_ |= uint64(b_&0xFF) << shift
				break
			} else {
				v64_ |= uint64(b_&0x7F) << shift
				if b_ < 0x80 {
					break
				}
			}
		}
		id := StringRef(v64_)
		s_ = ""
		if pos >= l {
			return 0, io.ErrUnexpectedEOF
		}
		b_ = data[pos]
		pos++
		switch b_ {
		case 0:
			break
		case 1:
			break
		case 3:
			v32_ = uint32(0)
			for shift = uint(0); ; shift += 7 {
				if shift >= 32 {
					return 0, def.ErrIntOverflow
				}
				if pos >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b_ = data[pos]
				pos++
				v32_ |= uint32(b_&0x7F) << shift
				if b_ < 0x80 {
					break
				}
			}
			if pos+int(v32_) > l {
				return 0, io.ErrUnexpectedEOF
			}
			bs := data[pos : pos+int(v32_)]
			s_ = *(*string)(unsafe.Pointer(&bs))
			pos += int(v32_)
		case 5:
			v32_ = uint32(0)
			for shift = uint(0); ; shift += 7 {
				if shift >= 32 {
					return 0, def.ErrIntOverflow
				}
				if pos >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b_ = data[pos]
				pos++
				v32_ |= uint32(b_&0x7F) << shift
				if b_ < 0x80 {
					break
				}
			}
			if pos+int(v32_) > l {
				return 0, io.ErrUnexpectedEOF
			}
			bs := data[pos : pos+int(v32_)]
			bs, _ = typeMap.ISO8859_1Decoder.Bytes(bs)
			s_ = *(*string)(unsafe.Pointer(&bs))
			pos += int(v32_)
		case 4:
			v32_ = uint32(0)
			for shift = uint(0); ; shift += 7 {
				if shift >= 32 {
					return 0, def.ErrIntOverflow
				}
				if pos >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b_ = data[pos]
				pos++
				v32_ |= uint32(b_&0x7F) << shift
				if b_ < 0x80 {
					break
				}
			}
			bl := int(v32_)
			buf := make([]rune, bl)
			for i := 0; i < bl; i++ {
				v32_ = uint32(0)
				for shift = uint(0); ; shift += 7 {
					if shift >= 32 {
						return 0, def.ErrIntOverflow
					}
					if pos >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b_ = data[pos]
					pos++
					v32_ |= uint32(b_&0x7F) << shift
					if b_ < 0x80 {
						break
					}
				}
				buf[i] = rune(v32_)
			}
			s_ = string(buf)
		default:
			return 0, fmt.Errorf("unknown string type %d at %d", b_, pos)
		}
		bind.Temp.String = s_
		this.String[i] = bind.Temp
		this.IDMap[id] = uint32(i)
	}
	return pos, nil
}
