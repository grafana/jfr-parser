package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/pyroscope-io/jfr-parser/parser/types/def"
)

func main() {
	write("types/frametype.go", generate(&Type_jdk_types_FrameType, options{
		cpool: true,
	}))
	write("types/threadstate.go", generate(&Type_jdk_types_ThreadState, options{
		cpool: true,
	}))
	write("types/thread.go", generate(&Type_java_lang_Thread, options{
		cpool: true,
	}))
	write("types/class.go", generate(&Type_java_lang_Class, options{
		skipFields: []string{"classLoader", "modifiers"},
		cpool:      true,
	}))
	write("types/method.go", generate(&Type_jdk_types_Method, options{
		skipFields: []string{"hidden", "descriptor", "modifiers"},
		cpool:      true,
		sortedIDs:  true,
		Scratch:    true,
	}))
	write("types/package.go", generate(&Type_jdk_types_Package, options{
		cpool: true,
	}))
	write("types/symbol.go", generate(&Type_jdk_types_Symbol, options{
		cpool: true,
	}))
	write("types/loglevel.go", generate(&Type_profiler_types_LogLevel, options{
		cpool: true,
	}))
	write("types/stacktrace.go", generate(&Type_jdk_types_StackTrace, options{
		cpool: true,
	}))

	write("types/generic.go", generateGeneric())

	write("types/active_recording.go", generate(&Type_jdk_ActiveRecording, options{}))
	write("types/active_settings.go", generate(&Type_jdk_ActiveSetting, options{}))
	write("types/os_info.go", generate(&Type_jdk_OSInformation, options{}))
	write("types/jvm_info.go", generate(&Type_jdk_JVMInformation, options{}))
	write("types/initial_system_property.go", generate(&Type_jdk_InitialSystemProperty, options{}))
	write("types/native_libraries.go", generate(&Type_jdk_NativeLibrary, options{}))
	//todo cpuload

	write("types/execution_sample.go", generate(&Type_jdk_ExecutionSample, options{
		optionalField: "contextId",
	}))
	write("types/allocation_in_new_tlab.go", generate(&Type_jdk_ObjectAllocationInNewTLAB, options{
		optionalField: "contextId",
	}))
	write("types/allocation_outside_tlab.go", generate(&Type_jdk_ObjectAllocationOutsideTLAB, options{
		optionalField: "contextId",
	}))
	write("types/monitor_enter.go", generate(&Type_jdk_JavaMonitorEnter, options{
		optionalField: "contextId",
	}))
	write("types/thread_park.go", generate(&Type_jdk_ThreadPark, options{
		optionalField: "contextId",
	}))
	write("types/live_object.go", generate(&Type_profiler_LiveObject, options{}))

	write("types/log.go", generate(&Type_profiler_Log, options{}))
	write("types/cpu_load.go", generate(&Type_jdk_CPULoad, options{}))

}

func write(dst, s string) {
	err := os.WriteFile("parser/"+dst, []byte(s), 0666)
	if err != nil {
		panic(err)
	}
}

type options struct {
	skipFields []string

	optionalField string
	cpool         bool
	sortedIDs     bool
	Scratch       bool
}

func (o *options) skipField(field def.Field) bool {
	for _, skipField := range o.skipFields {
		if skipField == field.Name {
			return true
		}
	}
	return false
}

func TypeForCPoolID(ID def.TypeID) *def.Class {
	switch ID {
	case def.T_FRAME_TYPE:
		return &Type_jdk_types_FrameType
	case def.T_THREAD_STATE:
		return &Type_jdk_types_ThreadState
	case def.T_THREAD:
		return &Type_java_lang_Thread
	case def.T_CLASS:
		return &Type_java_lang_Class
	case def.T_METHOD:
		return &Type_jdk_types_Method
	case def.T_PACKAGE:
		return &Type_jdk_types_Package
	case def.T_SYMBOL:
		return &Type_jdk_types_Symbol
	case def.T_LOG_LEVEL:
		return &Type_profiler_types_LogLevel
	case def.T_STACK_TRACE:
		return &Type_jdk_types_StackTrace
	default:
		panic("unknown type " + def.TypeID2Sym(ID))
	}
}

