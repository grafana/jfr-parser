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
	if def.TypeID(typ) != def.T_CPOOL {
		return fmt.Errorf("expected ConstantPool type %d, got %d", def.T_CPOOL, typ)
	}
	if typeMask != 1 {
		return fmt.Errorf("expected ConstantPool typeMask 1, got %d", typeMask)
	}
	if n != 9 {
		return fmt.Errorf("expected ConstantPool n 9, got %d", n)
	}
	for i := 0; i < int(n); i++ {
		typ, err = p.varInt()
		if err != nil {
			return err
		}
		c := p.typeMap[def.TypeID(typ)]
		if c == nil {
			return fmt.Errorf("unknown type %s", def.TypeID2Sym(def.TypeID(typ)))
		}
		err = p.readConstants(c)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) readConstants(c *def.Class) error {
	switch c.ID {
	case def.T_FRAME_TYPE:
		o, err := p.FrameTypes.Parse(p.buf[p.pos:], c, p.typeMap)
		//o, err := types.Skip(p.buf[p.pos:], p.typeMap[c.ID], p.typeMap, true)
		p.pos += o
		return err
	case def.T_THREAD_STATE:
		o, err := p.ThreadStates.Parse(p.buf[p.pos:], c, p.typeMap)
		p.pos += o
		return err
	case def.T_THREAD:
		o, err := p.Threads.Parse(p.buf[p.pos:], c, p.typeMap)
		p.pos += o
		return err
	case def.T_CLASS:
		o, err := p.Classes.Parse(p.buf[p.pos:], c, p.typeMap)
		p.pos += o
		return err
	case def.T_METHOD:
		o, err := p.Methods.Parse(p.buf[p.pos:], c, p.typeMap)
		p.pos += o
		return err
	case def.T_PACKAGE:
		o, err := p.Packages.Parse(p.buf[p.pos:], c, p.typeMap)
		p.pos += o
		return err
	case def.T_SYMBOL:
		o, err := p.Symbols.Parse(p.buf[p.pos:], c, p.typeMap)
		p.pos += o
		return err
	case def.T_LOG_LEVEL:
		//o, err := types.Skip(p.buf[p.pos:], p.typeMap[c.ID], p.typeMap, true)
		o, err := p.LogLevels.Parse(p.buf[p.pos:], c, p.typeMap)
		p.pos += o
		return err
	case def.T_STACK_TRACE:
		sft := p.typeMap[def.T_STACK_FRAME]
		if sft == nil {
			return fmt.Errorf("unknown type %s", def.TypeID2Sym(def.T_STACK_FRAME))
		}
		o, err := p.Stacktrace.Parse(p.buf[p.pos:], c, sft, p.typeMap)
		p.pos += o
		return err
	default:
		//todo test wtih above
		o, err := types.Skip(p.buf[p.pos:], c, p.typeMap, true)
		p.pos += o
		return err
	}
}
