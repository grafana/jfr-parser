package types

import (
	"fmt"
	"io"
	"unsafe"

	"github.com/pyroscope-io/jfr-parser/parser/types/def"
)

type BindLogLevel struct {
	Temp   LogLevel
	Fields []BindFieldLogLevel
}

type BindFieldLogLevel struct {
	Field  *def.Field
	string *string
}

func NewBindLogLevel(typ *def.Class, typeMap *def.TypeMap) *BindLogLevel {
	res := new(BindLogLevel)
	for i := 0; i < len(typ.Fields); i++ {
		switch typ.Fields[i].Name {
		case "name":
			if typ.Fields[i].Equals(&def.Field{Name: "name", Type: typeMap.T_STRING, ConstantPool: false, Array: false}) {
				res.Fields = append(res.Fields, BindFieldLogLevel{Field: &typ.Fields[i], string: &res.Temp.Name})
			} else {
				res.Fields = append(res.Fields, BindFieldLogLevel{Field: &typ.Fields[i]}) // skip
			}
		default:
			res.Fields = append(res.Fields, BindFieldLogLevel{Field: &typ.Fields[i]}) // skip
		}
	}
	return res
}

type LogLevelRef uint32
type LogLevelList struct {
	IDMap    map[LogLevelRef]uint32
	LogLevel []LogLevel
}

type LogLevel struct {
	Name string
}

func (this *LogLevelList) Parse(data []byte, bind *BindLogLevel, typeMap *def.TypeMap) (pos int, err error) {
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
	this.IDMap = make(map[LogLevelRef]uint32, n)
	this.LogLevel = make([]LogLevel, n)
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
		id := LogLevelRef(v32_)
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
						if bind.Fields[bindFieldIndex].string != nil {
							*bind.Fields[bindFieldIndex].string = s_
						}
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
						// skipping
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
		this.LogLevel[i] = bind.Temp
		this.IDMap[id] = uint32(i)
	}
	return pos, nil
}
