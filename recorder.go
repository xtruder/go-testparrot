package testparrot

import "fmt"

type Record struct {
	// Name defines name of record
	Name string

	// Value defines recorded value
	Value interface{}
}

type RecordProvider interface {
	Records() []Record
}

type Recorder interface {
	// Init method initializes recorder by passing itself and recorder options
	Init(self Recorder)

	// Name method gets recorder name
	Name() string

	// Records returns list of recorded records
	Records() []Record

	// Load method loads list of records into recorder
	Load([]Record)

	// EnableRecording method enables or disables recording on recorder
	EnableRecording(bool)
}

type baseRecorder struct {
	name            string
	enableRecording bool
}

func (r *baseRecorder) Name() string {
	return r.name
}

func (r *baseRecorder) EnableRecording(enable bool) {
	r.enableRecording = enable
}

type KVRecorder struct {
	baseRecorder

	values map[string]interface{}
}

func NewKVRecorder(name string) *KVRecorder {
	return &KVRecorder{baseRecorder: baseRecorder{name: name}}
}

func (r *KVRecorder) Init(self Recorder) {
	r.values = make(map[string]interface{})
}

func (r *KVRecorder) Record(name string, val interface{}) interface{} {
	if r.enableRecording {
		r.values[name] = val
	} else if val, ok := r.values[name]; ok {
		return val
	} else {
		panic(fmt.Errorf("go-test-record: value '%s' not recorded in recorder '%s'", name, r.name))
	}

	return val
}

func (r *KVRecorder) Value(name string) interface{} {
	return r.values[name]
}

// Records method returns recroded records
func (r *KVRecorder) Records() []Record {
	records := []Record{}

	for n, v := range r.values {
		records = append(records, Record{n, v})
	}

	return records
}

// Load method loads records into recorder, does not load it in struct
// yet, because you need to use Record method for that
func (r *KVRecorder) Load(records []Record) {
	for _, v := range records {
		r.values[v.Name] = v.Value
	}
}
