package parser

import (
	"fmt"
	"io"
	"unsafe"

	types2 "github.com/pyroscope-io/jfr-parser/parser/types"
	"github.com/pyroscope-io/jfr-parser/parser/types/def"
)

const chunkHeaderSize = 68
const bufferSize = 1024 * 1024
const chunkMagic = 0x464c5200

type ChunkHeader struct {
	Magic              uint32
	Version            uint32
	Size               int
	OffsetConstantPool int
	OffsetMeta         int
	StartNanos         uint64
	DurationNanos      uint64
	StartTicks         uint64
	TicksPerSecond     uint64
	Features           uint32
}

func (c *ChunkHeader) String() string {
	return fmt.Sprintf("ChunkHeader{Magic: %x, Version: %x, Size: %d, OffsetConstantPool: %d, OffsetMeta: %d, StartNanos: %d, DurationNanos: %d, StartTicks: %d, TicksPerSecond: %d, Features: %d}", c.Magic, c.Version, c.Size, c.OffsetConstantPool, c.OffsetMeta, c.StartNanos, c.DurationNanos, c.StartTicks, c.TicksPerSecond, c.Features)
}

type SymbolProcessor func(ref *types2.Symbol)

type Options struct {
	ChunkSizeLimit  int
	SymbolProcessor SymbolProcessor
}

type Parser struct {
	FrameTypes   types2.FrameTypeList
	ThreadStates types2.ThreadStateList
	Threads      types2.ThreadList
	Classes      types2.ClassList
	Methods      types2.MethodList
	Packages     types2.PackageList
	Symbols      types2.SymbolList
	LogLevels    types2.LogLevelList
	Stacktrace   types2.StackTraceList

	ExecutionSample             types2.ExecutionSample
	ObjectAllocationInNewTLAB   types2.ObjectAllocationInNewTLAB
	ObjectAllocationOutsideTLAB types2.ObjectAllocationOutsideTLAB
	JavaMonitorEnter            types2.JavaMonitorEnter
	ThreadPark                  types2.ThreadPark
	LiveObject                  types2.LiveObject
	ActiveSetting               types2.ActiveSetting

	header   ChunkHeader
	options  Options
	buf      []byte
	pos      int
	metaSize uint32
	chunkEnd int

	TypeMap def.TypeMap

	executionSampleHaveContextId       bool
	allocationInNewTlabHaveContextId   bool
	allocationOutsideTlabHaveContextId bool
	monitorEnterHaveContextId          bool
	threadParkHaveContextId            bool

	typeExecutionSample  *def.Class
	typeAllocInNewTLAB   *def.Class
	typeALlocOutsideTLAB *def.Class
	typeMonitorEnter     *def.Class
	typeThreadPark       *def.Class
	typeLiveObject       *def.Class
	typeActiveSetting    *def.Class

	bindFrameType   *types2.BindFrameType
	bindThreadState *types2.BindThreadState
	bindThread      *types2.BindThread
	bindClass       *types2.BindClass
	bindMethod      *types2.BindMethod
	bindPackage     *types2.BindPackage
	bindSymbol      *types2.BindSymbol
	bindLogLevel    *types2.BindLogLevel
}

func NewParser(buf []byte, options Options) (res *Parser, err error) {
	p := &Parser{
		options: options,
		buf:     buf,
	}
	if err := p.readChunk(0); err != nil {
		return nil, err
	}
	return p, nil
}

