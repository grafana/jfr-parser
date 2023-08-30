package types

import (
	"fmt"
	"io"
	"unsafe"

	"github.com/pyroscope-io/jfr-parser/parser/types/def"
)

type BindThreadPark struct {
	Temp   ThreadPark
	Fields []BindFieldThreadPark
}

type BindFieldThreadPark struct {
	Field         *def.Field
	uint64        *uint64
	ThreadRef     *ThreadRef
	StackTraceRef *StackTraceRef
	ClassRef      *ClassRef
}

func NewBindThreadPark(typ *def.Class, typeMap *def.TypeMap) *BindThreadPark {
	res := new(BindThreadPark)
	for i := 0; i < len(typ.Fields); i++ {
		switch typ.Fields[i].Name {
		case "startTime":
			if typ.Fields[i].Equals(&def.Field{Name: "startTime", Type: typeMap.T_LONG, ConstantPool: false, Array: false}) {
				res.Fields = append(res.Fields, BindFieldThreadPark{Field: &typ.Fields[i], uint64: &res.Temp.StartTime})
			} else {
				res.Fields = append(res.Fields, BindFieldThreadPark{Field: &typ.Fields[i]}) // skip
			}
		case "duration":
			if typ.Fields[i].Equals(&def.Field{Name: "duration", Type: typeMap.T_LONG, ConstantPool: false, Array: false}) {
				res.Fields = append(res.Fields, BindFieldThreadPark{Field: &typ.Fields[i], uint64: &res.Temp.Duration})
			} else {
				res.Fields = append(res.Fields, BindFieldThreadPark{Field: &typ.Fields[i]}) // skip
			}
		case "eventThread":
			if typ.Fields[i].Equals(&def.Field{Name: "eventThread", Type: typeMap.T_THREAD, ConstantPool: true, Array: false}) {
				res.Fields = append(res.Fields, BindFieldThreadPark{Field: &typ.Fields[i], ThreadRef: &res.Temp.EventThread})
			} else {
				res.Fields = append(res.Fields, BindFieldThreadPark{Field: &typ.Fields[i]}) // skip
			}
		case "stackTrace":
			if typ.Fields[i].Equals(&def.Field{Name: "stackTrace", Type: typeMap.T_STACK_TRACE, ConstantPool: true, Array: false}) {
				res.Fields = append(res.Fields, BindFieldThreadPark{Field: &typ.Fields[i], StackTraceRef: &res.Temp.StackTrace})
			} else {
				res.Fields = append(res.Fields, BindFieldThreadPark{Field: &typ.Fields[i]}) // skip
			}
		case "parkedClass":
			if typ.Fields[i].Equals(&def.Field{Name: "parkedClass", Type: typeMap.T_CLASS, ConstantPool: true, Array: false}) {
				res.Fields = append(res.Fields, BindFieldThreadPark{Field: &typ.Fields[i], ClassRef: &res.Temp.ParkedClass})
			} else {
				res.Fields = append(res.Fields, BindFieldThreadPark{Field: &typ.Fields[i]}) // skip
			}
		case "timeout":
			if typ.Fields[i].Equals(&def.Field{Name: "timeout", Type: typeMap.T_LONG, ConstantPool: false, Array: false}) {
				res.Fields = append(res.Fields, BindFieldThreadPark{Field: &typ.Fields[i], uint64: &res.Temp.Timeout})
			} else {
				res.Fields = append(res.Fields, BindFieldThreadPark{Field: &typ.Fields[i]}) // skip
			}
		case "until":
			if typ.Fields[i].Equals(&def.Field{Name: "until", Type: typeMap.T_LONG, ConstantPool: false, Array: false}) {
				res.Fields = append(res.Fields, BindFieldThreadPark{Field: &typ.Fields[i], uint64: &res.Temp.Until})
			} else {
				res.Fields = append(res.Fields, BindFieldThreadPark{Field: &typ.Fields[i]}) // skip
			}
		case "address":
			if typ.Fields[i].Equals(&def.Field{Name: "address", Type: typeMap.T_LONG, ConstantPool: false, Array: false}) {
				res.Fields = append(res.Fields, BindFieldThreadPark{Field: &typ.Fields[i], uint64: &res.Temp.Address})
			} else {
				res.Fields = append(res.Fields, BindFieldThreadPark{Field: &typ.Fields[i]}) // skip
			}
		case "contextId":
			if typ.Fields[i].Equals(&def.Field{Name: "contextId", Type: typeMap.T_LONG, ConstantPool: false, Array: false}) {
				res.Fields = append(res.Fields, BindFieldThreadPark{Field: &typ.Fields[i], uint64: &res.Temp.ContextId})
			} else {
				res.Fields = append(res.Fields, BindFieldThreadPark{Field: &typ.Fields[i]}) // skip
			}
		default:
			res.Fields = append(res.Fields, BindFieldThreadPark{Field: &typ.Fields[i]}) // skip
		}
	}
	return res
}

