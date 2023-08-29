package types

import (
	"fmt"
	"io"
	"unsafe"

	"github.com/pyroscope-io/jfr-parser/parser/types/def"
)

type StackTraceRef uint32
type StackTraceList struct {
	IDMap      map[StackTraceRef]uint32
	StackTrace []StackTrace
}

var ExpectedMetaStackTrace = &def.Class{
	Name: "jdk.types.StackTrace",
	ID:   def.T_STACK_TRACE,
	Fields: []def.Field{
		{
			Name:         "truncated",
			Type:         def.T_BOOLEAN,
			ConstantPool: false,
			Array:        false,
		},
		{
			Name:         "frames",
			Type:         def.T_STACK_FRAME,
			ConstantPool: false,
			Array:        true,
		},
	},
}

type StackTrace struct {
	Truncated bool
	Frames    []StackFrame
}

func (this *StackTraceList) Parse(data []byte, typ *def.Class, stackFrameType *def.Class, typeMap map[def.TypeID]*def.Class) (pos int, err error) {
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
	nFields := len(ExpectedMetaStackTrace.Fields)
	skipFields := typ.Fields[nFields:]
	stackFrameSkipFields := stackFrameType.Fields[len(def.TypeStackFrame.Fields):]
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
	this.IDMap = make(map[StackTraceRef]uint32, n)
	this.StackTrace = make([]StackTrace, n)
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
		id := StackTraceRef(v32_)
		if pos >= l {
			return 0, io.ErrUnexpectedEOF
		}
		b_ = data[pos]
		pos++
		this.StackTrace[i].Truncated = b_ == 0
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
		m := int(v32_)
		this.StackTrace[i].Frames = make([]StackFrame, m)
		for j := 0; j < m; j++ {
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
			this.StackTrace[i].Frames[j].Method = MethodRef(v32_)

			// skipping added fields
			for skipFI := range stackFrameSkipFields {
				nSkip := int(1)
				if stackFrameSkipFields[skipFI].Array {
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
					if stackFrameSkipFields[skipFI].ConstantPool {
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
						switch stackFrameSkipFields[skipFI].Type {
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
							gt := typeMap[stackFrameSkipFields[skipFI].Type]
							if gt == nil {
								return 0, fmt.Errorf("unknown type %d", stackFrameSkipFields[skipFI].Type)
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
			this.StackTrace[i].Frames[j].LineNumber = v32_

			// skipping added fields
			for skipFI := range stackFrameSkipFields {
				nSkip := int(1)
				if stackFrameSkipFields[skipFI].Array {
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
					if stackFrameSkipFields[skipFI].ConstantPool {
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
						switch stackFrameSkipFields[skipFI].Type {
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
							gt := typeMap[stackFrameSkipFields[skipFI].Type]
							if gt == nil {
								return 0, fmt.Errorf("unknown type %d", stackFrameSkipFields[skipFI].Type)
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
			this.StackTrace[i].Frames[j].BytecodeIndex = v32_

			// skipping added fields
			for skipFI := range stackFrameSkipFields {
				nSkip := int(1)
				if stackFrameSkipFields[skipFI].Array {
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
					if stackFrameSkipFields[skipFI].ConstantPool {
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
						switch stackFrameSkipFields[skipFI].Type {
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
							gt := typeMap[stackFrameSkipFields[skipFI].Type]
							if gt == nil {
								return 0, fmt.Errorf("unknown type %d", stackFrameSkipFields[skipFI].Type)
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
			this.StackTrace[i].Frames[j].Type = FrameTypeRef(v32_)

			// skipping added fields
			for skipFI := range stackFrameSkipFields {
				nSkip := int(1)
				if stackFrameSkipFields[skipFI].Array {
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
					if stackFrameSkipFields[skipFI].ConstantPool {
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
						switch stackFrameSkipFields[skipFI].Type {
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
							gt := typeMap[stackFrameSkipFields[skipFI].Type]
							if gt == nil {
								return 0, fmt.Errorf("unknown type %d", stackFrameSkipFields[skipFI].Type)
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
		}

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
		this.IDMap[id] = uint32(i)
	}
	return pos, nil
}