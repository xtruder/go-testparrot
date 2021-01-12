//go:generate go test ./. -testparrot.record
package example

import (
	"testing"

	"github.com/xtruder/go-testparrot"
)

func TestMain(m *testing.M) {
	testparrot.Run(m)
}
