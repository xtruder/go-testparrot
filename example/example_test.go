package example

import (
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xtruder/go-testparrot"
)

type MyTestRecorder struct {
	testparrot.StructRecorder

	Val1 string
}

var recorder = &MyTestRecorder{}

func init() {
	testparrot.Register(recorder)
}

func TestExample(t *testing.T) {
	value := "value"
	recorder.RecordField(&recorder.Val1, value)

	require.Equal(t, recorder.Val1, value)
}

func TestMain(m *testing.M) {
	//testparrot.EnableRecording(true)
	code := m.Run()

	if code != 0 {
		os.Exit(code)
		return
	}

	_, filename, _, _ := runtime.Caller(0)
	genfilename := path.Join(path.Dir(filename), "gen.go")

	file, err := os.OpenFile(genfilename, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	generator := testparrot.GoGenerator{PkgName: "example"}
	generator.Generate(testparrot.R, file)
}
