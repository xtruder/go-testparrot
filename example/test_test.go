package example

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xtruder/go-testparrot"
)

type Dog struct {
	Name  string
	Breed string
}

func TestSequentialExample(t *testing.T) {
	dog1 := Dog{"Fido", "Terrier"}
	require.Equal(t, testparrot.RecordNext(t, dog1), dog1)

	dog2 := Dog{"Mika", "American Foxhound"}
	require.Equal(t, testparrot.RecordNext(t, dog2), dog2)
}

func TestKVExample(t *testing.T) {
	dog1 := Dog{"Lido", "Terrier"}
	require.Equal(t, testparrot.Record(t, "dog1", dog1), dog1)

	dog2 := Dog{"Rex", "Cavalier Kind Charles Spaniel"}
	require.Equal(t, testparrot.Record(t, "dog2", dog2), dog2)
}
