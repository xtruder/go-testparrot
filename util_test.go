package testparrot

import (
	"path"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestDecode(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		typ      interface{}
		expected interface{}
	}{
		{
			name:     "time",
			value:    "1999-01-02T03:04:05Z",
			typ:      time.Time{},
			expected: must(time.Parse(time.RFC3339, "1999-01-02T03:04:05Z")),
		}, {

			name:     "timeptr",
			value:    "1999-01-02T03:04:05Z",
			typ:      &time.Time{},
			expected: valToPtr(must(time.Parse(time.RFC3339, "1999-01-02T03:04:05Z"))),
		},
		{
			name:     "uuid",
			value:    "6ba7b814-9dad-11d1-80b4-00c04fd430c8",
			typ:      uuid.UUID{},
			expected: uuid.MustParse("6ba7b814-9dad-11d1-80b4-00c04fd430c8"),
		},
		{
			name:     "uuid ptr",
			value:    "6ba7b814-9dad-11d1-80b4-00c04fd430c8",
			typ:      &uuid.UUID{},
			expected: valToPtr(uuid.MustParse("6ba7b814-9dad-11d1-80b4-00c04fd430c8")),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			v := Decode(test.value, test.typ)
			require.IsType(t, test.typ, v)
			require.Equal(t, test.expected, v)
		})
	}
}

func TestValToPtr(t *testing.T) {
	v := "test"
	require.IsType(t, valToPtr("test"), &v)
}

func TestMust(t *testing.T) {
	a := func(val string) (string, error) {
		return val, nil
	}

	require.Equal(t, "test", must(a("test")))
}

func TestGetTestPath(t *testing.T) {
	t.Run("subtest", func(t *testing.T) {
		t.Run("subsubtest", func(t *testing.T) {
			var filename string
			var err error

			func() {
				filename, err = getTestPath(t)
			}()

			require.NoError(t, err)
			require.Equal(t, "util_test.go", path.Base(filename))
		})
	})
}
