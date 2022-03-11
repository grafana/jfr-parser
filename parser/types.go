package parser

import (
	"errors"
	"fmt"

	"github.com/pyroscope-io/jfr-parser/reader"
)

var types = map[string]func() ParseResolvable{
	"boolean": func() ParseResolvable { return new(Boolean) },
	// TODO: char
	"byte":  func() ParseResolvable { return new(Byte) },
	"short": func() ParseResolvable { return new(Short) },
	"int":   func() ParseResolvable { return new(Int) },
	"long":  func() ParseResolvable { return new(Long) },
	"float": func() ParseResolvable { return new(Float) },
	// TODO: double
	"java.lang.Class":       func() ParseResolvable { return new(Class) },
	"java.lang.String":      func() ParseResolvable { return new(String) },
	"java.lang.Thread":      func() ParseResolvable { return new(Thread) },
	"jdk.types.ClassLoader": func() ParseResolvable { return new(ClassLoader) },
	"jdk.types.FrameType":   func() ParseResolvable { return new(FrameType) },
	"jdk.types.Method":      func() ParseResolvable { return new(Method) },
	"jdk.types.Package":     func() ParseResolvable { return new(Package) },
	"jdk.types.StackFrame":  func() ParseResolvable { return new(StackFrame) },
	"jdk.types.StackTrace":  func() ParseResolvable { return new(StackTrace) },
	"jdk.types.Symbol":      func() ParseResolvable { return new(Symbol) },
	"jdk.types.ThreadState": func() ParseResolvable { return new(ThreadState) },
}

func ParseClass(r reader.Reader, classes ClassMap, cpools PoolMap, classID int64) (ParseResolvable, error) {
	class, ok := classes[int(classID)]
	if !ok {
		return nil, fmt.Errorf("unexpected class %d", classID)
	}
	typeFn, ok := types[class.Name]
	if !ok {
		return nil, fmt.Errorf("unknown type %s", class.Name)
	}
	v := typeFn()
	if err := v.Parse(r, classes, cpools, class); err != nil {
		return nil, err
	}
	return v, nil
}

type Parseable interface {
	Parse(reader.Reader, ClassMap, PoolMap, ClassMetadata) error
}

type Resolvable interface {
	Resolve(ClassMap, PoolMap) error
}

type ParseResolvable interface {
	Parseable
	Resolvable
}

type constant struct {
	classID int64
	field   string
	index   int64
}

func appendConstant(r reader.Reader, constants *[]constant, name string) error {
	i, err := r.VarLong()
	if err != nil {
		return fmt.Errorf("unable to read constant index")
	}
	*constants = append(*constants, constant{field: name, index: i})
	return nil
}

func parseFields(r reader.Reader, classes ClassMap, cpools PoolMap, class ClassMetadata, constants *[]constant, cb func(string, ParseResolvable) error) error {
	for _, f := range class.Fields {
		if f.ConstantPool {
			if constants != nil {
				if err := appendConstant(r, constants, f.Name); err != nil {
					return fmt.Errorf("failed to parse %s: unable to append constant: %w", class.Name, err)
				}
			} else {
				cpool, ok := cpools[int(f.Class)]
				if !ok {
					return fmt.Errorf("unknown constant pool class %d", f.Class)
				}
				i, err := r.VarLong()
				if err != nil {
					return fmt.Errorf("unable to read constant index")
				}
				p, ok := cpool[int(i)]
				if !ok {
					return fmt.Errorf("unknown constant pool index %d for class %d", i, f.Class)
				}
				if err := cb(f.Name, p); err != nil {
					return fmt.Errorf("unable to parse constant field %s: %w", f.Name, err)
				}
			}
		} else if f.Dimension == 1 {
			n, err := r.VarInt()
			if err != nil {
				return fmt.Errorf("failed to parse %s: unable to read array length: %w", class.Name, err)
			}
			// TODO: assert n is small enough
			for i := 0; i < int(n); i++ {
				p, err := ParseClass(r, classes, cpools, f.Class)
				if err != nil {
					return fmt.Errorf("failed to parse %s: unable to read an array element: %w", class.Name, err)
				}
				if err := cb(f.Name, p); err != nil {
					return fmt.Errorf("failed to parse %s: unable to parse an array element: %w", class.Name, err)
				}
			}
		} else {
			p, err := ParseClass(r, classes, cpools, f.Class)
			if err != nil {
				return fmt.Errorf("failed to parse %s: unable to read a field: %w", class.Name, err)
			}
			if err := cb(f.Name, p); err != nil {
				return fmt.Errorf("failed to parse %s: unable to parse a field: %w", class.Name, err)
			}
		}
	}
	return nil
}