func (p *Parser) ParseEvent() (def.TypeID, error) {
	for {
		if p.pos == p.chunkEnd {
			if p.pos == len(p.buf) {
				return 0, io.EOF
			}
			if err := p.readChunk(p.pos); err != nil {
				return 0, err
			}
		}
		pp := p.pos
		size, err := p.varInt()
		if err != nil {
			return 0, err
		}
		if size == 0 {
			return 0, def.ErrIntOverflow
		}
		typ, err := p.varInt()
		if err != nil {
			return 0, err
		}
		_ = size

		ttyp := def.TypeID(typ)
		switch ttyp {
		case def.T_EXECUTION_SAMPLE:
			if p.typeExecutionSample == nil {
				p.pos = pp + int(size) // skip
				continue
			}
			_, err := p.ExecutionSample.Parse(p.buf[p.pos:], p.typeExecutionSample, p.TypeMap.IDMap, p.executionSampleHaveContextId)
			if err != nil {
				return 0, err
			}
			p.pos = pp + int(size)
			return ttyp, nil
		case def.T_ALLOC_IN_NEW_TLAB:
			if p.typeAllocInNewTLAB == nil {
				p.pos = pp + int(size) // skip
				continue
			}
			_, err := p.ObjectAllocationInNewTLAB.Parse(p.buf[p.pos:], p.typeAllocInNewTLAB, p.TypeMap.IDMap, p.allocationInNewTlabHaveContextId)
			if err != nil {
				return 0, err
			}
			p.pos = pp + int(size)
			return ttyp, nil
		case def.T_ALLOC_OUTSIDE_TLAB:
			if p.typeALlocOutsideTLAB == nil {
				p.pos = pp + int(size) // skip
				continue
			}
			_, err := p.ObjectAllocationOutsideTLAB.Parse(p.buf[p.pos:], p.typeALlocOutsideTLAB, p.TypeMap.IDMap, p.allocationOutsideTlabHaveContextId)
			if err != nil {
				return 0, err
			}
			p.pos = pp + int(size)
			return ttyp, nil
		case def.T_LIVE_OBJECT:
			if p.typeLiveObject == nil {
				p.pos = pp + int(size) // skip
				continue
			}
			_, err := p.LiveObject.Parse(p.buf[p.pos:], p.typeLiveObject, p.TypeMap.IDMap)
			if err != nil {
				return 0, err
			}
			p.pos = pp + int(size)
			return ttyp, nil
		case def.T_MONITOR_ENTER:
			if p.typeMonitorEnter == nil {
				p.pos = pp + int(size) // skip
				continue
			}
			_, err := p.JavaMonitorEnter.Parse(p.buf[p.pos:], p.typeMonitorEnter, p.TypeMap.IDMap, p.monitorEnterHaveContextId)
			if err != nil {
				return 0, err
			}
			p.pos = pp + int(size)
			return ttyp, nil
		case def.T_THREAD_PARK:
			if p.typeThreadPark == nil {
				p.pos = pp + int(size) // skip
				continue
			}
			_, err := p.ThreadPark.Parse(p.buf[p.pos:], p.typeThreadPark, p.TypeMap.IDMap, p.threadParkHaveContextId)
			if err != nil {
				return 0, err
			}
			p.pos = pp + int(size)
			return ttyp, nil
		default:
			//fmt.Printf("skipping %s %v\n", def.TypeID2Sym(ttyp), ttyp)
			p.pos = pp + int(size)
		}
	}
}
func (p *Parser) GetStacktrace(stID types2.StackTraceRef) *types2.StackTrace {
	idx, ok := p.Stacktrace.IDMap[stID]
	if !ok {
		return nil
	}
	return &p.Stacktrace.StackTrace[idx]
}

func (p *Parser) GetMethod(mID types2.MethodRef) *types2.Method {
	if mID == 0 {
		return nil
	}
	var idx int

	refIDX := int(mID)
	if refIDX < len(p.Methods.IDMap.Slice) {
		idx = int(p.Methods.IDMap.Slice[mID])
	} else {
		idx = p.Methods.IDMap.Get(mID)
	}

	if idx == -1 {
		return nil
	}
	return &p.Methods.Method[idx]
}

func (p *Parser) GetClass(cID types2.ClassRef) *types2.Class {
	idx, ok := p.Classes.IDMap[cID]
	if !ok {
		return nil
	}
	return &p.Classes.Class[idx]
}

func (p *Parser) GetSymbol(sID types2.SymbolRef) *types2.Symbol {
	idx, ok := p.Symbols.IDMap[sID]
	if !ok {
		return nil
	}
	return &p.Symbols.Symbol[idx]
}

func (p *Parser) GetSymbolString(sID types2.SymbolRef) string {
	idx, ok := p.Symbols.IDMap[sID]
	if !ok {
		return ""
	}
	return p.Symbols.Symbol[idx].String
}

func (p *Parser) readChunk(pos int) error {
	if err := p.readChunkHeader(pos); err != nil {
		return fmt.Errorf("error reading chunk header: %w", err)
	}

	if err := p.readMeta(pos + p.header.OffsetMeta); err != nil {
		return fmt.Errorf("error reading metadata: %w", err)
	}
	if err := p.readConstantPool(pos + p.header.OffsetConstantPool); err != nil {
		return fmt.Errorf("error reading CP: %w", err)
	}
	pp := p.options.SymbolProcessor
	if pp != nil {
		for i := range p.Symbols.Symbol {
			pp(&p.Symbols.Symbol[i])
		}
	}
	p.pos = pos + chunkHeaderSize + int(p.metaSize)
	return nil
}

func (p *Parser) seek(pos int) error {
	if pos < len(p.buf) {
		p.pos = pos
		return nil
	}
	return io.ErrUnexpectedEOF
}

