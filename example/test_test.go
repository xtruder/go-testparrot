package example

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xtruder/go-testparrot"
)

type Dog struct {
	Name  string
	Breed string
	Age   int
	Note  string
}

func TestSequentialExample(t *testing.T) {
	dog1 := Dog{"Fido", "Terrier", 4, ""}
	require.Equal(t, testparrot.RecordNext(t, dog1), dog1)

	dog2 := Dog{"Mika", "American Foxhound", 8,
		`this is an awesome nice dog and a good friend,
must really have it!`,
	}
	require.Equal(t, testparrot.RecordNext(t, dog2), dog2)
}

func TestKVExample(t *testing.T) {
	dog1 := Dog{"Lido", "Terrier", 9, ""}
	require.Equal(t, testparrot.Record(t, "dog1", dog1), dog1)

	dog2 := Dog{"Rex", "Cavalier Kind Charles Spaniel", 12, ""}
	require.Equal(t, testparrot.Record(t, "dog2", dog2), dog2)
}
