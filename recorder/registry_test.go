package recorder

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
	registered := r.MustRegister(recorder)

	require.Equal(t, recorder, registered)
	require.Contains(t, r.recorders, "test")
	require.Equal(t, recorder, r.recorders["test"])
}

func TestRegistryLoad(t *testing.T) {
	values := map[string][]Record{
		"test": {{"key", "value"}},
	}

	r := NewRegistry()
	r.Load(values)

	require.Equal(t, values, r.values)
}

func TestRegistryLoadShouldLoadIntoRecorder(t *testing.T) {
	recorder := NewKVRecorder("test")

	r := NewRegistry()
	r.MustRegister(recorder)

	values := map[string][]Record{
		"test": {{"key", "value"}},
	}

	r.Load(values)

	require.Equal(t, values["test"], recorder.Records())
}

func TestRegistryRegisterShouldLoadIntoRecorder(t *testing.T) {
	r := NewRegistry()

	values := map[string][]Record{
		"test": {{"key", "value"}},
	}

	r.Load(values)

	recorder := NewKVRecorder("test")
	r.MustRegister(recorder)

	require.Equal(t, values["test"], recorder.Records())
}

func TestEnableRecording(t *testing.T) {
	r := NewRegistry()
	recorder := NewKVRecorder("test")
	r.MustRegister(recorder)

	r.EnableRecording(true)

	require.Equal(t, recorder.enableRecording, true)
}
