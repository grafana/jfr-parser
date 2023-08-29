package main

import (
	"github.com/pyroscope-io/jfr-parser/parser/types/def"
)

var Type_boolean = def.Class{
	Name:   "boolean",
	ID:     def.T_BOOLEAN,
	Fields: []def.Field{},
}
var Type_char = def.Class{
	Name:   "char",
	ID:     def.T_CHAR,
	Fields: []def.Field{},
}
var Type_float = def.Class{
	Name:   "float",
	ID:     def.T_FLOAT,
	Fields: []def.Field{},
}
var Type_double = def.Class{
	Name:   "double",
	ID:     def.T_DOUBLE,
	Fields: []def.Field{},
}
var Type_byte = def.Class{
	Name:   "byte",
	ID:     def.T_BYTE,
	Fields: []def.Field{},
}
var Type_short = def.Class{
	Name:   "short",
	ID:     def.T_SHORT,
	Fields: []def.Field{},
}
var Type_int = def.Class{
	Name:   "int",
	ID:     def.T_INT,
	Fields: []def.Field{},
}
var Type_long = def.Class{
	Name:   "long",
	ID:     def.T_LONG,
	Fields: []def.Field{},
}
var Type_java_lang_String = def.Class{
	Name:   "java.lang.String",
	ID:     def.T_STRING,
	Fields: []def.Field{},
}
var Type_java_lang_Class = def.Class{
	Name: "java.lang.Class",
	ID:   def.T_CLASS,
	Fields: []def.Field{
		{Name: "classLoader", Type: def.T_CLASS_LOADER, ConstantPool: true},
		{Name: "name", Type: def.T_SYMBOL, ConstantPool: true},
		{Name: "package", Type: def.T_PACKAGE, ConstantPool: true},
		{Name: "modifiers", Type: def.T_INT, ConstantPool: false},
	},
}
var Type_java_lang_Thread = def.Class{
	Name: "java.lang.Thread",
	ID:   def.T_THREAD,
	Fields: []def.Field{
		{Name: "osName", Type: def.T_STRING, ConstantPool: false},
		{Name: "osThreadId", Type: def.T_LONG, ConstantPool: false},
		{Name: "javaName", Type: def.T_STRING, ConstantPool: false},
		{Name: "javaThreadId", Type: def.T_LONG, ConstantPool: false},
	},
}
var Type_jdk_types_ClassLoader = def.Class{
	Name: "jdk.types.ClassLoader",
	ID:   def.T_CLASS_LOADER,
	Fields: []def.Field{
		{Name: "type", Type: def.T_CLASS, ConstantPool: true},
		{Name: "name", Type: def.T_SYMBOL, ConstantPool: true},
	},
}
var Type_jdk_types_FrameType = def.Class{
	Name: "jdk.types.FrameType",
	ID:   def.T_FRAME_TYPE,
	Fields: []def.Field{
		{Name: "description", Type: def.T_STRING, ConstantPool: false},
	},
}
var Type_jdk_types_ThreadState = def.Class{
	Name: "jdk.types.ThreadState",
	ID:   def.T_THREAD_STATE,
	Fields: []def.Field{
		{Name: "name", Type: def.T_STRING, ConstantPool: false},
	},
}
var Type_jdk_types_StackTrace = def.Class{
	Name: "jdk.types.StackTrace",
	ID:   def.T_STACK_TRACE,
	Fields: []def.Field{
		{Name: "truncated", Type: def.T_BOOLEAN, ConstantPool: false},
		{Name: "frames", Type: def.T_STACK_FRAME, ConstantPool: false, Array: true},
	},
}
var Type_jdk_types_StackFrame = def.Class{
	Name: "jdk.types.StackFrame",
	ID:   def.T_STACK_FRAME,
	Fields: []def.Field{
		{Name: "method", Type: def.T_METHOD, ConstantPool: true},
		{Name: "lineNumber", Type: def.T_INT, ConstantPool: false},
		{Name: "bytecodeIndex", Type: def.T_INT, ConstantPool: false},
		{Name: "type", Type: def.T_FRAME_TYPE, ConstantPool: true},
	},
}
var Type_jdk_types_Method = def.Class{
	Name: "jdk.types.Method",
	ID:   def.T_METHOD,
	Fields: []def.Field{
		{Name: "type", Type: def.T_CLASS, ConstantPool: true},
		{Name: "name", Type: def.T_SYMBOL, ConstantPool: true},
		{Name: "descriptor", Type: def.T_SYMBOL, ConstantPool: true},
		{Name: "modifiers", Type: def.T_INT, ConstantPool: false},
		{Name: "hidden", Type: def.T_BOOLEAN, ConstantPool: false},
	},
}
var Type_jdk_types_Package = def.Class{
	Name: "jdk.types.Package",
	ID:   def.T_PACKAGE,
	Fields: []def.Field{
		{Name: "name", Type: def.T_SYMBOL, ConstantPool: true},
	},
}
var Type_jdk_types_Symbol = def.Class{
	Name: "jdk.types.Symbol",
	ID:   def.T_SYMBOL,
	Fields: []def.Field{
		{Name: "string", Type: def.T_STRING, ConstantPool: false},
	},
}
var Type_profiler_types_LogLevel = def.Class{
	Name: "profiler.types.LogLevel",
	ID:   def.T_LOG_LEVEL,
	Fields: []def.Field{
		{Name: "name", Type: def.T_STRING, ConstantPool: false},
	},
}
var Type_jdk_ExecutionSample = def.Class{
	Name: "jdk.ExecutionSample",
	ID:   def.T_EXECUTION_SAMPLE,
	Fields: []def.Field{
		{Name: "startTime", Type: def.T_LONG, ConstantPool: false},
		{Name: "sampledThread", Type: def.T_THREAD, ConstantPool: true},
		{Name: "stackTrace", Type: def.T_STACK_TRACE, ConstantPool: true},
		{Name: "state", Type: def.T_THREAD_STATE, ConstantPool: true},
		{Name: "contextId", Type: def.T_LONG, ConstantPool: false},
	},
}
var Type_jdk_ObjectAllocationInNewTLAB = def.Class{
	Name: "jdk.ObjectAllocationInNewTLAB",
	ID:   def.T_ALLOC_IN_NEW_TLAB,
	Fields: []def.Field{
		{Name: "startTime", Type: def.T_LONG, ConstantPool: false},
		{Name: "eventThread", Type: def.T_THREAD, ConstantPool: true},
		{Name: "stackTrace", Type: def.T_STACK_TRACE, ConstantPool: true},
		{Name: "objectClass", Type: def.T_CLASS, ConstantPool: true},
		{Name: "allocationSize", Type: def.T_LONG, ConstantPool: false},
		{Name: "tlabSize", Type: def.T_LONG, ConstantPool: false},
		{Name: "contextId", Type: def.T_LONG, ConstantPool: false},
	},
}
var Type_jdk_ObjectAllocationOutsideTLAB = def.Class{
	Name: "jdk.ObjectAllocationOutsideTLAB",
	ID:   def.T_ALLOC_OUTSIDE_TLAB,
	Fields: []def.Field{
		{Name: "startTime", Type: def.T_LONG, ConstantPool: false},
		{Name: "eventThread", Type: def.T_THREAD, ConstantPool: true},
		{Name: "stackTrace", Type: def.T_STACK_TRACE, ConstantPool: true},
		{Name: "objectClass", Type: def.T_CLASS, ConstantPool: true},
		{Name: "allocationSize", Type: def.T_LONG, ConstantPool: false},
		{Name: "contextId", Type: def.T_LONG, ConstantPool: false},
	},
}
var Type_jdk_JavaMonitorEnter = def.Class{
	Name: "jdk.JavaMonitorEnter",
	ID:   def.T_MONITOR_ENTER,
	Fields: []def.Field{
		{Name: "startTime", Type: def.T_LONG, ConstantPool: false},
		{Name: "duration", Type: def.T_LONG, ConstantPool: false},
		{Name: "eventThread", Type: def.T_THREAD, ConstantPool: true},
		{Name: "stackTrace", Type: def.T_STACK_TRACE, ConstantPool: true},
		{Name: "monitorClass", Type: def.T_CLASS, ConstantPool: true},
		{Name: "previousOwner", Type: def.T_THREAD, ConstantPool: true},
		{Name: "address", Type: def.T_LONG, ConstantPool: false},
		{Name: "contextId", Type: def.T_LONG, ConstantPool: false},
	},
}
var Type_jdk_ThreadPark = def.Class{
	Name: "jdk.ThreadPark",
	ID:   def.T_THREAD_PARK,
	Fields: []def.Field{
		{Name: "startTime", Type: def.T_LONG, ConstantPool: false},
		{Name: "duration", Type: def.T_LONG, ConstantPool: false},
		{Name: "eventThread", Type: def.T_THREAD, ConstantPool: true},
		{Name: "stackTrace", Type: def.T_STACK_TRACE, ConstantPool: true},
		{Name: "parkedClass", Type: def.T_CLASS, ConstantPool: true},
		{Name: "timeout", Type: def.T_LONG, ConstantPool: false},
		{Name: "until", Type: def.T_LONG, ConstantPool: false},
		{Name: "address", Type: def.T_LONG, ConstantPool: false},
		{Name: "contextId", Type: def.T_LONG, ConstantPool: false},
	},
}
var Type_jdk_CPULoad = def.Class{
	Name: "jdk.CPULoad",
	ID:   def.T_CPU_LOAD,
	Fields: []def.Field{
		{Name: "startTime", Type: def.T_LONG, ConstantPool: false},
		{Name: "jvmUser", Type: def.T_FLOAT, ConstantPool: false},
		{Name: "jvmSystem", Type: def.T_FLOAT, ConstantPool: false},
		{Name: "machineTotal", Type: def.T_FLOAT, ConstantPool: false},
	},
}
var Type_jdk_ActiveRecording = def.Class{
	Name: "jdk.ActiveRecording",
	ID:   def.T_ACTIVE_RECORDING,
	Fields: []def.Field{
		{Name: "startTime", Type: def.T_LONG, ConstantPool: false},
		{Name: "duration", Type: def.T_LONG, ConstantPool: false},
		{Name: "eventThread", Type: def.T_THREAD, ConstantPool: true},
		{Name: "id", Type: def.T_LONG, ConstantPool: false},
		{Name: "name", Type: def.T_STRING, ConstantPool: false},
		{Name: "destination", Type: def.T_STRING, ConstantPool: false},
		{Name: "maxAge", Type: def.T_LONG, ConstantPool: false},
		{Name: "maxSize", Type: def.T_LONG, ConstantPool: false},
		{Name: "recordingStart", Type: def.T_LONG, ConstantPool: false},
		{Name: "recordingDuration", Type: def.T_LONG, ConstantPool: false},
	},
}
var Type_jdk_ActiveSetting = def.Class{
	Name: "jdk.ActiveSetting",
	ID:   def.T_ACTIVE_SETTING,
	Fields: []def.Field{
		{Name: "startTime", Type: def.T_LONG, ConstantPool: false},
		{Name: "duration", Type: def.T_LONG, ConstantPool: false},
		{Name: "eventThread", Type: def.T_THREAD, ConstantPool: true},
		{Name: "stackTrace", Type: def.T_STACK_TRACE, ConstantPool: true},
		{Name: "id", Type: def.T_LONG, ConstantPool: false},
		{Name: "name", Type: def.T_STRING, ConstantPool: false},
		{Name: "value", Type: def.T_STRING, ConstantPool: false},
	},
}
var Type_jdk_OSInformation = def.Class{
	Name: "jdk.OSInformation",
	ID:   def.T_OS_INFORMATION,
	Fields: []def.Field{
		{Name: "startTime", Type: def.T_LONG, ConstantPool: false},
		{Name: "osVersion", Type: def.T_STRING, ConstantPool: false},
	},
}
var Type_jdk_CPUInformation = def.Class{
	Name: "jdk.CPUInformation",
	ID:   def.T_CPU_INFORMATION,
	Fields: []def.Field{
		{Name: "startTime", Type: def.T_LONG, ConstantPool: false},
		{Name: "cpu", Type: def.T_STRING, ConstantPool: false},
		{Name: "description", Type: def.T_STRING, ConstantPool: false},
		{Name: "sockets", Type: def.T_INT, ConstantPool: false},
		{Name: "cores", Type: def.T_INT, ConstantPool: false},
		{Name: "hwThreads", Type: def.T_INT, ConstantPool: false},
	},
}
var Type_jdk_JVMInformation = def.Class{
	Name: "jdk.JVMInformation",
	ID:   def.T_JVM_INFORMATION,
	Fields: []def.Field{
		{Name: "startTime", Type: def.T_LONG, ConstantPool: false},
		{Name: "jvmName", Type: def.T_STRING, ConstantPool: false},
		{Name: "jvmVersion", Type: def.T_STRING, ConstantPool: false},
		{Name: "jvmArguments", Type: def.T_STRING, ConstantPool: false},
		{Name: "jvmFlags", Type: def.T_STRING, ConstantPool: false},
		{Name: "javaArguments", Type: def.T_STRING, ConstantPool: false},
		{Name: "jvmStartTime", Type: def.T_LONG, ConstantPool: false},
		{Name: "pid", Type: def.T_LONG, ConstantPool: false},
	},
}
var Type_jdk_InitialSystemProperty = def.Class{
	Name: "jdk.InitialSystemProperty",
	ID:   def.T_INITIAL_SYSTEM_PROPERTY,
	Fields: []def.Field{
		{Name: "startTime", Type: def.T_LONG, ConstantPool: false},
		{Name: "key", Type: def.T_STRING, ConstantPool: false},
		{Name: "value", Type: def.T_STRING, ConstantPool: false},
	},
}
var Type_jdk_NativeLibrary = def.Class{
	Name: "jdk.NativeLibrary",
	ID:   def.T_NATIVE_LIBRARY,
	Fields: []def.Field{
		{Name: "startTime", Type: def.T_LONG, ConstantPool: false},
		{Name: "name", Type: def.T_STRING, ConstantPool: false},
		{Name: "baseAddress", Type: def.T_LONG, ConstantPool: false},
		{Name: "topAddress", Type: def.T_LONG, ConstantPool: false},
	},
}
var Type_profiler_Log = def.Class{
	Name: "profiler.Log",
	ID:   def.T_LOG,
	Fields: []def.Field{
		{Name: "startTime", Type: def.T_LONG, ConstantPool: false},
		{Name: "level", Type: def.T_LOG_LEVEL, ConstantPool: true},
		{Name: "message", Type: def.T_STRING, ConstantPool: false},
	},
}
var Type_profiler_LiveObject = def.Class{
	Name: "profiler.LiveObject",
	ID:   def.T_LIVE_OBJECT,
	Fields: []def.Field{
		{Name: "startTime", Type: def.T_LONG, ConstantPool: false},
		{Name: "eventThread", Type: def.T_THREAD, ConstantPool: true},
		{Name: "stackTrace", Type: def.T_STACK_TRACE, ConstantPool: true},
		{Name: "objectClass", Type: def.T_CLASS, ConstantPool: true},
		{Name: "allocationSize", Type: def.T_LONG, ConstantPool: false},
		{Name: "allocationTime", Type: def.T_LONG, ConstantPool: false},
	},
}
var Type_jdk_jfr_Label = def.Class{
	Name: "jdk.jfr.Label",
	ID:   def.T_LABEL,
	Fields: []def.Field{
		{Name: "value", Type: def.T_STRING, ConstantPool: false},
	},
}
var Type_jdk_jfr_Category = def.Class{
	Name: "jdk.jfr.Category",
	ID:   def.T_CATEGORY,
	Fields: []def.Field{
		{Name: "value", Type: def.T_STRING, ConstantPool: false, Array: true},
	},
}
var Type_jdk_jfr_Timestamp = def.Class{
	Name: "jdk.jfr.Timestamp",
	ID:   def.T_TIMESTAMP,
	Fields: []def.Field{
		{Name: "value", Type: def.T_STRING, ConstantPool: false},
	},
}
var Type_jdk_jfr_Timespan = def.Class{
	Name: "jdk.jfr.Timespan",
	ID:   def.T_TIMESPAN,
	Fields: []def.Field{
		{Name: "value", Type: def.T_STRING, ConstantPool: false},
	},
}
var Type_jdk_jfr_DataAmount = def.Class{
	Name: "jdk.jfr.DataAmount",
	ID:   def.T_DATA_AMOUNT,
	Fields: []def.Field{
		{Name: "value", Type: def.T_STRING, ConstantPool: false},
	},
}
var Type_jdk_jfr_MemoryAddress = def.Class{
	Name:   "jdk.jfr.MemoryAddress",
	ID:     def.T_MEMORY_ADDRESS,
	Fields: []def.Field{},
}
var Type_jdk_jfr_Unsigned = def.Class{
	Name:   "jdk.jfr.Unsigned",
	ID:     def.T_UNSIGNED,
	Fields: []def.Field{},
}
var Type_jdk_jfr_Percentage = def.Class{
	Name:   "jdk.jfr.Percentage",
	ID:     def.T_PERCENTAGE,
	Fields: []def.Field{},
}
