// Code generated by gen/main.go. DO NOT EDIT.

package types

import (
	"fmt"

	"github.com/grafana/jfr-parser/parser/types/def"
	"github.com/grafana/jfr-parser/util"
)

type BindActiveSetting struct {
	Temp   ActiveSetting
	Fields []BindFieldActiveSetting
}

type BindFieldActiveSetting struct {
	Field         *def.Field
	uint64        *uint64
	ThreadRef     *ThreadRef
	StackTraceRef *StackTraceRef
	string        *string
}

func NewBindActiveSetting(typ *def.Class, typeMap *def.TypeMap) *BindActiveSetting {
	res := new(BindActiveSetting)
	res.Fields = make([]BindFieldActiveSetting, 0, len(typ.Fields))
	for i := 0; i < len(typ.Fields); i++ {
		switch typ.Fields[i].Name {
		case "startTime":
			if typ.Fields[i].Equals(&def.Field{Name: "startTime", Type: typeMap.T_LONG, ConstantPool: false, Array: false}) {
				res.Fields = append(res.Fields, BindFieldActiveSetting{Field: &typ.Fields[i], uint64: &res.Temp.StartTime})
			} else {
				res.Fields = append(res.Fields, BindFieldActiveSetting{Field: &typ.Fields[i]}) // skip changed field
			}
		case "duration":
			if typ.Fields[i].Equals(&def.Field{Name: "duration", Type: typeMap.T_LONG, ConstantPool: false, Array: false}) {
				res.Fields = append(res.Fields, BindFieldActiveSetting{Field: &typ.Fields[i], uint64: &res.Temp.Duration})
			} else {
				res.Fields = append(res.Fields, BindFieldActiveSetting{Field: &typ.Fields[i]}) // skip changed field
			}
		case "eventThread":
			if typ.Fields[i].Equals(&def.Field{Name: "eventThread", Type: typeMap.T_THREAD, ConstantPool: true, Array: false}) {
				res.Fields = append(res.Fields, BindFieldActiveSetting{Field: &typ.Fields[i], ThreadRef: &res.Temp.EventThread})
			} else {
				res.Fields = append(res.Fields, BindFieldActiveSetting{Field: &typ.Fields[i]}) // skip changed field
			}
		case "stackTrace":
			if typ.Fields[i].Equals(&def.Field{Name: "stackTrace", Type: typeMap.T_STACK_TRACE, ConstantPool: true, Array: false}) {
				res.Fields = append(res.Fields, BindFieldActiveSetting{Field: &typ.Fields[i], StackTraceRef: &res.Temp.StackTrace})
			} else {
				res.Fields = append(res.Fields, BindFieldActiveSetting{Field: &typ.Fields[i]}) // skip changed field
			}
		case "id":
			if typ.Fields[i].Equals(&def.Field{Name: "id", Type: typeMap.T_LONG, ConstantPool: false, Array: false}) {
				res.Fields = append(res.Fields, BindFieldActiveSetting{Field: &typ.Fields[i], uint64: &res.Temp.Id})
			} else {
				res.Fields = append(res.Fields, BindFieldActiveSetting{Field: &typ.Fields[i]}) // skip changed field
			}
		case "name":
			if typ.Fields[i].Equals(&def.Field{Name: "name", Type: typeMap.T_STRING, ConstantPool: false, Array: false}) {
				res.Fields = append(res.Fields, BindFieldActiveSetting{Field: &typ.Fields[i], string: &res.Temp.Name})
			} else {
				res.Fields = append(res.Fields, BindFieldActiveSetting{Field: &typ.Fields[i]}) // skip changed field
			}
		case "value":
			if typ.Fields[i].Equals(&def.Field{Name: "value", Type: typeMap.T_STRING, ConstantPool: false, Array: false}) {
				res.Fields = append(res.Fields, BindFieldActiveSetting{Field: &typ.Fields[i], string: &res.Temp.Value})
			} else {
				res.Fields = append(res.Fields, BindFieldActiveSetting{Field: &typ.Fields[i]}) // skip changed field
			}
		default:
			res.Fields = append(res.Fields, BindFieldActiveSetting{Field: &typ.Fields[i]}) // skip unknown new field
		}
	}
	return res
}

type ActiveSetting struct {
	StartTime   uint64
	Duration    uint64
	EventThread ThreadRef
	StackTrace  StackTraceRef
	Id          uint64
	Name        string
	Value       string
}