func resolveConstants(classes ClassMap, cpools PoolMap, constants *[]constant, cb func(string, ParseResolvable) error) error {
	for _, c := range *constants {
		if err := ResolveConstants(classes, cpools, int(c.classID)); err != nil {
			return fmt.Errorf("unable to resolve contants: %w", err)
		}
		p, ok := cpools[int(c.classID)][int(c.index)]
		if !ok {
			return fmt.Errorf("unknown constant pool %d for class %d", c.index, c.classID)
		}
		if err := cb(c.field, p); err != nil {
			return fmt.Errorf("unable to resolve constants for field %s: %w", c.field, c.classID)
		}
	}
	*constants = nil
	return nil
}

type Boolean bool

func (b *Boolean) Parse(r reader.Reader, _ ClassMap, _ PoolMap, _ ClassMetadata) error {
	// TODO: Assert simpletype, no fields, etc.
	x, err := r.Boolean()
	*b = Boolean(x)
	return err
}

func (Boolean) Resolve(ClassMap, PoolMap) error { return nil }

func toBoolean(p Parseable) (bool, error) {
	x, ok := p.(*Boolean)
	if !ok {
		return false, errors.New("not a Boolean")
	}
	return bool(*x), nil
}

type Byte int8

func (b *Byte) Parse(r reader.Reader, _ ClassMap, _ PoolMap, _ ClassMetadata) error {
	x, err := r.Byte()
	*b = Byte(x)
	return err
}

func (Byte) Resolve(ClassMap, PoolMap) error { return nil }

type Short int16

func (s *Short) Parse(r reader.Reader, _ ClassMap, _ PoolMap, _ ClassMetadata) error {
	x, err := r.VarShort()
	*s = Short(x)
	return err
}

func (Short) Resolve(ClassMap, PoolMap) error { return nil }

type Int int32

func (i *Int) Parse(r reader.Reader, _ ClassMap, _ PoolMap, _ ClassMetadata) error {
	x, err := r.VarInt()
	*i = Int(x)
	return err
}

func (Int) Resolve(ClassMap, PoolMap) error { return nil }

func toInt(p Parseable) (int32, error) {
	x, ok := p.(*Int)
	if !ok {
		return 0, errors.New("not an Int")
	}
	return int32(*x), nil
}

type Long int64

func (l *Long) Parse(r reader.Reader, _ ClassMap, _ PoolMap, _ ClassMetadata) error {
	x, err := r.VarLong()
	*l = Long(x)
	return err
}

func (Long) Resolve(ClassMap, PoolMap) error { return nil }

func toLong(p Parseable) (int64, error) {
	x, ok := p.(*Long)
	if !ok {
		return 0, errors.New("not a Long")
	}
	return int64(*x), nil
}

type Float float32

func (f *Float) Parse(r reader.Reader, _ ClassMap, _ PoolMap, _ ClassMetadata) error {
	x, err := r.Float()
	*f = Float(x)
	return err
}

func (Float) Resolve(ClassMap, PoolMap) error { return nil }

func toFloat(p Parseable) (float32, error) {
	x, ok := p.(*Float)
	if !ok {
		return 0, errors.New("not a Float")
	}
	return float32(*x), nil
}

// TODO: rest of builtin types

type String string

func (s *String) Parse(r reader.Reader, _ ClassMap, _ PoolMap, _ ClassMetadata) error {
	x, err := r.String()
	*s = String(x)
	return err
}

func (s String) Resolve(_ ClassMap, _ PoolMap) error { return nil }

func toString(p Parseable) (string, error) {
	s, ok := p.(*String)
	if !ok {
		return "", errors.New("not a String")
	}
	return string(*s), nil
}

type FrameType struct {
	Description string
	constants   []constant
}

func (ft *FrameType) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "description":
		ft.Description, err = toString(p)
	}
	return err
}

func (ft *FrameType) Parse(r reader.Reader, classes ClassMap, cpools PoolMap, class ClassMetadata) error {
	return parseFields(r, classes, cpools, class, &ft.constants, ft.parseField)
}

