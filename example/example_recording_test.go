// Code generated by testparrot. DO NOT EDIT.

package example

import gotestparrot "github.com/xtruder/go-testparrot"

func init() {
	gotestparrot.R.Load("TestKVExample", []gotestparrot.Recording{{
		Key: "dog1",
		Value: Dog{
			Breed: "Terrier",
			Name:  "Lido",
		},
	}, {
		Key: "dog2",
		Value: Dog{
			Breed: "Cavalier Kind Charles Spaniel",
			Name:  "Rex",
		},
	}})
	gotestparrot.R.Load("TestSequentialExample", []gotestparrot.Recording{{
		Key: 0,
		Value: Dog{
			Breed: "Terrier",
			Name:  "Fido",
		},
	}, {
		Key: 1,
		Value: Dog{
			Breed: "American Foxhound",
			Name:  "Mika",
		},
	}})
}