func generate(typ *def.Class, opt options) string {
	res := ""

	if opt.cpool {
		res += fmt.Sprintf("type %sRef uint32\n", name(typ))
		res += fmt.Sprintf("type %sList struct {\n", name(typ))
		if opt.sortedIDs {
			res += fmt.Sprintf("	IDMap IDMap[%sRef]\n", name(typ))
		} else {
			res += fmt.Sprintf("	IDMap map[%sRef]uint32\n", name(typ))
		}
		res += fmt.Sprintf("	%s []%s\n", name(typ), name(typ))
		res += fmt.Sprintf("}\n\n")
	}

	res += fmt.Sprintf("var ExpectedMeta%s = &def.Class{\n", name(typ))
	res += fmt.Sprintf("	Name: \"%s\",\n", typ.Name)
	res += fmt.Sprintf("	ID: def.%s,\n", def.TypeID2Sym(typ.ID))
	res += fmt.Sprintf("	Fields: []def.Field{\n")
	for _, field := range typ.Fields {
		res += fmt.Sprintf("		{\n")
		res += fmt.Sprintf("			Name: \"%s\",\n", field.Name)
		res += fmt.Sprintf("			Type: def.%s,\n", def.TypeID2Sym(field.Type))
		res += fmt.Sprintf("			ConstantPool: %t,\n", field.ConstantPool)
		res += fmt.Sprintf("			Array: %t,\n", field.Array)
		res += fmt.Sprintf("		},\n")
	}
	res += fmt.Sprintf("	},\n")
	res += fmt.Sprintf("}\n\n")
	res += fmt.Sprintf("type %s struct {\n", name(typ))
	for _, field := range typ.Fields {
		if opt.skipField(field) {
			continue
		}
		if field.Array {
			if field.ConstantPool {
				panic("cp array not implemented")
			} else {
				res += fmt.Sprintf("	%s []%s\n", capitalize(field.Name), goTypeName(field))
			}
		} else {
			res += fmt.Sprintf("	%s %s\n", capitalize(field.Name), goTypeName(field))
		}
	}
	if opt.Scratch {
		res += fmt.Sprintf("	Scratch string\n")
	}
	res += fmt.Sprintf("}\n\n")
	res += fmt.Sprintf("\n")

	var receiver string
	if opt.cpool {
		optinoal := ""
		if typ.ID == def.T_STACK_TRACE {
			optinoal = " stackFrameType *def.Class,"
		}
		res += fmt.Sprintf("func (this *%sList) Parse(data []byte, typ *def.Class,%s typeMap map[def.TypeID]*def.Class) (pos int, err error) {\n", name(typ), optinoal)
		receiver = fmt.Sprintf("this.%s[i]", name(typ))
	} else {
		receiver = fmt.Sprintf("this")
		optional := ""
		if opt.optionalField != "" {
			optional = ", have" + opt.optionalField + " bool"
		}
		res += fmt.Sprintf("func (this *%s) Parse(data []byte, typ *def.Class, typeMap map[def.TypeID]*def.Class%s) (pos int, err error) {\n", name(typ), optional)
	}
	_ = receiver

	res += "	var (\n"
	res += "		v64_ uint64\n"
	res += "		v32_ uint32\n"
	res += "		s_   string\n"
	res += "		b_   byte\n"
	res += "		shift = uint(0)\n"
	res += "		l = len(data)\n"
	res += "	)\n"
	res += "	_ = v64_\n"
	res += "	_ = v32_\n"
	res += "	_ = s_\n"

	res += pad(1) + fmt.Sprintf("nFields := len(ExpectedMeta%s.Fields)\n", name(typ))
	if opt.optionalField != "" {
		res += pad(1) + fmt.Sprintf("if !have%s {\n", opt.optionalField)
		res += pad(1) + fmt.Sprintf("	nFields -= 1\n")
		res += pad(1) + fmt.Sprintf("}\n")
	}
	res += pad(1) + fmt.Sprintf("skipFields := typ.Fields[nFields:]\n")
	if typ.ID == def.T_STACK_TRACE {
		res += pad(1) + fmt.Sprintf("stackFrameSkipFields := stackFrameType.Fields[len(def.TypeStackFrame.Fields):]\n")
	}

	depth := 2
	if opt.cpool {
		res += emitReadI32(1)
		res += pad(1) + "n := int(v32_)\n"
		if opt.sortedIDs {
			res += pad(1) + fmt.Sprintf("this.IDMap = NewIDMap[%sRef](n)\n", name(typ))
		} else {
			res += pad(1) + fmt.Sprintf("this.IDMap = make(map[%sRef]uint32, n)\n", name(typ))
		}
		res += pad(1) + fmt.Sprintf("this.%s = make([]%s, n)\n", name(typ), name(typ))
		res += "	for i := 0; i < n; i++ {\n"
	} else {
		depth = 1
	}

	if opt.cpool {
		res += emitReadI32(depth)
		res += pad(depth) + fmt.Sprintf("id := %sRef(v32_)\n", name(typ))
	}
	for _, field := range typ.Fields {
		if field.Name == opt.optionalField {
			res += fmt.Sprintf("		if have%s {\n", opt.optionalField)
		}
		if field.ConstantPool {
			if field.Array {
				panic("TODO " + field.String())
			}
			res += emitReadI32(depth)
			if opt.skipField(field) {
				res += fmt.Sprintf("\t\t// skipped %s\n", field.Name)
				continue
			}
			res += fmt.Sprintf("		%s.%s = %sRef(v32_)\n", receiver, capitalize(field.Name), name(TypeForCPoolID(field.Type)))
		} else {
			if field.Array {
				res += emitReadI32(depth)
				res += pad(depth) + "m := int(v32_)\n"
				res += pad(depth) + fmt.Sprintf("%s.%s = make([]%s, m)\n", receiver, capitalize(field.Name), goTypeName(field))
				res += pad(depth) + "for j := 0; j < m; j++ {\n"
				//
				switch field.Type {
				case def.T_STACK_FRAME:
					st := &Type_jdk_types_StackFrame
					for _, sf := range st.Fields {
						if sf.ConstantPool {
							res += emitReadI32(depth + 1)
							res += fmt.Sprintf("			%s.%s[j].%s = %sRef(v32_)\n", receiver, capitalize(field.Name), capitalize(sf.Name), name(TypeForCPoolID(sf.Type)))
						} else {
							switch sf.Type {
							case def.T_STRING:
								res += emitString(depth + 1)
								res += fmt.Sprintf("			%s.%s[j].%s  = s_\n", receiver, capitalize(field.Name), capitalize(sf.Name))
							case def.T_LONG:
								res += emitReadU64(depth + 1)
								res += fmt.Sprintf("			%s.%s[j].%s = v64_\n", receiver, capitalize(field.Name), capitalize(sf.Name))
							case def.T_INT:
								res += emitReadI32(depth + 1)
								res += fmt.Sprintf("			%s.%s[j].%s = v32_\n", receiver, capitalize(field.Name), capitalize(sf.Name))
							case def.T_BOOLEAN:
								res += emitReadByte(depth + 1)
								res += fmt.Sprintf("			%s.%s[j].%s = b_ == 0\n", receiver, capitalize(field.Name), capitalize(sf.Name))
							default:
								panic("TODO " + field.String())
							}
						}
						res += emitSkipFields("stackFrameSkipFields", "typeMap", 3)

					}
				default:
					panic("TODO " + field.String())
				}
				res += "		}\n"
			} else {
				switch field.Type {
				case def.T_STRING:

					res += emitString(depth)

					if opt.skipField(field) {
						res += pad(depth) + fmt.Sprintf("// skipped %s\n", field.Name)
						continue
					}
					res += pad(depth) + fmt.Sprintf("%s.%s  = s_\n", receiver, capitalize(field.Name))

				case def.T_LONG:
					res += emitReadU64(depth)

					if opt.skipField(field) {
						res += pad(depth) + fmt.Sprintf("// skipped %s\n", field.Name)
						continue
					}
					res += pad(depth) + fmt.Sprintf("%s.%s = v64_\n", receiver, capitalize(field.Name))
				case def.T_INT:
					res += emitReadI32(depth)

					if opt.skipField(field) {
						res += pad(depth) + fmt.Sprintf("// skipped %s\n", field.Name)
						continue
					}
					res += pad(depth) + fmt.Sprintf("%s.%s = v32_\n", receiver, capitalize(field.Name))
				case def.T_FLOAT:
					res += emitReadI32(depth)

					if opt.skipField(field) {
						res += pad(depth) + fmt.Sprintf("// skipped %s\n", field.Name)
						continue
					}
					res += pad(depth) + fmt.Sprintf("%s.%s = *(*float32)(unsafe.Pointer(&v32_))\n", receiver, capitalize(field.Name))
				case def.T_BOOLEAN:
					res += emitReadByte(2)

					if opt.skipField(field) {
						res += pad(depth) + fmt.Sprintf("// skipped %s\n", field.Name)
						continue
					}
					res += pad(depth) + fmt.Sprintf("this.%s[i].%s = b_ == 0\n", name(typ), capitalize(field.Name))
				default:
					panic("TODO " + field.String())
				}
			}

		}
		if field.Name == opt.optionalField {
			res += pad(depth) + fmt.Sprintf("}\n")
		}
	}
	res += emitSkipFields("skipFields", "typeMap", 2)

	if opt.cpool {
		if opt.sortedIDs {
			res += pad(depth) + "this.IDMap.Set(id, i)\n"
		} else {
			res += pad(depth) + "this.IDMap[id] = uint32(i)\n"
		}
		res += pad(1) + "}\n"
	}
	res += "	return pos, nil\n"
	res += fmt.Sprintf("}\n")

	imports := "package types\n"
	imports += "\n"

	imports += "import (\n\t\"fmt\"\n\t\"io\"\n\t\"unsafe\"\n\t\"github.com/pyroscope-io/jfr-parser/parser/types/def\"\n\n)"

	imports += "\n"
	res = imports + res

	fmt.Println("types2.ExpectedMeta" + name(typ) + ",")
	//fmt.Println(res)
	return res
}