func (ft *FrameType) Resolve(classes ClassMap, cpools PoolMap) error {
	return resolveConstants(classes, cpools, &ft.constants, ft.parseField)
}

func toFrameType(p Parseable) (*FrameType, error) {
	ft, ok := p.(*FrameType)
	if !ok {
		return nil, errors.New("not a FrameType")
	}
	return ft, nil
}

type ThreadState struct {
	Name      string
	constants []constant
}

func (ts *ThreadState) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "name":
		ts.Name, err = toString(p)
	}
	return err
}

func (ts *ThreadState) Parse(r reader.Reader, classes ClassMap, cpools PoolMap, class ClassMetadata) error {
	return parseFields(r, classes, cpools, class, &ts.constants, ts.parseField)
}

func (ts *ThreadState) Resolve(classes ClassMap, cpools PoolMap) error {
	return resolveConstants(classes, cpools, &ts.constants, ts.parseField)
}

func toThreadState(p ParseResolvable) (*ThreadState, error) {
	ts, ok := p.(*ThreadState)
	if !ok {
		return nil, errors.New("not a ThreadState")
	}
	return ts, nil
}

type Thread struct {
	OsName       string
	OsThreadID   int64
	JavaName     string
	JavaThreadID int64
	constants    []constant
}

func (t *Thread) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "osName":
		t.OsName, err = toString(p)
	case "osThreadId":
		t.OsThreadID, err = toLong(p)
	case "javaName":
		t.JavaName, err = toString(p)
	case "javaThreadId":
		t.JavaThreadID, err = toLong(p)
	}
	return err
}

func (t *Thread) Parse(r reader.Reader, classes ClassMap, cpools PoolMap, class ClassMetadata) error {
	return parseFields(r, classes, cpools, class, &t.constants, t.parseField)
}

func (t *Thread) Resolve(classes ClassMap, cpools PoolMap) error {
	return resolveConstants(classes, cpools, &t.constants, t.parseField)
}

func toThread(p ParseResolvable) (*Thread, error) {
	t, ok := p.(*Thread)
	if !ok {
		return nil, errors.New("not a Thread")
	}
	return t, nil
}

type StackFrame struct {
	Method        *Method
	LineNumber    int32
	ByteCodeIndex int32
	Type          *FrameType
	constants     []constant
}

func (sf *StackFrame) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "method":
		sf.Method, err = toMethod(p)
	case "lineNumber":
		sf.LineNumber, err = toInt(p)
	case "byteCodeIndex":
		sf.ByteCodeIndex, err = toInt(p)
	case "type":
		sf.Type, err = toFrameType(p)
	}
	return err
}

func (sf *StackFrame) Parse(r reader.Reader, classes ClassMap, cpools PoolMap, class ClassMetadata) error {
	return parseFields(r, classes, cpools, class, &sf.constants, sf.parseField)
}

func (sf *StackFrame) Resolve(classes ClassMap, cpools PoolMap) error {
	return resolveConstants(classes, cpools, &sf.constants, sf.parseField)
}

func toStackFrame(p ParseResolvable) (*StackFrame, error) {
	sf, ok := p.(*StackFrame)
	if !ok {
		return nil, errors.New("not a StackFrame")
	}
	return sf, nil
}

type StackTrace struct {
	Truncated bool
	Frames    []*StackFrame
	constants []constant
}

func (st *StackTrace) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "truncated":
		st.Truncated, err = toBoolean(p)
	case "frames":
		var sf *StackFrame
		sf, err := toStackFrame(p)
		if err != nil {
			return err
		}
		st.Frames = append(st.Frames, sf)
	}
	return err
}

func (st *StackTrace) Parse(r reader.Reader, classes ClassMap, cpools PoolMap, class ClassMetadata) error {
	return parseFields(r, classes, cpools, class, &st.constants, st.parseField)
}

func (st *StackTrace) Resolve(classes ClassMap, cpools PoolMap) error {
	return resolveConstants(classes, cpools, &st.constants, st.parseField)
}

func toStackTrace(p ParseResolvable) (*StackTrace, error) {
	st, ok := p.(*StackTrace)
	if !ok {
		return nil, errors.New("not a StackTrace")
	}
	return st, nil
}

