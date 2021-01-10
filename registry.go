package testparrot

import (
	"errors"
)

var DuplicateErr = errors.New("recorder with same name already exists")

// RecordRegistry registers recorders and provides recorder initialization
type Registry struct {
	recorders   map[string]Recorder
	recordCache map[string][]Record
}

// NewRegistry creates new Registry
func NewRegistry() *Registry {
	return &Registry{
		recorders:   map[string]Recorder{},
		recordCache: map[string][]Record{},
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
	if records, ok := r.recordCache[recorder.Name()]; ok {
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

// Load method loads records into recorder by name
func (r *Registry) Load(name string, records []Record) {
	recorder, ok := r.recorders[name]

	// if recorder is not yet registered, store records in cache
	// so it can be loaded when registering
	if !ok {
		r.recordCache[name] = records
		return
	}

	recorder.Load(records)
}

func (r *Registry) Recorders() []Recorder {
	recorders := []Recorder{}
	for _, recorder := range r.recorders {
		recorders = append(recorders, recorder)
	}

	return recorders
}
