package testparrot

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type testStruct struct {
	V1 string
	V2 int
	V3 *testStruct
	V4 []string
}

func TestGenerate(t *testing.T) {
	registry := NewRegistry()
	registry.MustRegister(NewKVRecorder("recorder1"))
	registry.MustRegister(NewKVRecorder("recorder2"))

	registry.Load("recorder1", []Record{
		{"key1", "value1"},
		{"key2", 1},
	})

	registry.Load("recorder2", []Record{
		{"key1", "value1"},
		{"key2", 1},
	})

	buf := &bytes.Buffer{}
	generator := GoGenerator{}
	err := generator.Generate(registry, buf)

	require.NoError(t, err)
	require.Equal(t, "", buf.String())
}

func TestValToCode(t *testing.T) {
	type simpleStruct struct {
		V1 string
		V2 int
		V3 *float64
		V4 bool
	}

	type nestedStruct struct {
		V1 string
		V2 *nestedStruct
	}

	type Key struct {
		K1 string
		K2 string
	}

	type Value struct {
		V1 string
		V2 int
	}

	tests := []struct {
		name     string
		value    interface{}
		expected string
	}{
		{
			name:     "literal str",
			value:    "test",
			expected: "\"test\"",
		},
		{
			name:     "ptr literal",
			value:    ValToPtr("test"),
			expected: "gotestparrot.ValToPtr(\"test\")",
		},
		{
			name:     "ptr to ptr",
			value:    ValToPtr(ValToPtr("test")),
			expected: "gotestparrot.ValToPtr(gotestparrot.ValToPtr(\"test\"))",
		},
		{
			name: "simple map",
			value: map[string]string{
				"key": "value",
			},
			expected: "map[string]string{\"key\": \"value\"}",
		},
		{
			name: "double kv map",
			value: map[Key]Value{
				{"k1", "k2"}: {"v1", 1},
				{"k1", "k3"}: {"v2", 3},
			},
			expected: "map[Key]Value{\n\tKey{\n\t\tK1: \"k1\",\n\t\tK2: \"k2\",\n\t}: Value{\n\t\tV1: \"v1\",\n\t\tV2: 1,\n\t},\n\tKey{\n\t\tK1: \"k1\",\n\t\tK2: \"k3\",\n\t}: Value{\n\t\tV1: \"v2\",\n\t\tV2: 3,\n\t},\n}",
		},
		{
			name:     "simple slice",
			value:    []string{"val1", "val2"},
			expected: "[]string{\"val1\", \"val2\"}",
		},
		{
			name:     "struct slice",
			value:    []Value{{"v1", 1}, {"v2", 2}},
			expected: "[]Value{Value{\n\tV1: \"v1\",\n\tV2: 1,\n}, Value{\n\tV1: \"v2\",\n\tV2: 2,\n}}",
		},
		{
			name: "interface map",
			value: map[interface{}]interface{}{
				"key": "value",
				10:    "value",
			},
			expected: "map[interface{}]interface{}{\n\t\"key\": \"value\",\n\t10:    \"value\",\n}",
		},
		{
			name: "simple struct",
			value: simpleStruct{
				V1: "test",
				V2: 10,
				V3: ValToPtr(1.1).(*float64),
			},
			expected: "simpleStruct{\n\tV1: \"test\",\n\tV2: 10,\n\tV3: gotestparrot.ValToPtr(1.1),\n}",
		},
		{
			name: "nested struct",
			value: nestedStruct{
				V1: "string",
				V2: &nestedStruct{
					V1: "value",
				},
			},
			expected: "nestedStruct{\n\tV1: \"string\",\n\tV2: &nestedStruct{V1: \"value\"},\n}",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			g := &GoGenerator{PkgName: "github.com/xtruder/go-testparrot"}
			code, err := valToCode(g, reflect.ValueOf(test.value))
			require.NoError(t, err)
			require.Equal(t, test.expected, fmt.Sprintf("%#v", code))
		})
	}
}
