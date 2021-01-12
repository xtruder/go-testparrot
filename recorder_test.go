package testparrot

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewRecorder(t *testing.T) {
	recorder := NewRecorder()
	require.IsType(t, &Recorder{}, recorder)
}

func TestLoad(t *testing.T) {
	name := t.Name()

	t.Run("should load recording", func(t *testing.T) {
		recorder := NewRecorder()
		recorder.Load(name, []Recording{{"key", "value"}})
		require.Contains(t, recorder.allRecordings, name)
	})

	t.Run("should panic on duplicate recordings for single test", func(t *testing.T) {
		recorder := NewRecorder()
		require.PanicsWithError(t,
			"testparrot: recordings already loaded for test 'TestLoad'",
			func() {
				recorder.Load(name, []Recording{{"key", "value"}})
				recorder.Load(name, []Recording{{"key", "value"}})
			},
		)
	})
}

func TestRecorderReset(t *testing.T) {
	recorder := NewRecorder()
	recorder.allRecordings["test"] = []Recording{{"key", "value"}}
	recorder.counters["test"] = 0

	recorder.Reset()
	require.Empty(t, recorder.allRecordings)
	require.Empty(t, recorder.counters)
}

func TestRecorderRecord(t *testing.T) {
	t.Run("recording enabled", func(t *testing.T) {
		recorder := NewRecorder()
		recorder.EnableRecording(true)
		require.Equal(t, "value1", recorder.Record(t, "key1", "value1"))
		require.Equal(t, "value2", recorder.Record(t, "key2", "value2"))
		require.Contains(t, recorder.allRecordings, t.Name())
		require.Contains(t, recorder.allRecordings[t.Name()], Recording{"key1", "value1"})
		require.Contains(t, recorder.allRecordings[t.Name()], Recording{"key2", "value2"})
	})

	t.Run("panics on duplicate key when recording enabled", func(t *testing.T) {
		recorder := NewRecorder()
		recorder.EnableRecording(true)
		require.PanicsWithError(t,
			fmt.Sprintf("recording with key '%s' already exists for test '%s'", "key1", t.Name()),
			func() {
				recorder.Record(t, "key1", "value1")
				recorder.Record(t, "key1", "value1")
			},
		)
	})

	t.Run("panics on missing key if not recording", func(t *testing.T) {
		recorder := NewRecorder()
		require.PanicsWithError(t,
			fmt.Sprintf("testparrot: recording with key '%s' not found for test '%s'", "key1", t.Name()),
			func() {
				recorder.Record(t, "key1", "value1")
			},
		)
	})

	t.Run("reply values", func(t *testing.T) {
		value := "value"
		recorder := NewRecorder()
		recorder.Load(t.Name(), []Recording{{"key", value}})
		require.Equal(t, recorder.Record(t, "key", "value1"), value)
	})
}

func TestRecorderRecordNext(t *testing.T) {
	t.Run("recording enabled", func(t *testing.T) {
		recorder := NewRecorder()
		recorder.EnableRecording(true)
		require.Equal(t, "value1", recorder.RecordNext(t, "value1"))
		require.Equal(t, "value2", recorder.RecordNext(t, "value2"))
		require.Contains(t, recorder.allRecordings, t.Name())
		require.Contains(t, recorder.allRecordings[t.Name()], Recording{0, "value1"})
		require.Contains(t, recorder.allRecordings[t.Name()], Recording{1, "value2"})
	})

	t.Run("panics on missing key if not recording", func(t *testing.T) {
		recorder := NewRecorder()
		require.PanicsWithError(t,
			fmt.Sprintf("testparrot: recording with key '%s' not found for test '%s'", "0", t.Name()),
			func() {
				recorder.RecordNext(t, "key1")
			},
		)
	})

	t.Run("reply values", func(t *testing.T) {
		value := "value"
		recorder := NewRecorder()
		recorder.Load(t.Name(), []Recording{{0, value}})
		require.Equal(t, recorder.RecordNext(t, "value"), value)
	})
}

func TestEnableRecording(t *testing.T) {
	recorder := NewRecorder()
	recorder.EnableRecording(true)
	require.True(t, recorder.RecordingEnabled())
}