func emitSkipFields(skipFieldsSliceName string, meta string, depth int) string {
	res := pad(depth) + ""
	res += pad(depth) + "\n\n"
	res += pad(depth) + "// skipping added fields \n"
	res += pad(depth) + fmt.Sprintf("for skipFI := range %s {\n", skipFieldsSliceName)
	res += pad(depth) + "	nSkip := int(1)\n"
	res += pad(depth) + fmt.Sprintf("	if %s[skipFI].Array {\n", skipFieldsSliceName)
	res += emitReadI32(depth + 2)
	res += pad(depth) + "		nSkip = int(v32_)\n"
	res += pad(depth) + "	}\n"
	res += pad(depth) + "	for iSkip := 0; iSkip < nSkip; iSkip++ {\n"
	res += pad(depth) + fmt.Sprintf("		if %s[skipFI].ConstantPool {\n", skipFieldsSliceName)
	res += emitReadI32(depth + 3)
	res += pad(depth) + "		} else {\n"
	res += pad(depth) + fmt.Sprintf("			switch %s[skipFI].Type {\n", skipFieldsSliceName)
	res += pad(depth) + "			case def.T_STRING:\n"
	res += emitString(depth + 4)
	res += pad(depth) + "			case def.T_LONG:\n"
	res += emitReadU64(depth + 4)
	res += pad(depth) + "			case def.T_INT:\n"
	res += emitReadI32(depth + 4)
	res += pad(depth) + "			case def.T_FLOAT:\n"
	res += emitReadI32(depth + 4)
	res += pad(depth) + "			case def.T_BOOLEAN:\n"
	res += emitReadByte(depth + 4)
	res += pad(depth) + "			default:\n"
	res += pad(depth) + fmt.Sprintf("				gt := %s[%s[skipFI].Type]\n", meta, skipFieldsSliceName)
	res += pad(depth) + fmt.Sprintf("				if gt == nil {\n")
	res += pad(depth) + fmt.Sprintf("					return 0, fmt.Errorf(\"unknown type %%d\", %s[skipFI].Type)\n", skipFieldsSliceName)
	res += pad(depth) + fmt.Sprintf("				}\n")

	res += pad(depth) + "				for gti := 0; gti < len(gt.Fields); gti++ {\n"
	res += pad(depth) + "					if gt.Fields[gti].Array {\n"
	res += pad(depth) + "						return 0, fmt.Errorf(\"two dimentional array not supported\")"
	res += pad(depth) + "					}\n"
	res += pad(depth) + "					if gt.Fields[gti].ConstantPool {\n"
	res += emitReadI32(depth + 6)
	res += pad(depth) + "					} else {\n"

	res += pad(depth) + "						switch gt.Fields[gti].Type {\n"
	res += pad(depth) + "						case def.T_STRING:\n"
	res += emitString(depth + 7)
	res += pad(depth) + "						case def.T_LONG:\n"
	res += emitReadU64(depth + 7)
	res += pad(depth) + "						case def.T_INT:\n"
	res += emitReadI32(depth + 7)
	res += pad(depth) + "						case def.T_FLOAT:\n"
	res += emitReadI32(depth + 7)
	res += pad(depth) + "						case def.T_BOOLEAN:\n"
	res += emitReadByte(depth + 7)
	res += pad(depth) + "						default:\n"
	res += pad(depth) + "							return 0, fmt.Errorf(\"unknown type %d\", gt.Fields[gti].Type)\n"
	res += pad(depth) + "						}\n"
	res += pad(depth) + "					}\n"
	res += pad(depth) + "				}\n"
	res += pad(depth) + "			}\n"
	res += pad(depth) + "		}\n"
	res += pad(depth) + "	}\n"
	res += pad(depth) + "}\n"
	return res
}

