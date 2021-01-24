package testparrot

import (
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDecode(t *testing.T) {
	v := Decode([]byte("1999-01-02T03:04:05Z"), &time.Time{})
	require.IsType(t, &time.Time{}, v)

	v = Decode([]byte("1999-01-02T03:04:05Z"), time.Time{})
	require.IsType(t, time.Time{}, v)
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
