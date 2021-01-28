package testparrot

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type testStruct struct {
	V1 string
	V2 int
	V3 *testStruct
	V4 []string
}

func TestNewGenerator(t *testing.T) {
	generator := NewGenerator(pkgPath, pkgName)
	require.IsType(t, &Generator{}, generator)
	require.Equal(t, pkgName, generator.pkgName)
	require.Equal(t, pkgPath, generator.pkgPath)
}

func TestGenerate(t *testing.T) {
	recorder := NewRecorder()

	recorder.Load("recorder1", []Recording{
		{"key1", "value1"},
		{"key2", 1},
	})

	recorder.Load("recorder2", []Recording{
		{"key1", "value1"},
		{"key2", 1},
	})

	buf := &bytes.Buffer{}
	generator := NewGenerator(pkgPath, pkgName)

	expected :=
		"// Code generated by testparrot. DO NOT EDIT.\n\npackage testparrot\n\nfunc init() {\n\trecorder.Load(\"recorder1\", []Recording{" +
			"{\n\t\tKey:   \"key1\",\n\t\tValue: \"value1\",\n\t}, {\n\t\tKey:   \"key2\",\n\t\tValue: 1,\n\t}})\n\trecorder.Load(\"recorder2\", []Recording{" +
			"{\n\t\tKey:   \"key1\",\n\t\tValue: \"value1\",\n\t}, {\n\t\tKey:   \"key2\",\n\t\tValue: 1,\n\t}})\n}\n"

	t.Run("to string", func(t *testing.T) {
		err := generator.Generate(recorder, GenOptions{RecorderVar: "recorder"}, buf)

		require.NoError(t, err)
		require.Equal(t, expected, buf.String())
	})

	t.Run("to file", func(t *testing.T) {
		genPath := path.Join(t.TempDir(), "gen.go")
		err := generator.GenerateToFile(recorder, GenOptions{RecorderVar: "recorder"}, genPath)
		require.NoError(t, err)

		_, err = os.Stat(genPath)
		require.NoError(t, err)

		contents, err := ioutil.ReadFile(genPath)
		if err != nil {
			panic(err)
		}

		require.Equal(t, expected, string(contents))
	})
}