func (this *ActiveSetting) Parse(data []byte, bind *BindActiveSetting, typeMap *def.TypeMap) (pos int, err error) {
	for bindFieldIndex := 0; bindFieldIndex < len(bind.Fields); bindFieldIndex++ {
		bindArraySize := 1
		if bind.Fields[bindFieldIndex].Field.Array {

			v32_, err := util.ParseVarInt(data, &pos)
			if err != nil {
				return 0, err
			}
			_ = v32_

			bindArraySize = int(v32_)
		}
		for bindArrayIndex := 0; bindArrayIndex < bindArraySize; bindArrayIndex++ {
			if bind.Fields[bindFieldIndex].Field.ConstantPool {

				v32_, err := util.ParseVarInt(data, &pos)
				if err != nil {
					return 0, err
				}
				_ = v32_

				switch bind.Fields[bindFieldIndex].Field.Type {
				case typeMap.T_THREAD:
					if bind.Fields[bindFieldIndex].ThreadRef != nil {
						*bind.Fields[bindFieldIndex].ThreadRef = ThreadRef(v32_)
					}
				case typeMap.T_STACK_TRACE:
					if bind.Fields[bindFieldIndex].StackTraceRef != nil {
						*bind.Fields[bindFieldIndex].StackTraceRef = StackTraceRef(v32_)
					}
				}
			} else {
				bindFieldTypeID := bind.Fields[bindFieldIndex].Field.Type
				switch bindFieldTypeID {
				case typeMap.T_STRING:

					s_, err := util.ParseString(data, &pos)
					if err != nil {
						return 0, err
					}
					_ = s_

					if bind.Fields[bindFieldIndex].string != nil {
						*bind.Fields[bindFieldIndex].string = s_
					}
				case typeMap.T_INT:

					v32_, err := util.ParseVarInt(data, &pos)
					if err != nil {
						return 0, err
					}
					_ = v32_

					// skipping
				case typeMap.T_LONG:

					v64_, err := util.ParseVarLong(data, &pos)
					if err != nil {
						return 0, err
					}
					_ = v64_

					if bind.Fields[bindFieldIndex].uint64 != nil {
						*bind.Fields[bindFieldIndex].uint64 = v64_
					}
				case typeMap.T_BOOLEAN:

					b_, err := util.ParseByte(data, &pos)
					if err != nil {
						return 0, err
					}
					_ = b_

					// skipping
				case typeMap.T_FLOAT:

					v32_, err := util.ParseVarInt(data, &pos)
					if err != nil {
						return 0, err
					}
					_ = v32_

					// skipping
				default:
					bindFieldType := typeMap.IDMap[bind.Fields[bindFieldIndex].Field.Type]
					if bindFieldType == nil || len(bindFieldType.Fields) == 0 {
						return 0, fmt.Errorf("unknown type %d", bind.Fields[bindFieldIndex].Field.Type)
					}
					bindSkipObjects := 1
					if bind.Fields[bindFieldIndex].Field.Array {

						v32_, err := util.ParseVarInt(data, &pos)
						if err != nil {
							return 0, err
						}
						_ = v32_

						bindSkipObjects = int(v32_)
					}
					for bindSkipObjectIndex := 0; bindSkipObjectIndex < bindSkipObjects; bindSkipObjectIndex++ {
						for bindskipFieldIndex := 0; bindskipFieldIndex < len(bindFieldType.Fields); bindskipFieldIndex++ {
							bindSkipFieldType := bindFieldType.Fields[bindskipFieldIndex].Type
							if bindFieldType.Fields[bindskipFieldIndex].ConstantPool {

								v32_, err := util.ParseVarInt(data, &pos)
								if err != nil {
									return 0, err
								}
								_ = v32_

							} else if bindSkipFieldType == typeMap.T_STRING {

								s_, err := util.ParseString(data, &pos)
								if err != nil {
									return 0, err
								}
								_ = s_

							} else if bindSkipFieldType == typeMap.T_INT {

								v32_, err := util.ParseVarInt(data, &pos)
								if err != nil {
									return 0, err
								}
								_ = v32_

							} else if bindSkipFieldType == typeMap.T_FLOAT {

								v32_, err := util.ParseVarInt(data, &pos)
								if err != nil {
									return 0, err
								}
								_ = v32_

							} else if bindSkipFieldType == typeMap.T_LONG {

								v64_, err := util.ParseVarLong(data, &pos)
								if err != nil {
									return 0, err
								}
								_ = v64_

							} else if bindSkipFieldType == typeMap.T_BOOLEAN {

								b_, err := util.ParseByte(data, &pos)
								if err != nil {
									return 0, err
								}
								_ = b_

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