func emitString(depth int) string {
	res := pad(depth) + "s_ = \"\"\n"
	res += pad(depth) + "if pos >= l {\n"
	res += pad(depth) + "	return 0, io.ErrUnexpectedEOF\n"
	res += pad(depth) + "}\n"
	res += pad(depth) + "b_ = data[pos]\n"
	res += pad(depth) + "pos++\n"
	res += pad(depth) + "switch b_ {\n"
	res += pad(depth) + "case 0:\n"
	res += pad(depth) + "case 1:\n"
	res += pad(depth) + "	break\n"
	res += pad(depth) + "case 3:\n"
	res += emitReadI32(depth + 1)
	res += pad(depth) + "	if pos+int(v32_) > l {\n"
	res += pad(depth) + "		return 0, io.ErrUnexpectedEOF\n"
	res += pad(depth) + "	}\n"

	res += pad(depth) + "	bs := data[pos : pos+int(v32_)]\n"
	res += pad(depth) + fmt.Sprintf("	s_ = *(*string)(unsafe.Pointer(&bs))\n")

	res += pad(depth) + "	pos += int(v32_)\n"
	res += pad(depth) + "default:\n"
	res += pad(depth) + "	return 0, fmt.Errorf(\"unknown string type %d at %d\", b_, pos)\n"
	res += pad(depth) + "}\n"
	return res
}