func (p *Parser) byte() (byte, error) {
	if p.pos >= len(p.buf) {
		return 0, io.ErrUnexpectedEOF
	}
	b := p.buf[p.pos]
	p.pos++
	return b, nil
}
func (p *Parser) varInt() (uint32, error) {
	v := uint32(0)
	for shift := uint(0); ; shift += 7 {
		if shift >= 32 {
			return 0, def.ErrIntOverflow
		}
		if p.pos >= len(p.buf) {
			return 0, io.ErrUnexpectedEOF
		}
		b := p.buf[p.pos]
		p.pos++
		v |= uint32(b&0x7F) << shift
		if b < 0x80 {
			break
		}
	}
	return v, nil
}

func (p *Parser) varLong() (uint64, error) {
	var v uint64
	for shift := uint(0); ; shift += 7 {
		if shift >= 64 {
			return 0, def.ErrIntOverflow
		}
		if p.pos >= len(p.buf) {
			return 0, io.ErrUnexpectedEOF
		}
		b := p.buf[p.pos]
		p.pos++
		v |= uint64(b&0x7F) << shift
		if b < 0x80 {
			break
		}
	}
	return v, nil
}

func (p *Parser) string() (string, error) {
	if p.pos >= len(p.buf) {
		return "", io.ErrUnexpectedEOF
	}
	b := p.buf[p.pos]
	p.pos++
	switch b {
	case 0:
		return "", nil //todo this should be nil
	case 1:
		return "", nil
	case 3:
		bs, err := p.bytes()
		if err != nil {
			return "", err
		}
		str := *(*string)(unsafe.Pointer(&bs))
		return str, nil
	default:
		return "", fmt.Errorf("unknown string type %d", b)
	}

}

func (p *Parser) bytes() ([]byte, error) {
	l, err := p.varInt()
	if err != nil {
		return nil, err
	}
	if l < 0 {
		return nil, def.ErrIntOverflow
	}
	if p.pos+int(l) > len(p.buf) {
		return nil, io.ErrUnexpectedEOF
	}
	bs := p.buf[p.pos : p.pos+int(l)]
	p.pos += int(l)
	return bs, nil
}

