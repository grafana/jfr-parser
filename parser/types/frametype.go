package types

import (
	"fmt"
	"github.com/pyroscope-io/jfr-parser/parser/types/def"
	"io"
	"unsafe"
)

type BindFrameType struct {
	Temp   FrameType
	Fields []BindFieldFrameType
}

type BindFieldFrameType struct {
	Field  *def.Field
	string *string
}

func NewBindFrameType(typ *def.Class, typeMap *def.TypeMap) (*BindFrameType, error) {
	res := new(BindFrameType)
	for i := 0; i < len(typ.Fields); i++ {
		if typ.Fields[i].ConstantPool && typ.Fields[i].Array {
			return nil, fmt.Errorf("unimplemented cp && array")
		}
		switch typ.Fields[i].Name {
		case "description":
			if typ.Fields[i].Equals(&def.Field{Name: "description", Type: typeMap.T_STRING, ConstantPool: false, Array: false}) {
				res.Fields = append(res.Fields, BindFieldFrameType{Field: &typ.Fields[i], string: &res.Temp.Description})
			} else {
				res.Fields = append(res.Fields, BindFieldFrameType{Field: &typ.Fields[i]}) // skip
			}
		default:
			res.Fields = append(res.Fields, BindFieldFrameType{Field: &typ.Fields[i]}) // skip
		}
	}
	return res, nil
}

type FrameTypeRef uint32
type FrameTypeList struct {
	IDMap     map[FrameTypeRef]uint32
	FrameType []FrameType
}

type FrameType struct {
	Description string
}

func (this *FrameTypeList) Parse(data []byte, bind *BindFrameType, typeMap *def.TypeMap) (pos int, err error) {
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
	this.IDMap = make(map[FrameTypeRef]uint32, n)
	this.FrameType = make([]FrameType, n)
	for i := 0; i < n; i++ {
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
		id := FrameTypeRef(v32_)
		for bindFieldIndex := 0; bindFieldIndex < len(bind.Fields); bindFieldIndex++ {
			if bind.Fields[bindFieldIndex].Field.ConstantPool {
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
			} else {
				bft := bind.Fields[bindFieldIndex].Field.Type
				if bft == typeMap.T_STRING {
					s_ = ""
					if pos >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b_ = data[pos]
					pos++
					switch b_ {
					case 0:
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
					default:
						return 0, fmt.Errorf("unknown string type %d at %d", b_, pos)
					}
					if bind.Fields[bindFieldIndex].string != nil {
						*bind.Fields[bindFieldIndex].string = s_
					}
				} else if bft == typeMap.T_INT {
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
					// skipping
				} else if bft == typeMap.T_LONG {
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
					// skipping
				} else if bft == typeMap.T_BOOLEAN {
					if pos >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b_ = data[pos]
					pos++
					// skipping
				} else if bft == typeMap.T_FLOAT {
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
					// skipping
				} else {
					fieldTyp := typeMap.IDMap[bind.Fields[bindFieldIndex].Field.Type]
					if fieldTyp == nil {
						return 0, fmt.Errorf("unknown type %d", bind.Fields[bindFieldIndex].Field.Type)
					}
					for skipFieldIndex := 0; skipFieldIndex < len(fieldTyp.Fields); skipFieldIndex++ {
						skipFieldType := fieldTyp.Fields[skipFieldIndex].Type
						if skipFieldType == typeMap.T_STRING {
							s_ = ""
							if pos >= l {
								return 0, io.ErrUnexpectedEOF
							}
							b_ = data[pos]
							pos++
							switch b_ {
							case 0:
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
							default:
								return 0, fmt.Errorf("unknown string type %d at %d", b_, pos)
							}
						} else if skipFieldType == typeMap.T_INT {
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
						} else if skipFieldType == typeMap.T_FLOAT {
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
						} else if skipFieldType == typeMap.T_LONG {
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
						} else if skipFieldType == typeMap.T_BOOLEAN {
							if pos >= l {
								return 0, io.ErrUnexpectedEOF
							}
							b_ = data[pos]
							pos++
						} else {
							return 0, fmt.Errorf("nested objects not implemented. ")
						}
					}
				}
			}
		}
		this.FrameType[i] = bind.Temp
		this.IDMap[id] = uint32(i)
	}
	return pos, nil
}
