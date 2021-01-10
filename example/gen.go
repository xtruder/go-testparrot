package example

import gotestparrot "github.com/xtruder/go-testparrot"

func init() {
	gotestparrot.Load("MyTestRecorder", []gotestparrot.Record{gotestparrot.Record{
		Name:  "Val1",
		Value: "value",
	}})
}
