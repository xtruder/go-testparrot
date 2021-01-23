package testparrot

import (
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValToPtr(t *testing.T) {
	v := "test"
	require.IsType(t, valToPtr("test"), &v)
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
