package def

import "fmt"

type TypeID uint32

var (
	T_METADATA                = TypeID(0)
	T_CPOOL                   = TypeID(1)
	T_BOOLEAN                 = TypeID(4)
	T_CHAR                    = TypeID(5)
	T_FLOAT                   = TypeID(6)
	T_DOUBLE                  = TypeID(7)
	T_BYTE                    = TypeID(8)
	T_SHORT                   = TypeID(9)
	T_INT                     = TypeID(10)
	T_LONG                    = TypeID(11)
	T_STRING                  = TypeID(20)
	T_CLASS                   = TypeID(21)
	T_THREAD                  = TypeID(22)
	T_CLASS_LOADER            = TypeID(23)
	T_FRAME_TYPE              = TypeID(24)
	T_THREAD_STATE            = TypeID(25)
	T_STACK_TRACE             = TypeID(26)
	T_STACK_FRAME             = TypeID(27)
	T_METHOD                  = TypeID(28)
	T_PACKAGE                 = TypeID(29)
	T_SYMBOL                  = TypeID(30)
	T_LOG_LEVEL               = TypeID(31)
	T_EVENT                   = TypeID(100)
	T_EXECUTION_SAMPLE        = TypeID(101)
	T_ALLOC_IN_NEW_TLAB       = TypeID(102)
	T_ALLOC_OUTSIDE_TLAB      = TypeID(103)
	T_MONITOR_ENTER           = TypeID(104)
	T_THREAD_PARK             = TypeID(105)
	T_CPU_LOAD                = TypeID(106)
	T_ACTIVE_RECORDING        = TypeID(107)
	T_ACTIVE_SETTING          = TypeID(108)
	T_OS_INFORMATION          = TypeID(109)
	T_CPU_INFORMATION         = TypeID(110)
	T_JVM_INFORMATION         = TypeID(111)
	T_INITIAL_SYSTEM_PROPERTY = TypeID(112)
	T_NATIVE_LIBRARY          = TypeID(113)
	T_LOG                     = TypeID(114)
	T_LIVE_OBJECT             = TypeID(115)
	T_ANNOTATION              = TypeID(200)
	T_LABEL                   = TypeID(201)
	T_CATEGORY                = TypeID(202)
	T_TIMESTAMP               = TypeID(203)
	T_TIMESPAN                = TypeID(204)
	T_DATA_AMOUNT             = TypeID(205)
	T_MEMORY_ADDRESS          = TypeID(206)
	T_UNSIGNED                = TypeID(207)
	T_PERCENTAGE              = TypeID(208)
)

func TypeID2Sym(id TypeID) string {
	switch id {
	case T_METADATA:
		return "T_METADATA"
	case T_CPOOL:
		return "T_CPOOL"
	case T_BOOLEAN:
		return "T_BOOLEAN"
	case T_CHAR:
		return "T_CHAR"
	case T_FLOAT:
		return "T_FLOAT"
	case T_DOUBLE:
		return "T_DOUBLE"
	case T_BYTE:
		return "T_BYTE"
	case T_SHORT:
		return "T_SHORT"
	case T_INT:
		return "T_INT"
	case T_LONG:
		return "T_LONG"
	case T_STRING:
		return "T_STRING"
	case T_CLASS:
		return "T_CLASS"
	case T_THREAD:
		return "T_THREAD"
	case T_CLASS_LOADER:
		return "T_CLASS_LOADER"
	case T_FRAME_TYPE:
		return "T_FRAME_TYPE"
	case T_THREAD_STATE:
		return "T_THREAD_STATE"
	case T_STACK_TRACE:
		return "T_STACK_TRACE"
	case T_STACK_FRAME:
		return "T_STACK_FRAME"
	case T_METHOD:
		return "T_METHOD"
	case T_PACKAGE:
		return "T_PACKAGE"
	case T_SYMBOL:
		return "T_SYMBOL"
	case T_LOG_LEVEL:
		return "T_LOG_LEVEL"
	case T_EVENT:
		return "T_EVENT"
	case T_EXECUTION_SAMPLE:
		return "T_EXECUTION_SAMPLE"
	case T_ALLOC_IN_NEW_TLAB:
		return "T_ALLOC_IN_NEW_TLAB"
	case T_ALLOC_OUTSIDE_TLAB:
		return "T_ALLOC_OUTSIDE_TLAB"
	case T_MONITOR_ENTER:
		return "T_MONITOR_ENTER"
	case T_THREAD_PARK:
		return "T_THREAD_PARK"
	case T_CPU_LOAD:
		return "T_CPU_LOAD"
	case T_ACTIVE_RECORDING:
		return "T_ACTIVE_RECORDING"
	case T_ACTIVE_SETTING:
		return "T_ACTIVE_SETTING"
	case T_OS_INFORMATION:
		return "T_OS_INFORMATION"
	case T_CPU_INFORMATION:
		return "T_CPU_INFORMATION"
	case T_JVM_INFORMATION:
		return "T_JVM_INFORMATION"
	case T_INITIAL_SYSTEM_PROPERTY:
		return "T_INITIAL_SYSTEM_PROPERTY"
	case T_NATIVE_LIBRARY:
		return "T_NATIVE_LIBRARY"
	case T_LOG:
		return "T_LOG"
	case T_LIVE_OBJECT:
		return "T_LIVE_OBJECT"
	case T_ANNOTATION:
		return "T_ANNOTATION"
	case T_LABEL:
		return "T_LABEL"
	case T_CATEGORY:
		return "T_CATEGORY"
	case T_TIMESTAMP:
		return "T_TIMESTAMP"
	case T_TIMESPAN:
		return "T_TIMESPAN"
	case T_DATA_AMOUNT:
		return "T_DATA_AMOUNT"
	case T_MEMORY_ADDRESS:
		return "T_MEMORY_ADDRESS"
	case T_UNSIGNED:
		return "T_UNSIGNED"
	case T_PERCENTAGE:
		return "T_PERCENTAGE"
	default:
		return fmt.Sprintf("unknown type %d", id)
	}
}

var TypeStackFrame = Class{
	Name: "jdk.types.StackFrame",
	ID:   T_STACK_FRAME,
	Fields: []Field{
		{
			Name:         "method",
			Type:         T_METHOD,
			ConstantPool: true,
			Array:        false,
		},
		{
			Name:         "lineNumber",
			Type:         T_INT,
			ConstantPool: false,
			Array:        false,
		},
		{
			Name:         "bytecodeIndex",
			Type:         T_INT,
			ConstantPool: false,
			Array:        false,
		},
		{
			Name:         "type",
			Type:         T_FRAME_TYPE,
			ConstantPool: true,
			Array:        false,
		},
	},
}

type TypeMap struct {
	IDMap   map[TypeID]*Class
	NameMap map[string]*Class

	T_STRING  TypeID
	T_INT     TypeID
	T_LONG    TypeID
	T_FLOAT   TypeID
	T_BOOLEAN TypeID

	T_CLASS        TypeID
	T_THREAD       TypeID
	T_FRAME_TYPE   TypeID
	T_THREAD_STATE TypeID
	T_STACK_TRACE  TypeID
	T_METHOD       TypeID
	T_PACKAGE      TypeID
	T_SYMBOL       TypeID
	T_LOG_LEVEL    TypeID

	T_STACK_FRAME  TypeID
	T_CLASS_LOADER TypeID
}