func TestValToCode(t *testing.T) {
	type wrappedBytes []byte

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

	type Enum string

	tests := []struct {
		name     string
		value    interface{}
		expected string
	}{
		{
			name:     "nil",
			value:    nil,
			expected: "nil",
		},
		{
			name:     "literal str",
			value:    "test",
			expected: "\"test\"",
		},
		{
			name:     "multiline str",
			value:    "long\nmulti\nline\nstrings\nare\ngenerated\nas\nmuliline\nstrings\nand\nthis\nmakes\neverything\nmore\nreadable",
			expected: "`long\nmulti\nline\nstrings\nare\ngenerated\nas\nmuliline\nstrings\nand\nthis\nmakes\neverything\nmore\nreadable`",
		},
		{
			name:     "ptr literal",
			value:    Ptr("test"),
			expected: "gotestparrot.Ptr(\"test\").(*string)",
		},
		{
			name:     "ptr to ptr",
			value:    Ptr(Ptr("test")),
			expected: "gotestparrot.Ptr(gotestparrot.Ptr(\"test\").(*string))",
		},
		{
			name:     "slice uint8",
			value:    []uint8{0x62, 0x79, 0x74, 0x65, 0x61},
			expected: "[]uint8{uint8(0x62), uint8(0x79), uint8(0x74), uint8(0x65), uint8(0x61)}",
		},
		{
			name:     "wrapped slice ptr",
			value:    Ptr(wrappedBytes("test")).(*wrappedBytes),
			expected: "gotestparrot.Ptr(wrappedBytes(\"test\")).(*gotestparrot.wrappedBytes)",
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
			expected: "[]Value{{\n\tV1: \"v1\",\n\tV2: 1,\n}, {\n\tV1: \"v2\",\n\tV2: 2,\n}}",
		},
		{
			name:     "struct interface",
			value:    []interface{}{1, "test", Value{"v1", 1}},
			expected: "[]interface{}{1, \"test\", Value{\n\tV1: \"v1\",\n\tV2: 1,\n}}",
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
				V3: Ptr(1.1).(*float64),
			},
			expected: "simpleStruct{\n\tV1: \"test\",\n\tV2: 10,\n\tV3: gotestparrot.Ptr(1.1).(*float64),\n}",
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
		{
			name:     "time",
			value:    time.Date(1999, 1, 2, 3, 4, 5, 0, time.FixedZone("UTC-8", -8*60*60)),
			expected: "gotestparrot.Decode(\"1999-01-02T03:04:05-08:00\", time.Time{}).(time.Time)",
		},
		{
			name:     "timeptr",
			value:    valToPtr(time.Date(1999, 1, 2, 3, 4, 5, 0, time.FixedZone("UTC-8", -8*60*60))).(*time.Time),
			expected: "gotestparrot.Decode(\"1999-01-02T03:04:05-08:00\", &time.Time{}).(*time.Time)",
		},
		{
			name:     "wrapped slice bytes",
			value:    wrappedBytes{0xff, 0xfe, 0xfd},
			expected: `wrappedBytes{uint8(0xff), uint8(0xfe), uint8(0xfd)}`,
		},
		{
			name:     "wrapped slice string",
			value:    wrappedBytes(`test`),
			expected: `wrappedBytes("test")`,
		},
		{
			name:     "uuid",
			value:    uuid.MustParse("6ba7b814-9dad-11d1-80b4-00c04fd430c8"),
			expected: "gotestparrot.Decode(\"6ba7b814-9dad-11d1-80b4-00c04fd430c8\", uuid.UUID{}).(uuid.UUID)",
		},
		{
			name:     "uuid ptr",
			value:    valToPtr(uuid.MustParse("6ba7b814-9dad-11d1-80b4-00c04fd430c8")),
			expected: "gotestparrot.Decode(\"6ba7b814-9dad-11d1-80b4-00c04fd430c8\", &uuid.UUID{}).(*uuid.UUID)",
		},
		{
			name: "time interface slice",
			value: []interface{}{
				must(time.Parse(time.RFC3339, "1999-01-02T03:04:05Z")),
			},
			expected: "[]interface{}{gotestparrot.Decode(\"1999-01-02T03:04:05Z\", time.Time{}).(time.Time)}",
		},
		{
			name:     "enum",
			value:    Enum("test"),
			expected: `Enum("test")`,
		},
		{
			name:     "enum ptr",
			value:    valToPtr(Enum("test")).(*Enum),
			expected: `gotestparrot.Ptr(Enum("test")).(*Enum)`,
		},
		{
			name: "anonymous struct",
			value: struct {
				Field1 string `json:"key1"`
				Field2 struct {
					Field3 string
					Field4 int
				} `json:"key2"`
				Field3 []struct {
					Field1 string
					Field2 int
				}
				Field4 *struct {
					Field1 int
				}
				Field5 Value
				Field6 **[]Value
			}{},
			expected: "struct {\n\tField1 string `json:\"key1\"`\n\tField2 struct {\n\t\tField3 string\n\t\tField4 int\n\t} " +
				"`json:\"key2\"`\n\tField3 []struct {\n\t\tField1 string\n\t\tField2 int\n\t}\n\tField4 *struct {\n\t\tField1 " +
				"int\n\t}\n\tField5 Value\n\tField6 **[]Value\n}{}",
		},
		{
			name:     "anonymous slice struct",
			value:    []struct{ Field1 string }{},
			expected: "[]struct {\n\tField1 string\n}{}",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			g := NewGenerator(pkgPath, pkgName)
			code, err := valToCode(g, reflect.ValueOf(test.value), reflect.Value{})
			require.NoError(t, err)
			require.Equal(t, test.expected, fmt.Sprintf("%#v", code))
		})
	}
}
