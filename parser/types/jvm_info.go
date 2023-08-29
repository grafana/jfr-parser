package types

import (
	"fmt"
	"io"
	"unsafe"

	"github.com/pyroscope-io/jfr-parser/parser/types/def"
)

var ExpectedMetaJVMInformation = &def.Class{
	Name: "jdk.JVMInformation",
	ID:   def.T_JVM_INFORMATION,
	Fields: []def.Field{
		{
			Name:         "startTime",
			Type:         def.T_LONG,
			ConstantPool: false,
			Array:        false,
		},
		{
			Name:         "jvmName",
			Type:         def.T_STRING,
			ConstantPool: false,
			Array:        false,
		},
		{
			Name:         "jvmVersion",
			Type:         def.T_STRING,
			ConstantPool: false,
			Array:        false,
		},
		{
			Name:         "jvmArguments",
			Type:         def.T_STRING,
			ConstantPool: false,
			Array:        false,
		},
		{
			Name:         "jvmFlags",
			Type:         def.T_STRING,
			ConstantPool: false,
			Array:        false,
		},
		{
			Name:         "javaArguments",
			Type:         def.T_STRING,
			ConstantPool: false,
			Array:        false,
		},
		{
			Name:         "jvmStartTime",
			Type:         def.T_LONG,
			ConstantPool: false,
			Array:        false,
		},
		{
			Name:         "pid",
			Type:         def.T_LONG,
			ConstantPool: false,
			Array:        false,
		},
	},
}

type JVMInformation struct {
	StartTime     uint64
	JvmName       string
	JvmVersion    string
	JvmArguments  string
	JvmFlags      string
	JavaArguments string
	JvmStartTime  uint64
	Pid           uint64
}

func (this *JVMInformation) Parse(data []byte, typ *def.Class, typeMap map[def.TypeID]*def.Class) (pos int, err error) {
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
	nFields := len(ExpectedMetaJVMInformation.Fields)
	skipFields := typ.Fields[nFields:]
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
	this.StartTime = v64_
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
	this.JvmName = s_
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
	this.JvmVersion = s_
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
	this.JvmArguments = s_
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
	this.JvmFlags = s_
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
	this.JavaArguments = s_
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
	this.JvmStartTime = v64_
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
	this.Pid = v64_

	// skipping added fields
	for skipFI := range skipFields {
		nSkip := int(1)
		if skipFields[skipFI].Array {
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
			nSkip = int(v32_)
		}
		for iSkip := 0; iSkip < nSkip; iSkip++ {
			if skipFields[skipFI].ConstantPool {
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
				switch skipFields[skipFI].Type {
				case def.T_STRING:
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
				case def.T_LONG:
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
				case def.T_INT:
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
				case def.T_FLOAT:
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
				case def.T_BOOLEAN:
					if pos >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b_ = data[pos]
					pos++
				default:
					gt := typeMap[skipFields[skipFI].Type]
					if gt == nil {
						return 0, fmt.Errorf("unknown type %d", skipFields[skipFI].Type)
					}
					for gti := 0; gti < len(gt.Fields); gti++ {
						if gt.Fields[gti].Array {
							return 0, fmt.Errorf("two dimentional array not supported")
						}
						if gt.Fields[gti].ConstantPool {
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
							switch gt.Fields[gti].Type {
							case def.T_STRING:
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
							case def.T_LONG:
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
							case def.T_INT:
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
							case def.T_FLOAT:
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
							case def.T_BOOLEAN:
								if pos >= l {
									return 0, io.ErrUnexpectedEOF
								}
								b_ = data[pos]
								pos++
							default:
								return 0, fmt.Errorf("unknown type %d", gt.Fields[gti].Type)
							}
						}
					}
				}
			}
		}
	}
	return pos, nil
}