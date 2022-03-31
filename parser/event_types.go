package parser

import (
	"fmt"

	"github.com/pyroscope-io/jfr-parser/reader"
)

var events = map[string]func() Parseable{
	"jdk.ActiveRecording":       func() Parseable { return new(ActiveRecording) },
	"jdk.ActiveSetting":         func() Parseable { return new(ActiveSetting) },
	"jdk.CPUInformation":        func() Parseable { return new(CPUInformation) },
	"jdk.CPULoad":               func() Parseable { return new(CPULoad) },
	"jdk.ExecutionSample":       func() Parseable { return new(ExecutionSample) },
	"jdk.InitialSystemProperty": func() Parseable { return new(InitialSystemProperty) },
	// TODO: jdk.JavaMonitorEnter
	"jdk.JVMInformation":              func() Parseable { return new(JVMInformation) },
	"jdk.ObjectAllocationInNewTLAB":   func() Parseable { return new(ObjectAllocationInNewTLAB) },
	"jdk.ObjectAllocationOutsideTLAB": func() Parseable { return new(ObjectAllocationOutsideTLAB) },
	"jdk.OSInformation":               func() Parseable { return new(OSInformation) },
	// TODO: jdk.ThreadPark
}

func ParseEvent(r reader.Reader, classes ClassMap, cpools PoolMap) (Parseable, error) {
	kind, err := r.VarLong()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve event type: %w", err)
	}
	return parseEvent(r, classes, cpools, int(kind))
}

func parseEvent(r reader.Reader, classes ClassMap, cpools PoolMap, classID int) (Parseable, error) {
	class, ok := classes[classID]
	if !ok {
		return nil, fmt.Errorf("unknown class %d", classID)
	}
	var v Parseable
	if typeFn, ok := events[class.Name]; ok {
		v = typeFn()
	} else {
		v = new(UnsupportedEvent)
	}
	if err := v.Parse(r, classes, cpools, class); err != nil {
		return nil, fmt.Errorf("unable to parse event %s: %w", class.Name, err)
	}
	return v, nil
}

type ActiveRecording struct {
	StartTime         int64
	Duration          int64
	EventThread       *Thread
	ID                int64
	Name              string
	Destination       string
	MaxAge            int64
	MaxSize           int64
	RecordingStart    int64
	RecordingDuration int64
}

func (ar *ActiveRecording) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		ar.StartTime, err = toLong(p)
	case "duration":
		ar.Duration, err = toLong(p)
	case "eventThread":
		ar.EventThread, err = toThread(p)
	case "id":
		ar.ID, err = toLong(p)
	case "name":
		ar.Name, err = toString(p)
	case "destination":
		ar.Destination, err = toString(p)
	case "maxAge":
		ar.MaxAge, err = toLong(p)
	case "maxSize":
		ar.MaxSize, err = toLong(p)
	case "recordingStart":
		ar.RecordingStart, err = toLong(p)
	case "recordingDuration":
		ar.RecordingDuration, err = toLong(p)
	}
	return err
}

func (ar *ActiveRecording) Parse(r reader.Reader, classes ClassMap, cpools PoolMap, class ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, ar.parseField)
}

type ActiveSetting struct {
	StartTime   int64
	Duration    int64
	EventThread *Thread
	ID          int64
	Name        string
	Value       string
}

func (as *ActiveSetting) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		as.StartTime, err = toLong(p)
	case "duration":
		as.Duration, err = toLong(p)
	case "eventThread":
		as.EventThread, err = toThread(p)
	case "id":
		as.ID, err = toLong(p)
	case "name":
		as.Name, err = toString(p)
	case "value":
		as.Value, err = toString(p)
	}
	return err
}

func (as *ActiveSetting) Parse(r reader.Reader, classes ClassMap, cpools PoolMap, class ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, as.parseField)
}

type CPUInformation struct {
	StartTime   int64
	CPU         string
	Description string
	Sockets     int32
	Cores       int32
	HWThreads   int32
}

func (ci *CPUInformation) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		ci.StartTime, err = toLong(p)
	case "duration":
		ci.CPU, err = toString(p)
	case "eventThread":
		ci.Description, err = toString(p)
	case "sockets":
		ci.Sockets, err = toInt(p)
	case "cores":
		ci.Cores, err = toInt(p)
	case "hwThreads":
		ci.HWThreads, err = toInt(p)
	}
	return err
}

func (ci *CPUInformation) Parse(r reader.Reader, classes ClassMap, cpools PoolMap, class ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, ci.parseField)
}

type CPULoad struct {
	StartTime    int64
	JVMUser      float32
	JVMSystem    float32
	MachineTotal float32
}

