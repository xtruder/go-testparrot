package recorder

import (
	"errors"
)

var DuplicateErr = errors.New("recorder with same name already exists")

// RecordRegistry registers recorders and provides recorder initialization
type Registry struct {
	recorders map[string]Recorder
	values    map[string][]Record
}

// NewRegistry creates new Registry
func NewRegistry() *Registry {
	return &Registry{
		recorders: map[string]Recorder{},
		values:    map[string][]Record{},
	}
}

// Register method registers recorder with recoreder registry
func (r *Registry) Register(recorder Recorder) (Recorder, error) {
	recorder.Init(recorder)

	if _, ok := r.recorders[recorder.Name()]; ok {
		return nil, DuplicateErr
	}

	r.recorders[recorder.Name()] = recorder

	// if values have already been provided for a speicfic recorder, load
	// them into recorder
	if records, ok := r.values[recorder.Name()]; ok {
		recorder.Load(records)
	}

	return recorder, nil
}

func (r *Registry) MustRegister(recorder Recorder) Recorder {
	recorder, err := r.Register(recorder)
	if err != nil {
		panic(err)
	}

	return recorder
}

// EnableRecording enables or disables recording on all recorders
func (r *Registry) EnableRecording(enable bool) {
	for _, r := range r.recorders {
		r.EnableRecording(enable)
	}
}

func (r *Registry) Load(values map[string][]Record) {
	r.values = values

	// load records for all recorders that have already been registered
	for name, records := range r.values {
		if recorder, ok := r.recorders[name]; ok {
			recorder.Load(records)
		}
	}
}