type Method struct {
	Type       *Class
	Name       *Symbol
	Descriptor *Symbol
	Modifiers  int32
	Hidden     bool
	constants  []constant
}

func (m *Method) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "type":
		m.Type, err = toClass(p)
	case "name":
		m.Name, err = toSymbol(p)
	case "descriptor":
		m.Descriptor, err = toSymbol(p)
	case "modifiers":
		m.Modifiers, err = toInt(p)
	case "hidden":
		m.Hidden, err = toBoolean(p)
	}
	return err
}

func (m *Method) Parse(r reader.Reader, classes ClassMap, cpools PoolMap, class ClassMetadata) error {
	return parseFields(r, classes, cpools, class, &m.constants, m.parseField)
}

func (m *Method) Resolve(classes ClassMap, cpools PoolMap) error {
	return resolveConstants(classes, cpools, &m.constants, m.parseField)
}

func toMethod(p ParseResolvable) (*Method, error) {
	m, ok := p.(*Method)
	if !ok {
		return nil, errors.New("not a Method")
	}
	return m, nil
}

type Class struct {
	ClassLoader *ClassLoader
	Name        *Symbol
	Package     *Package
	Modifiers   int64
	constants   []constant
}

func (c *Class) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "classLoader":
		c.ClassLoader, err = toClassLoader(p)
	case "name":
		c.Name, err = toSymbol(p)
	case "package":
		c.Package, err = toPackage(p)
	case "modifers":
		c.Modifiers, err = toLong(p)
	}
	return err
}

func (c *Class) Parse(r reader.Reader, classes ClassMap, cpools PoolMap, class ClassMetadata) error {
	return parseFields(r, classes, cpools, class, &c.constants, c.parseField)
}

func (c *Class) Resolve(classes ClassMap, cpools PoolMap) error {
	return resolveConstants(classes, cpools, &c.constants, c.parseField)
}

func toClass(p ParseResolvable) (*Class, error) {
	c, ok := p.(*Class)
	if !ok {
		// TODO
		return nil, errors.New("")
	}
	return c, nil
}

type ClassLoader struct {
	Type      *Class
	Name      *Symbol
	constants []constant
}

func (cl *ClassLoader) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "type":
		cl.Type, err = toClass(p)
	case "name":
		cl.Name, err = toSymbol(p)
	}
	return err
}

func (cl *ClassLoader) Parse(r reader.Reader, classes ClassMap, cpools PoolMap, class ClassMetadata) error {
	return parseFields(r, classes, cpools, class, &cl.constants, cl.parseField)
}

func (cl *ClassLoader) Resolve(classes ClassMap, cpools PoolMap) error {
	return resolveConstants(classes, cpools, &cl.constants, cl.parseField)
}

func toClassLoader(p ParseResolvable) (*ClassLoader, error) {
	c, ok := p.(*ClassLoader)
	if !ok {
		// TODO
		return nil, errors.New("")
	}
	return c, nil
}

type Package struct {
	Name      *Symbol
	constants []constant
}

func (pkg *Package) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "name":
		pkg.Name, err = toSymbol(p)
	}
	return err
}

func (p *Package) Parse(r reader.Reader, classes ClassMap, cpools PoolMap, class ClassMetadata) error {
	return parseFields(r, classes, cpools, class, &p.constants, p.parseField)
}

func (p *Package) Resolve(classes ClassMap, cpools PoolMap) error {
	return resolveConstants(classes, cpools, &p.constants, p.parseField)
}

func toPackage(p ParseResolvable) (*Package, error) {
	pkg, ok := p.(*Package)
	if !ok {
		// TODO
		return nil, errors.New("")
	}
	return pkg, nil
}

type Symbol struct {
	String    string
	constants []constant
}

func (s *Symbol) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "string":
		s.String, err = toString(p)
	}
	return err
}

func (s *Symbol) Parse(r reader.Reader, classes ClassMap, cpools PoolMap, class ClassMetadata) error {
	return parseFields(r, classes, cpools, class, &s.constants, s.parseField)
}

func (s *Symbol) Resolve(classes ClassMap, cpools PoolMap) error {
	return resolveConstants(classes, cpools, &s.constants, s.parseField)
}

func toSymbol(p ParseResolvable) (*Symbol, error) {
	s, ok := p.(*Symbol)
	if !ok {
		// TODO
		return nil, errors.New("")
	}
	return s, nil
}