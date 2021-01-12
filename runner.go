package testparrot

import (
	"flag"
	"fmt"
	"os"
	"path"
	"testing"
)

// Helper method to use in TestMain for running tests
func Run(m *testing.M) {
	beforeTests(R)

	code := m.Run()

	if code != 0 {
		os.Exit(code)
		return
	}

	afterTests(R, "", 1)
}

// BeforeTests is method to use in TestMain before running tests
func BeforeTests() {
	beforeTests(R)
}

// AfterTests is method to use in TestMain after running tests
func AfterTests(recorder *Recorder) {
	afterTests(R, "", 1)
}

func beforeTests(recorder *Recorder) {

	// define additional test flag for enabling testparrot recording and parse it
	enableRecording := flag.Bool("testparrot.record", false, "whether to enable testparrot recording")
	flag.Parse()

	recorder.EnableRecording(*enableRecording)

	// reset loaded values if recording is enabled
	if recorder.RecordingEnabled() {
		recorder.Reset()
	}
}

func afterTests(recorder *Recorder, recorderVar string, skip int) {
	// nothing to do if recording is not enabled
	if !recorder.RecordingEnabled() {
		return
	}

	// get package path and name, so we know where to put and name generated file
	pkgPath, pkgName, fsPath, err := getPkgInfo(skip+1, true)
	if err != nil {
		panic(newErr(err))
	}

	genFileName := fmt.Sprintf("%s_recording_test.go", pkgName)
	genFilePath := path.Join(fsPath, genFileName)

	generator := NewGenerator(pkgPath, pkgName)
	err = generator.GenerateToFile(recorder, recorderVar, genFilePath)

	if err != nil {
		panic(newErr(err))
	}
}
