package testparrot

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewRegistry(t *testing.T) {
	r := NewRegistry()
	require.IsType(t, &Registry{}, r)
}

func TestRegistryRegister(t *testing.T) {
	recorder := NewKVRecorder("test")

	r := NewRegistry()
	registered, err := r.Register(recorder)
	require.NoError(t, err)

	require.Equal(t, recorder, registered)
	require.Contains(t, r.recorders, "test")
	require.Equal(t, recorder, r.recorders["test"])
}

func TestRegistryShouldNotRegisterWithSameName(t *testing.T) {
	r := NewRegistry()
	_, err := r.Register(NewKVRecorder("recorder"))
	require.NoError(t, err)

	_, err = r.Register(NewKVRecorder("recorder"))
	require.EqualError(t, err, DuplicateErr.Error())
}

func TestRegistryLoadShouldStoreInCache(t *testing.T) {
	records := []Record{{"key", "value"}}

	r := NewRegistry()
	r.Load("recorder", records)

	require.Contains(t, r.recordCache, "recorder")
	require.Equal(t, records, r.recordCache["recorder"])
}

func TestRegistryLoadShouldStoreInRecorder(t *testing.T) {
	recorder := NewKVRecorder("recorder")

	r := NewRegistry()
	r.MustRegister(recorder)

	values := []Record{{"key", "value"}}

	r.Load("recorder", values)

	require.Equal(t, values, recorder.Records())
}

func TestRegistryRegisterShouldLoadIntoRecorder(t *testing.T) {
	r := NewRegistry()

	values := []Record{{"key", "value"}}

	r.Load("recorder", values)

	recorder := NewKVRecorder("recorder")
	r.MustRegister(recorder)

	require.Equal(t, values, recorder.Records())
}

func TestEnableRecording(t *testing.T) {
	r := NewRegistry()
	recorder := NewKVRecorder("test")
	r.MustRegister(recorder)

	r.EnableRecording(true)

	require.Equal(t, recorder.enableRecording, true)
}
