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

	typeMap map[def.TypeID]*def.Class

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
		if ttyp == def.T_EXECUTION_SAMPLE {
			_, err := p.ExecutionSample.Parse(p.buf[p.pos:], p.typeExecutionSample, p.typeMap, p.executionSampleHaveContextId)
			if err != nil {
				return 0, err
			}
			p.pos = pp + int(size)
			return ttyp, nil
		} else if ttyp == def.T_ALLOC_IN_NEW_TLAB {
			_, err := p.ObjectAllocationInNewTLAB.Parse(p.buf[p.pos:], p.typeAllocInNewTLAB, p.typeMap, p.allocationInNewTlabHaveContextId)
			if err != nil {
				return 0, err
			}
			p.pos = pp + int(size)
			return ttyp, nil
		} else if ttyp == def.T_ALLOC_OUTSIDE_TLAB {
			_, err := p.ObjectAllocationOutsideTLAB.Parse(p.buf[p.pos:], p.typeALlocOutsideTLAB, p.typeMap, p.allocationOutsideTlabHaveContextId)
			if err != nil {
				return 0, err
			}
			p.pos = pp + int(size)
			return ttyp, nil
		} else if ttyp == def.T_LIVE_OBJECT {
			_, err := p.LiveObject.Parse(p.buf[p.pos:], p.typeLiveObject, p.typeMap)
			if err != nil {
				return 0, err
			}
			p.pos = pp + int(size)
			return ttyp, nil
		} else if ttyp == def.T_MONITOR_ENTER {
			_, err := p.JavaMonitorEnter.Parse(p.buf[p.pos:], p.typeMonitorEnter, p.typeMap, p.monitorEnterHaveContextId)
			if err != nil {
				return 0, err
			}
			p.pos = pp + int(size)
			return ttyp, nil
		} else if ttyp == def.T_THREAD_PARK {
			_, err := p.ThreadPark.Parse(p.buf[p.pos:], p.typeThreadPark, p.typeMap, p.threadParkHaveContextId)
			if err != nil {
				return 0, err
			}
			p.pos = pp + int(size)
			return ttyp, nil
		} else {
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
	var expected = []*def.Class{
		types2.ExpectedMetaFrameType,
		types2.ExpectedMetaThreadState,
		types2.ExpectedMetaThread,
		types2.ExpectedMetaClass,
		types2.ExpectedMetaMethod,
		types2.ExpectedMetaPackage,
		types2.ExpectedMetaSymbol,
		types2.ExpectedMetaLogLevel,
		types2.ExpectedMetaStackTrace,
		types2.ExpectedMetaActiveRecording,
		types2.ExpectedMetaActiveSetting,
		types2.ExpectedMetaOSInformation,
		types2.ExpectedMetaJVMInformation,
		types2.ExpectedMetaInitialSystemProperty,
		types2.ExpectedMetaNativeLibrary,
		types2.ExpectedMetaExecutionSample,
		types2.ExpectedMetaObjectAllocationInNewTLAB,
		types2.ExpectedMetaObjectAllocationOutsideTLAB,
		types2.ExpectedMetaJavaMonitorEnter,
		types2.ExpectedMetaThreadPark,
		types2.ExpectedMetaLiveObject,
		types2.ExpectedMetaLog,
		types2.ExpectedMetaCPULoad,
	}
	p.typeExecutionSample = p.typeMap[def.T_EXECUTION_SAMPLE]
	p.typeAllocInNewTLAB = p.typeMap[def.T_ALLOC_IN_NEW_TLAB]
	p.typeALlocOutsideTLAB = p.typeMap[def.T_ALLOC_OUTSIDE_TLAB]
	p.typeMonitorEnter = p.typeMap[def.T_MONITOR_ENTER]
	p.typeThreadPark = p.typeMap[def.T_THREAD_PARK]
	p.typeLiveObject = p.typeMap[def.T_LIVE_OBJECT]
	p.typeActiveSetting = p.typeMap[def.T_ACTIVE_SETTING]
	p.executionSampleHaveContextId = false
	p.allocationInNewTlabHaveContextId = false
	p.allocationOutsideTlabHaveContextId = false
	p.monitorEnterHaveContextId = false
	p.threadParkHaveContextId = false

	for i := range expected {
		realTyp := p.typeMap[expected[i].ID]
		if realTyp == nil {
			continue
		}
		expectedHaveContextID := expected[i].Field("contextId") != nil
		if expectedHaveContextID {

			realHaveContextID := realTyp.Field("contextId") != nil
			if realHaveContextID {
				switch realTyp.ID {
				case def.T_EXECUTION_SAMPLE:
					p.executionSampleHaveContextId = true
				case def.T_ALLOC_IN_NEW_TLAB:
					p.allocationInNewTlabHaveContextId = true
				case def.T_ALLOC_OUTSIDE_TLAB:
					p.allocationOutsideTlabHaveContextId = true
				case def.T_MONITOR_ENTER:
					p.monitorEnterHaveContextId = true
				case def.T_THREAD_PARK:
					p.threadParkHaveContextId = true
				default:
					return fmt.Errorf("unexpected contextId in type %+v", realTyp)
				}
				if !def.CanParse(expected[i].Fields, realTyp.Fields) {
					return fmt.Errorf("unable to parse %+v", realTyp)
				}
			} else {
				if !def.CanParse(expected[i].TrimLastField("contextId"), realTyp.Fields) {
					return fmt.Errorf("unable to parse %+v", realTyp)
				}
			}
		} else {
			if !def.CanParse(expected[i].Fields, realTyp.Fields) {
				return fmt.Errorf("unable to parse %+v", realTyp)
			}
		}
	}

	p.FrameTypes.IDMap = nil
	p.ThreadStates.IDMap = nil
	p.Threads.IDMap = nil
	p.Classes.IDMap = nil
	p.Methods.IDMap.Slice = nil
	p.Methods.IDMap.Dict = nil
	p.Packages.IDMap = nil
	p.Symbols.IDMap = nil
	p.LogLevels.IDMap = nil
	p.Stacktrace.IDMap = nil
	return nil
}
