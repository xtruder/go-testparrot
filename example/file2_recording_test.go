// Code generated by testparrot. DO NOT EDIT.

package example

import gotestparrot "github.com/xtruder/go-testparrot"

func init() {
	gotestparrot.R.Load("TestKVExample", []gotestparrot.Recording{{
		Key: "dog1",
		Value: Dog{
			Age:   9,
			Breed: "Terrier",
			Name:  "Lido",
		},
	}, {
		Key: "dog2",
		Value: Dog{
			Age:   12,
			Breed: "Cavalier Kind Charles Spaniel",
			Name:  "Rex",
		},
	}})
}
