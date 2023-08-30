package types

import (
	"fmt"
	"github.com/pyroscope-io/jfr-parser/parser/types/def"
	"io"
	"unsafe"
)

type BindClass struct {
	Temp   Class
	Fields []BindFieldClass
}

type BindFieldClass struct {
	Field          *def.Field
	ClassLoaderRef *ClassLoaderRef
	SymbolRef      *SymbolRef
	PackageRef     *PackageRef
	uint32         *uint32
}

func NewBindClass(typ *def.Class, typeMap *def.TypeMap) (*BindClass, error) {
	res := new(BindClass)
	for i := 0; i < len(typ.Fields); i++ {
		if typ.Fields[i].ConstantPool && typ.Fields[i].Array {
			return nil, fmt.Errorf("unimplemented cp && array")
		}
		switch typ.Fields[i].Name {
		case "classLoader":
			if typ.Fields[i].Equals(&def.Field{Name: "classLoader", Type: typeMap.T_CLASS_LOADER, ConstantPool: true, Array: false}) {
				res.Fields = append(res.Fields, BindFieldClass{Field: &typ.Fields[i], ClassLoaderRef: &res.Temp.ClassLoader})
			} else {
				res.Fields = append(res.Fields, BindFieldClass{Field: &typ.Fields[i]}) // skip
			}
		case "name":
			if typ.Fields[i].Equals(&def.Field{Name: "name", Type: typeMap.T_SYMBOL, ConstantPool: true, Array: false}) {
				res.Fields = append(res.Fields, BindFieldClass{Field: &typ.Fields[i], SymbolRef: &res.Temp.Name})
			} else {
				res.Fields = append(res.Fields, BindFieldClass{Field: &typ.Fields[i]}) // skip
			}
		case "package":
			if typ.Fields[i].Equals(&def.Field{Name: "package", Type: typeMap.T_PACKAGE, ConstantPool: true, Array: false}) {
				res.Fields = append(res.Fields, BindFieldClass{Field: &typ.Fields[i], PackageRef: &res.Temp.Package})
			} else {
				res.Fields = append(res.Fields, BindFieldClass{Field: &typ.Fields[i]}) // skip
			}
		case "modifiers":
			if typ.Fields[i].Equals(&def.Field{Name: "modifiers", Type: typeMap.T_INT, ConstantPool: false, Array: false}) {
				res.Fields = append(res.Fields, BindFieldClass{Field: &typ.Fields[i], uint32: &res.Temp.Modifiers})
			} else {
				res.Fields = append(res.Fields, BindFieldClass{Field: &typ.Fields[i]}) // skip
			}
		default:
			res.Fields = append(res.Fields, BindFieldClass{Field: &typ.Fields[i]}) // skip
		}
	}
	return res, nil
}

type ClassRef uint32
type ClassList struct {
	IDMap map[ClassRef]uint32
	Class []Class
}

type Class struct {
	ClassLoader ClassLoaderRef
	Name        SymbolRef
	Package     PackageRef
	Modifiers   uint32
}

func (this *ClassList) Parse(data []byte, bind *BindClass, typeMap *def.TypeMap) (pos int, err error) {
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
	this.IDMap = make(map[ClassRef]uint32, n)
	this.Class = make([]Class, n)
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
		id := ClassRef(v32_)
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
				case typeMap.T_CLASS_LOADER:
					if bind.Fields[bindFieldIndex].ClassLoaderRef != nil {
						*bind.Fields[bindFieldIndex].ClassLoaderRef = ClassLoaderRef(v32_)
					}
				case typeMap.T_SYMBOL:
					if bind.Fields[bindFieldIndex].SymbolRef != nil {
						*bind.Fields[bindFieldIndex].SymbolRef = SymbolRef(v32_)
					}
				case typeMap.T_PACKAGE:
					if bind.Fields[bindFieldIndex].PackageRef != nil {
						*bind.Fields[bindFieldIndex].PackageRef = PackageRef(v32_)
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
		this.Class[i] = bind.Temp
		this.IDMap[id] = uint32(i)
	}
	return pos, nil
}
