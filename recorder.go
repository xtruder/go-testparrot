package testparrot

import (
	"fmt"
	"path"
	"testing"
)

type Recording struct {
	// Key defines record key
	Key interface{}

	// Value defines record value
	Value interface{}
}

type Recorder struct {
	allRecordings map[string][]Recording

	// counter defines counter for sequential recordings
	counters map[string]int

	// testFilenames where individual tests are
	testFilenames map[string]string

	// recordingEnabled defines whether recording is enabled
	recordingEnabled bool
}

// NewRecoder creates a new Recorder
func NewRecorder() *Recorder {
	return &Recorder{
		allRecordings: map[string][]Recording{},
		counters:      map[string]int{},
		testFilenames: map[string]string{},
	}
}

// Reset method resets recorder
func (r *Recorder) Reset() {
	r.allRecordings = map[string][]Recording{}
	r.counters = map[string]int{}
	r.testFilenames = map[string]string{}
}

// Recorder method records value under specified key. If recording is enabled
// Record returns provided value, otherwise it returns already recorded value.
func (r *Recorder) Record(t *testing.T, key interface{}, value interface{}) interface{} {
	name := t.Name()

	testPath, err := getTestPath(t)
	if err != nil {
		panic(err)
	}

	r.testFilenames[name] = path.Base(testPath)

	value, err = r.record(name, key, value)
	if err != nil {
		panic(err)
	}

	return value
}

// RecordNext method records next value in sequence. If recording is enabled
// Record returns provided value, otherwise it returns alreday recorded value.
func (r *Recorder) RecordNext(t *testing.T, value interface{}) interface{} {
	name := t.Name()

	testPath, err := getTestPath(t)
	if err != nil {
		panic(err)
	}

	r.testFilenames[name] = path.Base(testPath)

	if _, ok := r.counters[name]; !ok {
		r.counters[name] = 0
	}

	value, err = r.record(name, r.counters[name], value)
	if err != nil {
		panic(err)
	}

	// increase counter
	r.counters[name]++

	return value
}

// EnableRecording enables test recording
func (r *Recorder) EnableRecording(enable bool) {
	r.recordingEnabled = enable
}

// RecordingEnabled returns whether recording is enabled
func (r *Recorder) RecordingEnabled() bool {
	return r.recordingEnabled
}

// Load method loads recording for a specific test name
func (r *Recorder) Load(name string, recordings []Recording) {
	if _, ok := r.allRecordings[name]; ok {
		panic(newErr(fmt.Errorf("recordings already loaded for test '%s'", name)))
	}

	r.allRecordings[name] = recordings
}

func (r *Recorder) record(name string, key interface{}, value interface{}) (interface{}, error) {
	if !r.recordingEnabled {
		value, err := r.getRecordValue(name, key)
		if err != nil {
			return nil, err
		}

		return value, nil
	}

	if err := r.setRecordValue(name, key, value); err != nil {
		return nil, err
	}

	return value, nil
}

func (r *Recorder) getRecordValue(name string, key interface{}) (interface{}, error) {
	if records, ok := r.allRecordings[name]; ok {
		for _, record := range records {
			if record.Key == key {
				return record.Value, nil
			}
		}
	}

	return nil, newErr(fmt.Errorf("recording with key '%v' not found for test '%s'", key, name))
}

func (r *Recorder) setRecordValue(name string, key interface{}, value interface{}) error {
	if records, ok := r.allRecordings[name]; ok {
		for _, record := range records {
			if record.Key == key {
				return fmt.Errorf("recording with key '%v' already exists for test '%s'", key, name)
			}
		}

		r.allRecordings[name] = append(records, Recording{key, value})
	} else {
		r.allRecordings[name] = []Recording{{key, value}}
	}

	return nil
}