func (p *Parser) checkTypes() error {
	var err error
	tint := p.TypeMap.NameMap["int"]
	tlong := p.TypeMap.NameMap["long"]
	tfloat := p.TypeMap.NameMap["float"]
	tboolean := p.TypeMap.NameMap["boolean"]
	tstring := p.TypeMap.NameMap["java.lang.String"]

	if tint == nil {
		return fmt.Errorf("missing \"int\"")
	}
	if tlong == nil {
		return fmt.Errorf("missing \"long\"")
	}
	if tfloat == nil {
		return fmt.Errorf("missing \"float\"")
	}
	if tboolean == nil {
		return fmt.Errorf("missing \"boolean\"")
	}
	if tstring == nil {
		return fmt.Errorf("missing \"java.lang.String\"")
	}
	p.TypeMap.T_INT = tint.ID
	p.TypeMap.T_LONG = tlong.ID
	p.TypeMap.T_FLOAT = tfloat.ID
	p.TypeMap.T_BOOLEAN = tboolean.ID
	p.TypeMap.T_STRING = tstring.ID

	p.typeExecutionSample = p.TypeMap.NameMap["jdk.ExecutionSample"]
	p.typeAllocInNewTLAB = p.TypeMap.NameMap["jdk.ObjectAllocationInNewTLAB"]
	p.typeALlocOutsideTLAB = p.TypeMap.NameMap["jdk.ObjectAllocationOutsideTLAB"]
	p.typeMonitorEnter = p.TypeMap.NameMap["jdk.JavaMonitorEnter"]
	p.typeThreadPark = p.TypeMap.NameMap["jdk.ThreadPark"]
	p.typeLiveObject = p.TypeMap.NameMap["profiler.LiveObject"]
	p.typeActiveSetting = p.TypeMap.NameMap["jdk.ActiveSetting"]

	p.executionSampleHaveContextId = p.typeExecutionSample != nil && p.typeExecutionSample.Field("contextId") != nil
	p.allocationInNewTlabHaveContextId = p.typeAllocInNewTLAB != nil && p.typeAllocInNewTLAB.Field("contextId") != nil
	p.allocationOutsideTlabHaveContextId = p.typeALlocOutsideTLAB != nil && p.typeALlocOutsideTLAB.Field("contextId") != nil
	p.monitorEnterHaveContextId = p.typeMonitorEnter != nil && p.typeMonitorEnter.Field("contextId") != nil
	p.threadParkHaveContextId = p.typeThreadPark != nil && p.typeThreadPark.Field("contextId") != nil

	typeCPFrameType := p.TypeMap.NameMap["jdk.types.FrameType"]
	typeCPThreadState := p.TypeMap.NameMap["jdk.types.ThreadState"]
	typeCPThread := p.TypeMap.NameMap["java.lang.Thread"]
	typeCPClass := p.TypeMap.NameMap["java.lang.Class"]
	typeCPMethod := p.TypeMap.NameMap["jdk.types.Method"]
	typeCPPackage := p.TypeMap.NameMap["jdk.types.Package"]
	typeCPSymbol := p.TypeMap.NameMap["jdk.types.Symbol"]
	typeCPLogLevel := p.TypeMap.NameMap["profiler.types.LogLevel"]
	typeCPStackTrace := p.TypeMap.NameMap["jdk.types.StackTrace"]
	typeCPClassLoader := p.TypeMap.NameMap["jdk.types.ClassLoader"]

	if typeCPFrameType == nil {
		return fmt.Errorf("missing \"jdk.types.FrameType\"")
	}
	if typeCPThreadState == nil {
		return fmt.Errorf("missing \"jdk.types.ThreadState\"")
	}
	if typeCPThread == nil {
		return fmt.Errorf("missing \"java.lang.Thread\"")
	}
	if typeCPClass == nil {
		return fmt.Errorf("missing \"java.lang.Class\"")
	}
	if typeCPMethod == nil {
		return fmt.Errorf("missing \"jdk.types.Method\"")
	}
	if typeCPPackage == nil {
		return fmt.Errorf("missing \"jdk.types.Package\"")
	}
	if typeCPSymbol == nil {
		return fmt.Errorf("missing \"jdk.types.Symbol\"")
	}
	if typeCPLogLevel == nil {
		return fmt.Errorf("missing \"profiler.types.LogLevel\"")
	}
	if typeCPStackTrace == nil {
		return fmt.Errorf("missing \"jdk.types.StackTrace\"")
	}
	if typeCPClassLoader == nil {
		return fmt.Errorf("missing \"jdk.types.ClassLoader\"")
	}
	p.TypeMap.T_FRAME_TYPE = typeCPFrameType.ID
	p.TypeMap.T_THREAD_STATE = typeCPThreadState.ID
	p.TypeMap.T_THREAD = typeCPThread.ID
	p.TypeMap.T_CLASS = typeCPClass.ID
	p.TypeMap.T_METHOD = typeCPMethod.ID
	p.TypeMap.T_PACKAGE = typeCPPackage.ID
	p.TypeMap.T_SYMBOL = typeCPSymbol.ID
	p.TypeMap.T_LOG_LEVEL = typeCPLogLevel.ID
	p.TypeMap.T_STACK_TRACE = typeCPStackTrace.ID
	p.TypeMap.T_CLASS_LOADER = typeCPClassLoader.ID

	typeStackFrame := p.TypeMap.NameMap["jdk.types.StackFrame"]

	if typeStackFrame == nil {
		return fmt.Errorf("missing \"jdk.types.StackFrame\"")
	}
	p.TypeMap.T_STACK_FRAME = typeStackFrame.ID

	p.bindFrameType, err = types2.NewBindFrameType(typeCPFrameType, &p.TypeMap)
	if err != nil {
		return fmt.Errorf("unsupported jdk.types.FrameType %w", err)
	}
	p.bindThreadState, err = types2.NewBindThreadState(typeCPThreadState, &p.TypeMap)
	if err != nil {
		return fmt.Errorf("unsupported jdk.types.ThreadState %w", err)
	}
	p.bindThread, err = types2.NewBindThread(typeCPThread, &p.TypeMap)
	if err != nil {
		return fmt.Errorf("unsupported java.lang.Thread %w", err)
	}
	p.bindClass, err = types2.NewBindClass(typeCPClass, &p.TypeMap)
	if err != nil {
		return fmt.Errorf("unsupported java.lang.Class %w", err)
	}
	p.bindMethod, err = types2.NewBindMethod(typeCPMethod, &p.TypeMap)
	if err != nil {
		return fmt.Errorf("unsupported jdk.types.Method %w", err)
	}
	p.bindPackage, err = types2.NewBindPackage(typeCPPackage, &p.TypeMap)
	if err != nil {
		return fmt.Errorf("unsupported jdk.types.Package %w", err)
	}
	p.bindSymbol, err = types2.NewBindSymbol(typeCPSymbol, &p.TypeMap)
	if err != nil {
		return fmt.Errorf("unsupported jdk.types.Symbol %w", err)
	}
	p.bindLogLevel, err = types2.NewBindLogLevel(typeCPLogLevel, &p.TypeMap)
	if err != nil {
		return fmt.Errorf("unsupported profiler.types.LogLevel %w", err)
	}

	p.FrameTypes.IDMap = nil
	p.ThreadStates.IDMap = nil
	p.Threads.IDMap = nil
	p.Classes.IDMap = nil
	p.Methods.IDMap.Slice = nil
	p.Packages.IDMap = nil
	p.Symbols.IDMap = nil
	p.LogLevels.IDMap = nil
	p.Stacktrace.IDMap = nil
	return nil
}
