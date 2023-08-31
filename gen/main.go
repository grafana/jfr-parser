package main

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/grafana/jfr-parser/parser/types/def"
)

func main() {
	write("types/frametype.go", generate(&Type_jdk_types_FrameType, options{
		cpool:     true,
		sortedIDs: false,
	}))
	write("types/stackframe.go", generate(&Type_jdk_types_StackFrame, options{
		skipFields: []string{
			"lineNumber", "bytecodeIndex", "type",
		},
		cpool: false,
	}))
	write("types/threadstate.go", generate(&Type_jdk_types_ThreadState, options{
		cpool: true,
	}))
	write("types/thread.go", generate(&Type_java_lang_Thread, options{
		cpool: true,
	}))
	write("types/class.go", generate(&Type_java_lang_Class, options{
		skipFields: []string{
			"classLoader",
			"package",
			"modifiers",
		},
		cpool: true,
	}))
	write("types/classloader.go", generate(&Type_jdk_types_ClassLoader, options{
		cpool: true,
	}))
	write("types/method.go", generate(&Type_jdk_types_Method, options{
		cpool:     true,
		sortedIDs: true,
		Scratch:   true,
		skipFields: []string{
			"hidden",
			"descriptor",
			"modifiers",
		},
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

	write("types/active_settings.go", generate(&Type_jdk_ActiveSetting, options{}))

	write("types/execution_sample.go", generate(&Type_jdk_ExecutionSample, options{}))
	write("types/allocation_in_new_tlab.go", generate(&Type_jdk_ObjectAllocationInNewTLAB, options{}))
	write("types/allocation_outside_tlab.go", generate(&Type_jdk_ObjectAllocationOutsideTLAB, options{}))
	write("types/monitor_enter.go", generate(&Type_jdk_JavaMonitorEnter, options{}))
	write("types/thread_park.go", generate(&Type_jdk_ThreadPark, options{}))
	write("types/live_object.go", generate(&Type_profiler_LiveObject, options{}))
	write("types/skipper.go", generate(&def.Class{
		Name:   "SkipConstantPool",
		ID:     0,
		Fields: []def.Field{},
	}, options{
		cpool:         true,
		doNotKeepData: true,
	}))

}

func write(dst, s string) {
	err := os.WriteFile("parser/"+dst, []byte(s), 0666)
	if err != nil {
		panic(err)
	}
}

type options struct {
	cpool         bool
	sortedIDs     bool
	Scratch       bool
	doNotKeepData bool
	skipFields    []string //todo make skip fields runtime option, but still saving memory - explode struct to fields
}

func TypeForCPoolID(ID def.TypeID) *def.Class {
	switch ID {
	case T_FRAME_TYPE:
		return &Type_jdk_types_FrameType
	case T_THREAD_STATE:
		return &Type_jdk_types_ThreadState
	case T_THREAD:
		return &Type_java_lang_Thread
	case T_CLASS:
		return &Type_java_lang_Class
	case T_METHOD:
		return &Type_jdk_types_Method
	case T_PACKAGE:
		return &Type_jdk_types_Package
	case T_SYMBOL:
		return &Type_jdk_types_Symbol
	case T_LOG_LEVEL:
		return &Type_profiler_types_LogLevel
	case T_STACK_TRACE:
		return &Type_jdk_types_StackTrace
	case T_CLASS_LOADER:
		return &Type_jdk_types_ClassLoader
	case T_STACK_FRAME:
		return &Type_jdk_types_StackFrame
	default:
		panic("unknown type " + TypeID2Sym(ID))
	}
}

func generate(typ *def.Class, opt options) string {
	res := ""
	res += generateBinding(typ, opt)

	if opt.cpool {
		res += fmt.Sprintf("type %s uint32\n", refName(typ))
		res += fmt.Sprintf("type %s struct {\n", listName(typ))
		if opt.doNotKeepData {

		} else {
			if opt.sortedIDs {
				res += fmt.Sprintf("	IDMap IDMap[%s]\n", refName(typ))
			} else {
				res += fmt.Sprintf("	IDMap map[%s]uint32\n", refName(typ))
			}
			res += fmt.Sprintf("	%s []%s\n", name(typ), name(typ))
		}
		res += fmt.Sprintf("}\n\n")
	}

	res += fmt.Sprintf("type %s struct {\n", name(typ))
	for _, field := range typ.Fields {
		if slices.Contains(opt.skipFields, field.Name) {
			res += fmt.Sprintf("	// skip %s\n", field.Name)
		} else {
			if field.Array {
				res += fmt.Sprintf("	%s []%s\n", capitalize(field.Name), goTypeName(field))
			} else {
				res += fmt.Sprintf("	%s %s\n", capitalize(field.Name), goTypeName(field))
			}
		}
	}
	if opt.Scratch {
		res += fmt.Sprintf("	Scratch []byte\n")
	}
	res += fmt.Sprintf("}\n\n")
	res += fmt.Sprintf("\n")

	var receiver string
	if opt.cpool {
		extraBinds := ""
		compoudBindings := getNonBasicFields(typ)
		for _, binding := range compoudBindings {
			extraBinds += fmt.Sprintf(", bind%s *%s", name(TypeForCPoolID(binding.Type)), bindName(TypeForCPoolID(binding.Type)))
		}
		res += fmt.Sprintf("func (this *%sList) Parse(data []byte, bind *%s %s , typeMap *def.TypeMap) (pos int, err error) {\n", name(typ), bindName(typ), extraBinds)
		receiver = fmt.Sprintf("this.%s[i]", name(typ))
	} else {
		receiver = fmt.Sprintf("this")

		res += fmt.Sprintf("func (this *%s) Parse(data []byte, bind *%s, typeMap *def.TypeMap) (pos int, err error) {\n", name(typ), bindName(typ))
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

	depth := 2
	if opt.cpool {
		res += emitReadI32(1)
		res += pad(1) + "n := int(v32_)\n"
		if opt.doNotKeepData {

		} else {
			if opt.sortedIDs {
				res += pad(1) + fmt.Sprintf("this.IDMap = NewIDMap[%s](n)\n", refName(typ))
			} else {
				res += pad(1) + fmt.Sprintf("this.IDMap = make(map[%s]uint32, n)\n", refName(typ))
			}
			res += pad(1) + fmt.Sprintf("this.%s = make([]%s, n)\n", name(typ), name(typ))
		}
		res += "	for i := 0; i < n; i++ {\n"
	} else {
		depth = 1
	}

	if opt.cpool {
		res += emitReadI32(depth)
		if opt.doNotKeepData {

		} else {
			res += pad(depth) + fmt.Sprintf("id := %s(v32_)\n", refName(typ))
		}
	}

	res += generateBindLoop(typ, "bind", depth, true)

	if opt.cpool {
		if opt.doNotKeepData {

		} else {

			res += pad(depth) + fmt.Sprintf("this.%s[i] = bind.Temp\n", name(typ))

			if opt.sortedIDs {
				res += pad(depth) + "this.IDMap.Set(id, i)\n"
			} else {
				res += pad(depth) + "this.IDMap[id] = uint32(i)\n"
			}
		}
		res += pad(1) + "}\n"
	} else {
		res += pad(depth) + fmt.Sprintf("*this = bind.Temp\n")
	}
	res += "	return pos, nil\n"
	res += fmt.Sprintf("}\n")

	imports := "package types\n"
	imports += "\n"

	imports += "import (\n\t\"fmt\"\n\t\"io\"\n\t\"unsafe\"\n\t\"github.com/grafana/jfr-parser/parser/types/def\"\n\n)"

	imports += "\n"
	res = imports + res

	fmt.Println("types2.ExpectedMeta" + name(typ) + ",")
	//fmt.Println(res)
	return res
}

func generateBindLoop(typ *def.Class, bindName string, depth int, nestedAllowed bool) string {
	fs := getUniqueFields(typ)
	cpoolFields := getUniqueCpoolFields(typ)
	complexFields := getNonBasicFields(typ)
	_ = complexFields
	res := ""

	res += pad(depth) + fmt.Sprintf("for %sFieldIndex := 0; %sFieldIndex < len(%s.Fields); %sFieldIndex++ {\n", bindName, bindName, bindName, bindName)
	res += pad(depth) + fmt.Sprintf("	%sArraySize := 1\n", bindName)
	res += pad(depth) + fmt.Sprintf("	if %s.Fields[%sFieldIndex].Field.Array {\n", bindName, bindName)
	res += emitReadI32(depth + 2)
	res += pad(depth) + fmt.Sprintf("		%sArraySize = int(v32_)\n", bindName)
	if len(complexFields) > 0 {
		res += pad(depth) + fmt.Sprintf("		if %s.Fields[%sFieldIndex].Field.Type == typeMap.%s {\n", bindName, bindName, TypeID2Sym(complexFields[0].Type))
		res += pad(depth) + fmt.Sprintf("			*%s.Fields[%sFieldIndex].%s = make([]%s, 0, %sArraySize)\n",
			bindName, bindName, name(TypeForCPoolID(complexFields[0].Type)), name(TypeForCPoolID(complexFields[0].Type)), bindName)
		res += pad(depth) + fmt.Sprintf("		}\n")
	}
	res += pad(depth) + fmt.Sprintf("	}\n")
	res += pad(depth) + fmt.Sprintf("	for %sArrayIndex := 0; %sArrayIndex < %sArraySize; %sArrayIndex++ {\n", bindName, bindName, bindName, bindName)
	res += pad(depth) + fmt.Sprintf("	if %s.Fields[%sFieldIndex].Field.ConstantPool {\n", bindName, bindName)
	res += emitReadI32(depth + 2)
	if len(cpoolFields) > 0 {
		res += pad(depth) + fmt.Sprintf("		switch %s.Fields[%sFieldIndex].Field.Type {\n", bindName, bindName)
		for _, field := range cpoolFields {
			res += pad(depth) + fmt.Sprintf("		case typeMap.%s:\n", TypeID2Sym(field.Type))
			res += pad(depth) + fmt.Sprintf("			if %s.Fields[%sFieldIndex].%s != nil {\n", bindName, bindName, goTypeName(field))
			res += pad(depth) + fmt.Sprintf("				*%s.Fields[%sFieldIndex].%s = %s(v32_)\n", bindName, bindName, goTypeName(field), goTypeName(field))
			res += pad(depth) + fmt.Sprintf("			}\n")
		}
		res += pad(depth) + fmt.Sprintf("		}\n")
	}
	res += pad(depth) + fmt.Sprintf("	} else {\n")
	res += pad(depth) + fmt.Sprintf("		%sFieldTypeID := %s.Fields[%sFieldIndex].Field.Type\n", bindName, bindName, bindName)

	res += pad(depth) + fmt.Sprintf("		switch %sFieldTypeID {\n", bindName)
	res += pad(depth) + fmt.Sprintf("		case  typeMap.T_STRING:\n")
	res += emitString(depth + 3)
	if fieldsHas(fs, T_STRING) {
		res += pad(depth) + fmt.Sprintf("			if %s.Fields[%sFieldIndex].string != nil {\n", bindName, bindName)
		res += pad(depth) + fmt.Sprintf("				*%s.Fields[%sFieldIndex].string = s_\n", bindName, bindName)
		res += pad(depth) + fmt.Sprintf("			}\n")
	} else {
		res += pad(depth) + fmt.Sprintf("			// skipping\n")
	}

	res += pad(depth) + fmt.Sprintf("		case typeMap.T_INT:\n")
	res += emitReadI32(depth + 3)
	if fieldsHas(fs, T_INT) {
		res += pad(depth) + fmt.Sprintf("			if %s.Fields[%sFieldIndex].uint32 != nil {\n", bindName, bindName)
		res += pad(depth) + fmt.Sprintf("				*%s.Fields[%sFieldIndex].uint32 = v32_\n", bindName, bindName)
		res += pad(depth) + fmt.Sprintf("			}\n")
	} else {
		res += pad(depth) + fmt.Sprintf("			// skipping\n")
	}
	res += pad(depth) + fmt.Sprintf("		case typeMap.T_LONG:\n")
	res += emitReadU64(depth + 3)
	if fieldsHas(fs, T_LONG) {
		res += pad(depth) + fmt.Sprintf("			if %s.Fields[%sFieldIndex].uint64 != nil {\n", bindName, bindName)
		res += pad(depth) + fmt.Sprintf("				*%s.Fields[%sFieldIndex].uint64 = v64_\n", bindName, bindName)
		res += pad(depth) + fmt.Sprintf("			}\n")
	} else {
		res += pad(depth) + fmt.Sprintf("			// skipping\n")
	}

	res += pad(depth) + fmt.Sprintf("		case typeMap.T_BOOLEAN:\n")
	res += emitReadByte(depth + 3)
	if fieldsHas(fs, T_BOOLEAN) {
		res += pad(depth) + fmt.Sprintf("			if %s.Fields[%sFieldIndex].bool != nil {\n", bindName, bindName)
		res += pad(depth) + fmt.Sprintf("				*%s.Fields[%sFieldIndex].bool = b_ != 0\n", bindName, bindName)
		res += pad(depth) + fmt.Sprintf("			}\n")
	} else {
		res += pad(depth) + fmt.Sprintf("			// skipping\n")
	}
	res += pad(depth) + fmt.Sprintf("		case typeMap.T_FLOAT:\n")
	res += emitReadI32(depth + 3)
	if fieldsHas(fs, T_FLOAT) {
		res += pad(depth) + fmt.Sprintf("			if %s.Fields[%sFieldIndex].float32 != nil {\n", bindName, bindName)
		res += pad(depth) + fmt.Sprintf("				*%s.Fields[%sFieldIndex].float32 = *(*float32)(unsafe.Pointer(&v32_))\n", bindName, bindName)
		res += pad(depth) + fmt.Sprintf("			}\n")
	} else {
		res += pad(depth) + fmt.Sprintf("			// skipping\n")
	}
	if nestedAllowed {
		for _, field := range complexFields {
			nestedType := TypeForCPoolID(field.Type)
			res += pad(depth) + fmt.Sprintf("		case typeMap.%s:\n", TypeID2Sym(field.Type))
			res += generateBindLoop(nestedType, "bind"+name(nestedType), depth+3, false)
			if field.Array {
				res += pad(depth) + fmt.Sprintf("			if %s.Fields[%sFieldIndex].%s != nil {\n", bindName, bindName, name(nestedType))
				res += pad(depth) + fmt.Sprintf("				*%s.Fields[%sFieldIndex].%s = append(*%s.Fields[%sFieldIndex].%s, bind%s.Temp)\n", bindName, bindName, name(nestedType), bindName, bindName, name(nestedType), name(nestedType))
				res += pad(depth) + fmt.Sprintf("			}\n")
			} else {
				panic("TODO " + field.String())
			}
		}
	}
	res += pad(depth) + fmt.Sprintf("		default:\n")
	//todo array
	res += pad(depth) + fmt.Sprintf("			%sFieldType := typeMap.IDMap[%s.Fields[%sFieldIndex].Field.Type]\n", bindName, bindName, bindName)
	res += pad(depth) + fmt.Sprintf("			if %sFieldType == nil || len(%sFieldType.Fields) == 0 {\n", bindName, bindName)
	res += pad(depth) + fmt.Sprintf("				return 0, fmt.Errorf(\"unknown type %%d\", %s.Fields[%sFieldIndex].Field.Type)\n", bindName, bindName)
	res += pad(depth) + fmt.Sprintf("			}\n")
	res += pad(depth) + fmt.Sprintf("			%sSkipObjects := 1\n", bindName)
	res += pad(depth) + fmt.Sprintf("			if %s.Fields[%sFieldIndex].Field.Array {\n", bindName, bindName)
	res += emitReadI32(depth + 4)
	res += pad(depth) + fmt.Sprintf("				%sSkipObjects = int(v32_)\n", bindName)
	res += pad(depth) + fmt.Sprintf("			}\n")
	res += pad(depth) + fmt.Sprintf("			for %sSkipObjectIndex := 0; %sSkipObjectIndex < %sSkipObjects; %sSkipObjectIndex++ {\n", bindName, bindName, bindName, bindName)
	res += pad(depth) + fmt.Sprintf("				for %sskipFieldIndex := 0; %sskipFieldIndex < len(%sFieldType.Fields); %sskipFieldIndex++ {\n", bindName, bindName, bindName, bindName)
	res += pad(depth) + fmt.Sprintf("					%sSkipFieldType :=  %sFieldType.Fields[%sskipFieldIndex].Type\n", bindName, bindName, bindName)
	res += pad(depth) + fmt.Sprintf("					if %sFieldType.Fields[%sskipFieldIndex].ConstantPool {\n", bindName, bindName)
	res += emitReadI32(depth + 7)
	res += pad(depth) + fmt.Sprintf("					} else if %sSkipFieldType == typeMap.T_STRING{\n", bindName)
	res += emitString(depth + 7)
	res += pad(depth) + fmt.Sprintf("					} else if %sSkipFieldType == typeMap.T_INT {\n", bindName)
	res += emitReadI32(depth + 7)
	res += pad(depth) + fmt.Sprintf("					} else if %sSkipFieldType == typeMap.T_FLOAT {\n", bindName)
	res += emitReadI32(depth + 7)
	res += pad(depth) + fmt.Sprintf("					} else if %sSkipFieldType == typeMap.T_LONG {\n", bindName)
	res += emitReadU64(depth + 7)
	res += pad(depth) + fmt.Sprintf("					} else if %sSkipFieldType == typeMap.T_BOOLEAN {\n", bindName)
	res += emitReadByte(depth + 7)
	res += pad(depth) + fmt.Sprintf("					} else {\n")
	res += pad(depth) + fmt.Sprintf("							return 0, fmt.Errorf(\"nested objects not implemented. \")\n")
	res += pad(depth) + fmt.Sprintf("					}\n")
	res += pad(depth) + fmt.Sprintf("				}\n")
	res += pad(depth) + fmt.Sprintf("			}\n")
	res += pad(depth) + fmt.Sprintf("			}\n")
	res += pad(depth) + fmt.Sprintf("		}\n")
	res += pad(depth) + fmt.Sprintf("	}\n")
	res += pad(depth) + fmt.Sprintf("}\n")
	return res
}

func fieldsHas(fs []def.Field, tString def.TypeID) bool {
	for _, f := range fs {
		if f.Type == tString {
			return true
		}
	}
	return false
}

func getUniqueFields(typ *def.Class) []def.Field {
	res := make([]def.Field, 0, len(typ.Fields))
	for j := range typ.Fields {
		found := false
		fieldCopy := typ.Fields[j]
		fieldCopy.Name = ""
		for i := 0; i < len(res); i++ {
			if res[i].Equals(&fieldCopy) {
				found = true
				break
			}
		}
		if !found {
			res = append(res, fieldCopy)
		}
	}
	return res
}

func getUniqueCpoolFields(typ *def.Class) []def.Field {
	fs := getUniqueFields(typ)
	res := make([]def.Field, 0, len(fs))
	for _, f := range fs {
		if f.ConstantPool {
			res = append(res, f)
		}
	}
	return res
}

func getNonBasicFields(typ *def.Class) []def.Field {

	res := make([]def.Field, 0, len(typ.Fields))
	for _, f := range typ.Fields {
		if f.ConstantPool {
			continue
		}
		if f.Type == T_INT || f.Type == T_LONG || f.Type == T_FLOAT || f.Type == T_BOOLEAN || f.Type == T_STRING {
			continue
		}
		res = append(res, f)
	}
	return res
}

func generateBinding(typ *def.Class, opt options) string {
	res := ""
	res += fmt.Sprintf("type %s struct {\n", bindName(typ))
	res += fmt.Sprintf("	Temp %s\n", name(typ))
	res += fmt.Sprintf("	Fields []%s\n", bindFieldName(typ))
	res += fmt.Sprintf("}\n\n")
	res += fmt.Sprintf("\n")

	res += fmt.Sprintf("type %s struct {\n", bindFieldName(typ))
	res += fmt.Sprintf("	Field 	*def.Field\n")

	uniqueFields := getUniqueFields(typ)
	//written := map[string]bool{}
	for _, uf := range uniqueFields {
		//if !written[goTypeName(id)] {
		arr := ""
		if uf.Array {
			arr = "[]"
		}
		res += fmt.Sprintf("\t%s *%s%s\n", goTypeName(uf), arr, goTypeName(uf))
		//written[goTypeName(id)] = true
		//}
	}

	res += fmt.Sprintf("}\n\n")
	res += fmt.Sprintf("\n")

	res += fmt.Sprintf("func New%s(typ *def.Class, typeMap *def.TypeMap) *%s {\n", bindName(typ), bindName(typ))
	res += fmt.Sprintf("	res := new(%s)\n", bindName(typ))
	res += fmt.Sprintf("	res.Fields = make([]%s, 0, len(typ.Fields))\n", bindFieldName(typ))
	res += fmt.Sprintf("	for i := 0; i < len(typ.Fields); i++ {\n")
	res += fmt.Sprintf("		switch typ.Fields[i].Name {\n")
	for i := 0; i < len(typ.Fields); i++ {

		res += fmt.Sprintf("		case \"%s\":\n", typ.Fields[i].Name)
		if slices.Contains(opt.skipFields, typ.Fields[i].Name) {
			res += fmt.Sprintf("			res.Fields = append(res.Fields, %s{Field: &typ.Fields[i]}) // skip to save mem\n", bindFieldName(typ))
		} else {
			res += fmt.Sprintf("			if typ.Fields[i].Equals(&def.Field{Name: \"%s\", Type: typeMap.%s, ConstantPool: %v, Array: %v}) {\n", typ.Fields[i].Name, TypeID2Sym(typ.Fields[i].Type), typ.Fields[i].ConstantPool, typ.Fields[i].Array)
			res += fmt.Sprintf("				res.Fields = append(res.Fields, %s{Field: &typ.Fields[i], %s: &res.Temp.%s}) \n", bindFieldName(typ), goTypeName(typ.Fields[i]), capitalize(typ.Fields[i].Name))
			res += fmt.Sprintf("			} else {\n")
			res += fmt.Sprintf("				res.Fields = append(res.Fields, %s{Field: &typ.Fields[i]}) // skip changed field\n", bindFieldName(typ))
			res += fmt.Sprintf("			}\n")
		}
	}
	res += fmt.Sprintf("		default:\n")
	res += fmt.Sprintf("			res.Fields = append(res.Fields, %s{Field: &typ.Fields[i]}) // skip unknown new field\n", bindFieldName(typ))
	res += fmt.Sprintf("		}\n")
	res += fmt.Sprintf("	}\n")
	res += fmt.Sprintf("	return res\n")
	res += fmt.Sprintf("}\n")
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
		//todo array is not tested
		return name(TypeForCPoolID(field.Type)) + "Ref"
	}
	switch field.Type {
	case T_STRING:
		return "string"
	case T_LONG:
		return "uint64"
	case T_INT:
		return "uint32"
	case T_FLOAT:
		return "float32"
	case T_BOOLEAN:
		return "bool"
	case T_STACK_FRAME:
		return "StackFrame" //todo make it generic
	default:
		panic("TODO " + field.String())
	}
}

func name(typ *def.Class) string {
	fs := strings.Split(typ.Name, ".")
	s := fs[len(fs)-1]
	return capitalize(s)
}

func bindName(typ *def.Class) string {
	return "Bind" + name(typ)
}
func bindFieldName(typ *def.Class) string {
	return "BindField" + name(typ)
}
func refName(typ *def.Class) string {
	return name(typ) + "Ref"
}

func listName(typ *def.Class) string {
	return name(typ) + "List"
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
