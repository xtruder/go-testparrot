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
