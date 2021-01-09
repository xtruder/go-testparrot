package generator

import (
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

func valToPtr(val interface{}) interface{} {
	switch v := val.(type) {
	case bool:
		return &v
	case int:
		return &v
	case int8:
		return &v
	case int16:
		return &v
	case int32:
		return &v
	case int64:
		return &v
	case uint:
		return &v
	case uint8:
		return &v
	case uint16:
		return &v
	case uint32:
		return &v
	case uint64:
		return &v
	case float32:
		return &v
	case float64:
		return &v
	case complex64:
		return &v
	case complex128:
		return &v
	default:
		return &v
	}
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
			value:    valToPtr("test"),
			expected: "generator.ValToPtr(\"test\")",
		},
		{
			name:     "ptr to ptr",
			value:    valToPtr(valToPtr("test")),
			expected: "generator.ValToPtr(generator.ValToPtr(\"test\"))",
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
			expected: "map[Key]Value{\n\tKey{\n\t\t\"K1\": \"k1\",\n\t\t\"K2\": \"k2\",\n\t}: Value{\n\t\t\"V1\": \"v1\",\n\t\t\"V2\": 1,\n\t},\n\tKey{\n\t\t\"K1\": \"k1\",\n\t\t\"K2\": \"k3\",\n\t}: Value{\n\t\t\"V1\": \"v2\",\n\t\t\"V2\": 3,\n\t},\n}",
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
				V3: valToPtr(1.1).(*float64),
			},
			expected: "simpleStruct{\n\t\"V1\": \"test\",\n\t\"V2\": 10,\n\t\"V3\": generator.ValToPtr(1.1),\n}",
		},
		{
			name: "nested struct",
			value: nestedStruct{
				V1: "string",
				V2: &nestedStruct{
					V1: "value",
				},
			},
			expected: "nestedStruct{\n\t\"V1\": \"string\",\n\t\"V2\": &nestedStruct{\"V1\": \"value\"},\n}",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			g := &generator{PkgName: "github.com/xtruder/go-test-record/generator"}
			code, err := valToCode(g, reflect.ValueOf(test.value))
			require.NoError(t, err)
			require.Equal(t, test.expected, fmt.Sprintf("%#v", code))
		})
	}
}
