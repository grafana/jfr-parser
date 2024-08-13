// Code generated by gen/main.go. DO NOT EDIT.

package types

import (
	"fmt"

	"github.com/grafana/jfr-parser/parser/types/def"
	"github.com/grafana/jfr-parser/util"
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
	res.Fields = make([]BindFieldMethod, 0, len(typ.Fields))
	for i := 0; i < len(typ.Fields); i++ {
		switch typ.Fields[i].Name {
		case "type":
			if typ.Fields[i].Equals(&def.Field{Name: "type", Type: typeMap.T_CLASS, ConstantPool: true, Array: false}) {
				res.Fields = append(res.Fields, BindFieldMethod{Field: &typ.Fields[i], ClassRef: &res.Temp.Type})
			} else {
				res.Fields = append(res.Fields, BindFieldMethod{Field: &typ.Fields[i]}) // skip changed field
			}
		case "name":
			if typ.Fields[i].Equals(&def.Field{Name: "name", Type: typeMap.T_SYMBOL, ConstantPool: true, Array: false}) {
				res.Fields = append(res.Fields, BindFieldMethod{Field: &typ.Fields[i], SymbolRef: &res.Temp.Name})
			} else {
				res.Fields = append(res.Fields, BindFieldMethod{Field: &typ.Fields[i]}) // skip changed field
			}
		case "descriptor":
			res.Fields = append(res.Fields, BindFieldMethod{Field: &typ.Fields[i]}) // skip to save mem
		case "modifiers":
			res.Fields = append(res.Fields, BindFieldMethod{Field: &typ.Fields[i]}) // skip to save mem
		case "hidden":
			res.Fields = append(res.Fields, BindFieldMethod{Field: &typ.Fields[i]}) // skip to save mem
		default:
			res.Fields = append(res.Fields, BindFieldMethod{Field: &typ.Fields[i]}) // skip unknown new field
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
	Type ClassRef
	Name SymbolRef
	// skip descriptor
	// skip modifiers
	// skip hidden
}

func (this *MethodList) Parse(data []byte, bind *BindMethod, typeMap *def.TypeMap) (pos int, err error) {

	v32_, err := util.ParseVarInt(data, &pos)
	if err != nil {
		return 0, err
	}
	_ = v32_

	n := int(v32_)
	this.IDMap = NewIDMap[MethodRef](n)
	this.Method = make([]Method, n)
	for i := 0; i < n; i++ {

		v32_, err := util.ParseVarInt(data, &pos)
		if err != nil {
			return 0, err
		}
		_ = v32_

		id := MethodRef(v32_)
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
					switch bindFieldTypeID {
					case typeMap.T_STRING:

						s_, err := util.ParseString(data, &pos)
						if err != nil {
							return 0, err
						}
						_ = s_

						// skipping
					case typeMap.T_INT:

						v32_, err := util.ParseVarInt(data, &pos)
						if err != nil {
							return 0, err
						}
						_ = v32_

						if bind.Fields[bindFieldIndex].uint32 != nil {
							*bind.Fields[bindFieldIndex].uint32 = v32_
						}
					case typeMap.T_LONG:

						v64_, err := util.ParseVarLong(data, &pos)
						if err != nil {
							return 0, err
						}
						_ = v64_

						// skipping
					case typeMap.T_BOOLEAN:

						b_, err := util.ParseByte(data, &pos)
						if err != nil {
							return 0, err
						}
						_ = b_

						if bind.Fields[bindFieldIndex].bool != nil {
							*bind.Fields[bindFieldIndex].bool = b_ != 0
						}
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
		this.Method[i] = bind.Temp
		this.IDMap.Set(id, i)
	}
	return pos, nil
}
