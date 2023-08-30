package parser

import (
	"fmt"

	"github.com/pyroscope-io/jfr-parser/parser/types"
	"github.com/pyroscope-io/jfr-parser/parser/types/def"
)

func (p *Parser) readConstantPool(pos int) error {
	if err := p.seek(pos); err != nil {
		return err
	}
	sz, err := p.varInt()
	if err != nil {
		return err
	}
	typ, err := p.varInt()
	if err != nil {
		return err
	}
	startTimeTicks, err := p.varLong()
	if err != nil {
		return err
	}
	duration, err := p.varInt()
	if err != nil {
		return err
	}
	delta, err := p.varInt()
	if err != nil {
		return err
	}
	typeMask, err := p.varInt()
	if err != nil {
		return err
	}
	n, err := p.varInt()
	if err != nil {
		return err
	}
	//fmt.Printf("ConstantPool: size %d type %d startTimeTicks %d duration %d delta %d typeMask %d n %d\n", sz, typ, startTimeTicks, duration, delta, typeMask, n)
	_ = startTimeTicks
	_ = duration
	_ = delta
	_ = sz

	if typeMask != 1 {
		return fmt.Errorf("expected ConstantPool typeMask 1, got %d", typeMask)
	}
	//if n != 9 {
	//	return fmt.Errorf("expected ConstantPool n 9, got %d", n)
	//}
	for i := 0; i < int(n); i++ {
		typ, err = p.varInt()
		if err != nil {
			return err
		}
		c := p.TypeMap.IDMap[def.TypeID(typ)]
		if c == nil {
			return fmt.Errorf("unknown type %d", def.TypeID(typ))
		}
		err = p.readConstants(c)
		if err != nil {
			return fmt.Errorf("error reading %+v %w", c, err)
		}
	}
	return nil
}

func (p *Parser) readConstants(c *def.Class) error {
	switch c.Name {
	case "jdk.types.FrameType":
		o, err := p.FrameTypes.Parse(p.buf[p.pos:], p.bindFrameType, &p.TypeMap)
		p.pos += o
		return err
	case "jdk.types.ThreadState":
		o, err := p.ThreadStates.Parse(p.buf[p.pos:], p.bindThreadState, &p.TypeMap)
		p.pos += o
		return err
	case "java.lang.Thread":
		o, err := p.Threads.Parse(p.buf[p.pos:], p.bindThread, &p.TypeMap)
		p.pos += o
		return err
	case "java.lang.Class":
		o, err := p.Classes.Parse(p.buf[p.pos:], p.bindClass, &p.TypeMap)
		p.pos += o
		return err
	case "jdk.types.Method":
		o, err := p.Methods.Parse(p.buf[p.pos:], p.bindMethod, &p.TypeMap)
		p.pos += o
		return err
	case "jdk.types.Package":
		o, err := p.Packages.Parse(p.buf[p.pos:], p.bindPackage, &p.TypeMap)
		p.pos += o
		return err
	case "jdk.types.Symbol":
		o, err := p.Symbols.Parse(p.buf[p.pos:], p.bindSymbol, &p.TypeMap)
		p.pos += o
		return err
	case "profiler.types.LogLevel":
		o, err := p.LogLevels.Parse(p.buf[p.pos:], p.bindLogLevel, &p.TypeMap)
		p.pos += o
		return err
	case "jdk.types.StackTrace":
		o, err := p.Stacktrace.Parse(p.buf[p.pos:], p.bindStackTrace, p.bindStackFrame, &p.TypeMap)
		p.pos += o
		return err
	default:
		b := types.NewBindSkipConstantPool(c, &p.TypeMap)
		skipper := types.SkipConstantPoolList{}
		o, err := skipper.Parse(p.buf[p.pos:], b, &p.TypeMap)
		p.pos += o
		return err
	}
}