func emitReadByte(depth int) string {
	code := ""
	code += pad(depth) + "if pos >= l {\n"
	code += pad(depth) + "	return 0, io.ErrUnexpectedEOF\n"
	code += pad(depth) + "}\n"
	code += pad(depth) + "b_ = data[pos]\n"
	code += pad(depth) + "pos++\n"
	return code

}
func emitReadI32(depth int) string {
	code := ""
	code += pad(depth) + "v32_ = uint32(0)\n"
	code += pad(depth) + "for shift = uint(0); ; shift += 7 {\n"
	code += pad(depth) + "	if shift >= 32 {\n"
	code += pad(depth) + "		return 0, def.ErrIntOverflow\n"
	code += pad(depth) + "	}\n"
	code += pad(depth) + "	if pos >= l {\n"
	code += pad(depth) + "		return 0, io.ErrUnexpectedEOF\n"
	code += pad(depth) + "	}\n"
	code += pad(depth) + "	b_ = data[pos]\n"
	code += pad(depth) + "	pos++\n"
	code += pad(depth) + "	v32_ |= uint32(b_&0x7F) << shift\n"
	code += pad(depth) + "	if b_ < 0x80 {\n"
	code += pad(depth) + "		break\n"
	code += pad(depth) + "	}\n"
	code += pad(depth) + "}\n"

	return code
}

