package types

import (
	"fmt"
	"github.com/pyroscope-io/jfr-parser/parser/types/def"
	"io"
	"unsafe"
)

type BindStackFrame struct {
	Temp   StackFrame
	Fields []BindFieldStackFrame
}

type BindFieldStackFrame struct {
	Field        *def.Field
	MethodRef    *MethodRef
	uint32       *uint32
	FrameTypeRef *FrameTypeRef
}

func NewBindStackFrame(typ *def.Class, typeMap *def.TypeMap) (*BindStackFrame, error) {
	res := new(BindStackFrame)
	for i := 0; i < len(typ.Fields); i++ {
		if typ.Fields[i].ConstantPool && typ.Fields[i].Array {
			return nil, fmt.Errorf("unimplemented cp && array")
		}
		switch typ.Fields[i].Name {
		case "method":
			if typ.Fields[i].Equals(&def.Field{Name: "method", Type: typeMap.T_METHOD, ConstantPool: true, Array: false}) {
				res.Fields = append(res.Fields, BindFieldStackFrame{Field: &typ.Fields[i], MethodRef: &res.Temp.Method})
			} else {
				res.Fields = append(res.Fields, BindFieldStackFrame{Field: &typ.Fields[i]}) // skip
			}
		case "lineNumber":
			if typ.Fields[i].Equals(&def.Field{Name: "lineNumber", Type: typeMap.T_INT, ConstantPool: false, Array: false}) {
				res.Fields = append(res.Fields, BindFieldStackFrame{Field: &typ.Fields[i], uint32: &res.Temp.LineNumber})
			} else {
				res.Fields = append(res.Fields, BindFieldStackFrame{Field: &typ.Fields[i]}) // skip
			}
		case "bytecodeIndex":
			if typ.Fields[i].Equals(&def.Field{Name: "bytecodeIndex", Type: typeMap.T_INT, ConstantPool: false, Array: false}) {
				res.Fields = append(res.Fields, BindFieldStackFrame{Field: &typ.Fields[i], uint32: &res.Temp.BytecodeIndex})
			} else {
				res.Fields = append(res.Fields, BindFieldStackFrame{Field: &typ.Fields[i]}) // skip
			}
		case "type":
			if typ.Fields[i].Equals(&def.Field{Name: "type", Type: typeMap.T_FRAME_TYPE, ConstantPool: true, Array: false}) {
				res.Fields = append(res.Fields, BindFieldStackFrame{Field: &typ.Fields[i], FrameTypeRef: &res.Temp.Type})
			} else {
				res.Fields = append(res.Fields, BindFieldStackFrame{Field: &typ.Fields[i]}) // skip
			}
		default:
			res.Fields = append(res.Fields, BindFieldStackFrame{Field: &typ.Fields[i]}) // skip
		}
	}
	return res, nil
}

type StackFrame struct {
	Method        MethodRef
	LineNumber    uint32
	BytecodeIndex uint32
	Type          FrameTypeRef
}

func (this *StackFrame) Parse(data []byte, bind *BindStackFrame, typeMap *def.TypeMap) (pos int, err error) {
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
			switch bind.Fields[bindFieldIndex].Field.Type {
			case typeMap.T_METHOD:
				if bind.Fields[bindFieldIndex].MethodRef != nil {
					*bind.Fields[bindFieldIndex].MethodRef = MethodRef(v32_)
				}
			case typeMap.T_FRAME_TYPE:
				if bind.Fields[bindFieldIndex].FrameTypeRef != nil {
					*bind.Fields[bindFieldIndex].FrameTypeRef = FrameTypeRef(v32_)
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
				// skipping
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
				if bind.Fields[bindFieldIndex].uint32 != nil {
					*bind.Fields[bindFieldIndex].uint32 = v32_
				}
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
	return pos, nil
}
