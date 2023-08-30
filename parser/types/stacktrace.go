package types

import (
	"fmt"
	"io"
	"unsafe"

	"github.com/grafana/jfr-parser/parser/types/def"
)

type BindStackTrace struct {
	Temp   StackTrace
	Fields []BindFieldStackTrace
}

type BindFieldStackTrace struct {
	Field      *def.Field
	bool       *bool
	StackFrame *[]StackFrame
}

func NewBindStackTrace(typ *def.Class, typeMap *def.TypeMap) *BindStackTrace {
	res := new(BindStackTrace)
	res.Fields = make([]BindFieldStackTrace, 0, len(typ.Fields))
	for i := 0; i < len(typ.Fields); i++ {
		switch typ.Fields[i].Name {
		case "truncated":
			if typ.Fields[i].Equals(&def.Field{Name: "truncated", Type: typeMap.T_BOOLEAN, ConstantPool: false, Array: false}) {
				res.Fields = append(res.Fields, BindFieldStackTrace{Field: &typ.Fields[i], bool: &res.Temp.Truncated})
			} else {
				res.Fields = append(res.Fields, BindFieldStackTrace{Field: &typ.Fields[i]}) // skip
			}
		case "frames":
			if typ.Fields[i].Equals(&def.Field{Name: "frames", Type: typeMap.T_STACK_FRAME, ConstantPool: false, Array: true}) {
				res.Fields = append(res.Fields, BindFieldStackTrace{Field: &typ.Fields[i], StackFrame: &res.Temp.Frames})
			} else {
				res.Fields = append(res.Fields, BindFieldStackTrace{Field: &typ.Fields[i]}) // skip
			}
		default:
			res.Fields = append(res.Fields, BindFieldStackTrace{Field: &typ.Fields[i]}) // skip
		}
	}
	return res
}

type StackTraceRef uint32
type StackTraceList struct {
	IDMap      map[StackTraceRef]uint32
	StackTrace []StackTrace
}

type StackTrace struct {
	Truncated bool
	Frames    []StackFrame
}

func (this *StackTraceList) Parse(data []byte, bind *BindStackTrace, bindStackFrame *BindStackFrame, typeMap *def.TypeMap) (pos int, err error) {
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
				if bind.Fields[bindFieldIndex].Field.Type == typeMap.T_STACK_FRAME {
					*bind.Fields[bindFieldIndex].StackFrame = make([]StackFrame, 0, bindArraySize)
				}
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
					switch bindFieldTypeID {
					case typeMap.T_STRING:
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
					case typeMap.T_INT:
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
					case typeMap.T_LONG:
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
					case typeMap.T_BOOLEAN:
						if pos >= l {
							return 0, io.ErrUnexpectedEOF
						}
						b_ = data[pos]
						pos++
						if bind.Fields[bindFieldIndex].bool != nil {
							*bind.Fields[bindFieldIndex].bool = b_ != 0
						}
					case typeMap.T_FLOAT:
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
					case typeMap.T_STACK_FRAME:
						for bindStackFrameFieldIndex := 0; bindStackFrameFieldIndex < len(bindStackFrame.Fields); bindStackFrameFieldIndex++ {
							bindStackFrameArraySize := 1
							if bindStackFrame.Fields[bindStackFrameFieldIndex].Field.Array {
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
								bindStackFrameArraySize = int(v32_)
							}
							for bindStackFrameArrayIndex := 0; bindStackFrameArrayIndex < bindStackFrameArraySize; bindStackFrameArrayIndex++ {
								if bindStackFrame.Fields[bindStackFrameFieldIndex].Field.ConstantPool {
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
									switch bindStackFrame.Fields[bindStackFrameFieldIndex].Field.Type {
									case typeMap.T_METHOD:
										if bindStackFrame.Fields[bindStackFrameFieldIndex].MethodRef != nil {
											*bindStackFrame.Fields[bindStackFrameFieldIndex].MethodRef = MethodRef(v32_)
										}
									case typeMap.T_FRAME_TYPE:
										if bindStackFrame.Fields[bindStackFrameFieldIndex].FrameTypeRef != nil {
											*bindStackFrame.Fields[bindStackFrameFieldIndex].FrameTypeRef = FrameTypeRef(v32_)
										}
									}
								} else {
									bindStackFrameFieldTypeID := bindStackFrame.Fields[bindStackFrameFieldIndex].Field.Type
									switch bindStackFrameFieldTypeID {
									case typeMap.T_STRING:
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
									case typeMap.T_INT:
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
										if bindStackFrame.Fields[bindStackFrameFieldIndex].uint32 != nil {
											*bindStackFrame.Fields[bindStackFrameFieldIndex].uint32 = v32_
										}
									case typeMap.T_LONG:
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
									case typeMap.T_BOOLEAN:
										if pos >= l {
											return 0, io.ErrUnexpectedEOF
										}
										b_ = data[pos]
										pos++
										// skipping
									case typeMap.T_FLOAT:
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
									default:
										bindStackFrameFieldType := typeMap.IDMap[bindStackFrame.Fields[bindStackFrameFieldIndex].Field.Type]
										if bindStackFrameFieldType == nil || len(bindStackFrameFieldType.Fields) == 0 {
											return 0, fmt.Errorf("unknown type %d", bindStackFrame.Fields[bindStackFrameFieldIndex].Field.Type)
										}
										bindStackFrameSkipObjects := 1
										if bindStackFrame.Fields[bindStackFrameFieldIndex].Field.Array {
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
											bindStackFrameSkipObjects = int(v32_)
										}
										for bindStackFrameSkipObjectIndex := 0; bindStackFrameSkipObjectIndex < bindStackFrameSkipObjects; bindStackFrameSkipObjectIndex++ {
											for bindStackFrameskipFieldIndex := 0; bindStackFrameskipFieldIndex < len(bindStackFrameFieldType.Fields); bindStackFrameskipFieldIndex++ {
												bindStackFrameSkipFieldType := bindStackFrameFieldType.Fields[bindStackFrameskipFieldIndex].Type
												if bindStackFrameFieldType.Fields[bindStackFrameskipFieldIndex].ConstantPool {
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
												} else if bindStackFrameSkipFieldType == typeMap.T_STRING {
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
												} else if bindStackFrameSkipFieldType == typeMap.T_INT {
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
												} else if bindStackFrameSkipFieldType == typeMap.T_FLOAT {
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
												} else if bindStackFrameSkipFieldType == typeMap.T_LONG {
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
												} else if bindStackFrameSkipFieldType == typeMap.T_BOOLEAN {
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
						if bind.Fields[bindFieldIndex].StackFrame != nil {
							*bind.Fields[bindFieldIndex].StackFrame = append(*bind.Fields[bindFieldIndex].StackFrame, bindStackFrame.Temp)
						}
					default:
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
		this.StackTrace[i] = bind.Temp
		this.IDMap[id] = uint32(i)
	}
	return pos, nil
}
