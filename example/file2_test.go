package example

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xtruder/go-testparrot"
)

func TestKVExample(t *testing.T) {
	dog1 := Dog{"Lido", "Terrier", 9, ""}
	require.Equal(t, testparrot.Record(t, "dog1", dog1), dog1)

	dog2 := Dog{"Rex", "Cavalier Kind Charles Spaniel", 12, ""}
	require.Equal(t, testparrot.Record(t, "dog2", dog2), dog2)
}