func emitReadU64(depth int) string {
	code := ""

	code += pad(depth) + "v64_ = 0 \n"
	code += pad(depth) + "for shift = uint(0); shift <= 56 ; shift += 7 {\n"
	code += pad(depth) + "	if pos >= l {\n"
	code += pad(depth) + "		return 0, io.ErrUnexpectedEOF\n"
	code += pad(depth) + "	}\n"
	code += pad(depth) + "	b_ = data[pos]\n"
	code += pad(depth) + "	pos++\n"
	code += pad(depth) + "	if shift == 56{\n"
	code += pad(depth) + "		v64_ |= uint64(b_&0xFF) << shift\n"
	code += pad(depth) + "		break\n"
	code += pad(depth) + "	} else {\n"
	code += pad(depth) + "		v64_ |= uint64(b_&0x7F) << shift\n"
	code += pad(depth) + "		if b_ < 0x80 {\n"
	code += pad(depth) + "			break\n"
	code += pad(depth) + "		}\n"
	code += pad(depth) + "	}\n"
	code += pad(depth) + "}\n"

	return code
}

func goTypeName(field def.Field) string {
	if field.ConstantPool {
		return name(TypeForCPoolID(field.Type)) + "Ref"
	}
	switch field.Type {
	case def.T_STRING:
		return "string"
	case def.T_LONG:
		return "uint64"
	case def.T_INT:
		return "uint32"
	case def.T_FLOAT:
		return "float32"
	case def.T_BOOLEAN:
		return "bool"
	case def.T_STACK_FRAME:
		return "StackFrame"
	default:
		panic("TODO " + field.String())
	}
}

func name(typ *def.Class) string {
	fs := strings.Split(typ.Name, ".")
	s := fs[len(fs)-1]
	return capitalize(s)
}

func capitalize(s string) string {
	return strings.ToUpper(s[:1]) + s[1:]
}

func pad(n int) string {
	res := ""
	for i := 0; i < n; i++ {
		res += "\t"
	}
	return res
}

func generateGeneric() string {
	res := "package types\n\nimport (\n\t\"fmt\"\n\t\"io\"\n\t\"unsafe\"\n\n\t\"github.com/pyroscope-io/jfr-parser/parser/types/def\"\n)\n\n"

	res += "func Skip(data []byte, typ *def.Class, typeMap map[def.TypeID]*def.Class, cpool bool) (pos int, err error) {\n"
	res += "	var (\n"
	res += "		v64_ uint64\n"
	res += "		v32_ uint32\n"
	res += "		s_   string\n"
	res += "		b_   byte\n"
	res += "		shift = uint(0)\n"
	res += "		l = len(data)\n"
	res += "	)\n"
	res += "	_ = v64_\n"
	res += "	_ = v32_\n"
	res += "	_ = s_\n"
	res += "	nObj := 1\n"
	res += "	if cpool {\n"
	res += emitReadI32(2)
	res += "		nObj = int(v32_)\n"
	res += "	}\n"
	res += "	for i := 0; i < nObj; i++ {\n"
	res += "		if cpool {\n"
	res += emitReadI32(3) // id
	res += "		}\n"

	res += emitSkipFields("typ.Fields", "typeMap", 2)
	res += "	}\n"
	res += "	return pos, nil\n"
	res += "}\n\n"
	return res
}
