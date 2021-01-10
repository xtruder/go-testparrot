package testparrot

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type testField struct {
	Val1 string
	Val2 string
}

type testRecorder struct {
	StructRecorder

	R1 string
	R2 testField
}

func TestRecordRecording(t *testing.T) {
	r := &testRecorder{}
	r.Init(r)
	r.EnableRecording(true)

	expectedStr := r.RecordField(&r.R1, "value")
	require.Equal(t, "value", expectedStr)
	require.Equal(t, "value", r.R1)

	val2 := testField{"val1", "val2"}

	expectedStuct := r.RecordField(&r.R2, val2)
	require.Equal(t, val2, expectedStuct)

	records := r.Records()

	require.Contains(t, records, Record{"R1", "value"})
	require.Contains(t, records, Record{"R2", val2})
}

func TestRecordLoad(t *testing.T) {
	r := &testRecorder{}
	r.Init(r)
	r.Load([]Record{{
		"R1", "value",
	}, {
		"R2", testField{"val1", "val2"},
	}})

	expected := r.RecordField(&r.R1, "value")
	require.Equal(t, "value", expected)
	require.Equal(t, "value", r.R1)
}
