package types

import (
	"fmt"
	"io"
	"unsafe"

	"github.com/pyroscope-io/jfr-parser/parser/types/def"
)

type BindMethod struct {
	Temp   Method
	Fields []BindFieldMethod
}

type BindFieldMethod struct {
	Field     *def.Field
	ClassRef  *ClassRef
	SymbolRef *SymbolRef
	uint32    *uint32
	bool      *bool
}

func NewBindMethod(typ *def.Class, typeMap *def.TypeMap) *BindMethod {
	res := new(BindMethod)
	for i := 0; i < len(typ.Fields); i++ {
		switch typ.Fields[i].Name {
		case "type":
			if typ.Fields[i].Equals(&def.Field{Name: "type", Type: typeMap.T_CLASS, ConstantPool: true, Array: false}) {
				res.Fields = append(res.Fields, BindFieldMethod{Field: &typ.Fields[i], ClassRef: &res.Temp.Type})
			} else {
				res.Fields = append(res.Fields, BindFieldMethod{Field: &typ.Fields[i]}) // skip
			}
		case "name":
			if typ.Fields[i].Equals(&def.Field{Name: "name", Type: typeMap.T_SYMBOL, ConstantPool: true, Array: false}) {
				res.Fields = append(res.Fields, BindFieldMethod{Field: &typ.Fields[i], SymbolRef: &res.Temp.Name})
			} else {
				res.Fields = append(res.Fields, BindFieldMethod{Field: &typ.Fields[i]}) // skip
			}
		case "descriptor":
			if typ.Fields[i].Equals(&def.Field{Name: "descriptor", Type: typeMap.T_SYMBOL, ConstantPool: true, Array: false}) {
				res.Fields = append(res.Fields, BindFieldMethod{Field: &typ.Fields[i], SymbolRef: &res.Temp.Descriptor})
			} else {
				res.Fields = append(res.Fields, BindFieldMethod{Field: &typ.Fields[i]}) // skip
			}
		case "modifiers":
			if typ.Fields[i].Equals(&def.Field{Name: "modifiers", Type: typeMap.T_INT, ConstantPool: false, Array: false}) {
				res.Fields = append(res.Fields, BindFieldMethod{Field: &typ.Fields[i], uint32: &res.Temp.Modifiers})
			} else {
				res.Fields = append(res.Fields, BindFieldMethod{Field: &typ.Fields[i]}) // skip
			}
		case "hidden":
			if typ.Fields[i].Equals(&def.Field{Name: "hidden", Type: typeMap.T_BOOLEAN, ConstantPool: false, Array: false}) {
				res.Fields = append(res.Fields, BindFieldMethod{Field: &typ.Fields[i], bool: &res.Temp.Hidden})
			} else {
				res.Fields = append(res.Fields, BindFieldMethod{Field: &typ.Fields[i]}) // skip
			}
		default:
			res.Fields = append(res.Fields, BindFieldMethod{Field: &typ.Fields[i]}) // skip
		}
	}
	return res
}

type MethodRef uint32
type MethodList struct {
	IDMap  IDMap[MethodRef]
	Method []Method
}

type Method struct {
	Type       ClassRef
	Name       SymbolRef
	Descriptor SymbolRef
	Modifiers  uint32
	Hidden     bool
	Scratch    string
}

func (this *MethodList) Parse(data []byte, bind *BindMethod, typeMap *def.TypeMap) (pos int, err error) {
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
	this.IDMap = NewIDMap[MethodRef](n)
	this.Method = make([]Method, n)
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
		id := MethodRef(v32_)
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
					case typeMap.T_CLASS:
						if bind.Fields[bindFieldIndex].ClassRef != nil {
							*bind.Fields[bindFieldIndex].ClassRef = ClassRef(v32_)
						}
					case typeMap.T_SYMBOL:
						if bind.Fields[bindFieldIndex].SymbolRef != nil {
							*bind.Fields[bindFieldIndex].SymbolRef = SymbolRef(v32_)
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
						if bind.Fields[bindFieldIndex].uint32 != nil {
							*bind.Fields[bindFieldIndex].uint32 = v32_
						}
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
						// skipping
					} else if bindFieldTypeID == typeMap.T_BOOLEAN {
						if pos >= l {
							return 0, io.ErrUnexpectedEOF
						}
						b_ = data[pos]
						pos++
						if bind.Fields[bindFieldIndex].bool != nil {
							*bind.Fields[bindFieldIndex].bool = b_ != 0
						}
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
		this.Method[i] = bind.Temp
		this.IDMap.Set(id, i)
	}
	return pos, nil
}