type ThreadPark struct {
	StartTime   uint64
	Duration    uint64
	EventThread ThreadRef
	StackTrace  StackTraceRef
	ParkedClass ClassRef
	Timeout     uint64
	Until       uint64
	Address     uint64
	ContextId   uint64
}

func (this *ThreadPark) Parse(data []byte, bind *BindThreadPark, typeMap *def.TypeMap) (pos int, err error) {
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
		bindArraySize := 1
		if bind.Fields[bindFieldIndex].Field.Array {
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
			bindArraySize = int(v32_)
		}
		for bindArrayIndex := 0; bindArrayIndex < bindArraySize; bindArrayIndex++ {
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
				case typeMap.T_THREAD:
					if bind.Fields[bindFieldIndex].ThreadRef != nil {
						*bind.Fields[bindFieldIndex].ThreadRef = ThreadRef(v32_)
					}
				case typeMap.T_STACK_TRACE:
					if bind.Fields[bindFieldIndex].StackTraceRef != nil {
						*bind.Fields[bindFieldIndex].StackTraceRef = StackTraceRef(v32_)
					}
				case typeMap.T_CLASS:
					if bind.Fields[bindFieldIndex].ClassRef != nil {
						*bind.Fields[bindFieldIndex].ClassRef = ClassRef(v32_)
					}
				}
			} else {
				bindFieldTypeID := bind.Fields[bindFieldIndex].Field.Type
				if bindFieldTypeID == typeMap.T_STRING {
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
				} else if bindFieldTypeID == typeMap.T_INT {
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
				} else if bindFieldTypeID == typeMap.T_LONG {
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
					if bind.Fields[bindFieldIndex].uint64 != nil {
						*bind.Fields[bindFieldIndex].uint64 = v64_
					}
				} else if bindFieldTypeID == typeMap.T_BOOLEAN {
					if pos >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b_ = data[pos]
					pos++
					// skipping
				} else if bindFieldTypeID == typeMap.T_FLOAT {
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
					bindFieldType := typeMap.IDMap[bind.Fields[bindFieldIndex].Field.Type]
					if bindFieldType == nil || len(bindFieldType.Fields) == 0 {
						return 0, fmt.Errorf("unknown type %d", bind.Fields[bindFieldIndex].Field.Type)
					}
					bindSkipObjects := 1
					if bind.Fields[bindFieldIndex].Field.Array {
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
						bindSkipObjects = int(v32_)
					}
					for bindSkipObjectIndex := 0; bindSkipObjectIndex < bindSkipObjects; bindSkipObjectIndex++ {
						for bindskipFieldIndex := 0; bindskipFieldIndex < len(bindFieldType.Fields); bindskipFieldIndex++ {
							bindSkipFieldType := bindFieldType.Fields[bindskipFieldIndex].Type
							if bindFieldType.Fields[bindskipFieldIndex].ConstantPool {
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
							} else if bindSkipFieldType == typeMap.T_STRING {
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
							} else if bindSkipFieldType == typeMap.T_INT {
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
							} else if bindSkipFieldType == typeMap.T_FLOAT {
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
							} else if bindSkipFieldType == typeMap.T_LONG {
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
							} else if bindSkipFieldType == typeMap.T_BOOLEAN {
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
		}
	}
	*this = bind.Temp
	return pos, nil
}