func (cl *CPULoad) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		cl.StartTime, err = toLong(p)
	case "jvmUser":
		cl.JVMUser, err = toFloat(p)
	case "jvmSystem":
		cl.JVMSystem, err = toFloat(p)
	case "machineTotal":
		cl.MachineTotal, err = toFloat(p)
	}
	return err
}

func (cl *CPULoad) Parse(r reader.Reader, classes ClassMap, cpools PoolMap, class ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, cl.parseField)
}

type ExecutionSample struct {
	StartTime     int64
	SampledThread *Thread
	StackTrace    *StackTrace
	State         *ThreadState
}

func (es *ExecutionSample) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		es.StartTime, err = toLong(p)
	case "sampledThread":
		es.SampledThread, err = toThread(p)
	case "stackTrace":
		es.StackTrace, err = toStackTrace(p)
	case "machineTotal":
		es.State, err = toThreadState(p)
	}
	return err
}

func (es *ExecutionSample) Parse(r reader.Reader, classes ClassMap, cpools PoolMap, class ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, es.parseField)
}

type InitialSystemProperty struct {
	StartTime int64
	Key       string
	Value     string
}

func (isp *InitialSystemProperty) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		isp.StartTime, err = toLong(p)
	case "key":
		isp.Key, err = toString(p)
	case "stackTrace":
		isp.Value, err = toString(p)
	}
	return err
}

func (isp *InitialSystemProperty) Parse(r reader.Reader, classes ClassMap, cpools PoolMap, class ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, isp.parseField)
}

type JVMInformation struct {
	StartTime     int64
	JVMName       string
	JVMVersion    string
	JVMArguments  string
	JVMFlags      string
	JavaArguments string
	JVMStartTime  int64
	PID           int64
}

func (ji *JVMInformation) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		ji.StartTime, err = toLong(p)
	case "jvmName":
		ji.JVMName, err = toString(p)
	case "jvmVersion":
		ji.JVMVersion, err = toString(p)
	case "jvmArguments":
		ji.JVMArguments, err = toString(p)
	case "jvmFlags":
		ji.JVMFlags, err = toString(p)
	case "javaArguments":
		ji.JavaArguments, err = toString(p)
	case "jvmStartTime":
		ji.JVMStartTime, err = toLong(p)
	case "pid":
		ji.PID, err = toLong(p)
	}
	return err
}

func (ji *JVMInformation) Parse(r reader.Reader, classes ClassMap, cpools PoolMap, class ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, ji.parseField)
}

type ObjectAllocationInNewTLAB struct {
	StartTime      int64
	EventThread    *Thread
	StackTrace     *StackTrace
	ObjectClass    *Class
	AllocationSize int64
	TLABSize       int64
}

func (oa *ObjectAllocationInNewTLAB) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		oa.StartTime, err = toLong(p)
	case "sampledThread":
		oa.EventThread, err = toThread(p)
	case "stackTrace":
		oa.StackTrace, err = toStackTrace(p)
	case "objectClass":
		oa.ObjectClass, err = toClass(p)
	case "allocationSize":
		oa.AllocationSize, err = toLong(p)
	case "tlabSize":
		oa.TLABSize, err = toLong(p)
	}
	return err
}

func (oa *ObjectAllocationInNewTLAB) Parse(r reader.Reader, classes ClassMap, cpools PoolMap, class ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, oa.parseField)
}

type ObjectAllocationOutsideTLAB struct {
	StartTime      int64
	EventThread    *Thread
	StackTrace     *StackTrace
	ObjectClass    *Class
	AllocationSize int64
}

func (oa *ObjectAllocationOutsideTLAB) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		oa.StartTime, err = toLong(p)
	case "sampledThread":
		oa.EventThread, err = toThread(p)
	case "stackTrace":
		oa.StackTrace, err = toStackTrace(p)
	case "objectClass":
		oa.ObjectClass, err = toClass(p)
	case "allocationSize":
		oa.AllocationSize, err = toLong(p)
	}
	return err
}

func (oa *ObjectAllocationOutsideTLAB) Parse(r reader.Reader, classes ClassMap, cpools PoolMap, class ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, oa.parseField)
}

type OSInformation struct {
	StartTime int64
	OSVersion string
}

func (os *OSInformation) parseField(name string, p ParseResolvable) (err error) {
	switch name {
	case "startTime":
		os.StartTime, err = toLong(p)
	case "osVersion":
		os.OSVersion, err = toString(p)
	}
	return err
}

func (os *OSInformation) Parse(r reader.Reader, classes ClassMap, cpools PoolMap, class ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, os.parseField)
}

type UnsupportedEvent struct{}

func (ue *UnsupportedEvent) parseField(name string, p ParseResolvable) error {
	return nil
}

func (ue *UnsupportedEvent) Parse(r reader.Reader, classes ClassMap, cpools PoolMap, class ClassMetadata) error {
	return parseFields(r, classes, cpools, class, nil, true, ue.parseField)
}
