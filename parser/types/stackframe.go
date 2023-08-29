package types

type StackFrame struct {
	Method        MethodRef
	LineNumber    uint32
	BytecodeIndex uint32
	Type          FrameTypeRef
}
