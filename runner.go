package testparrot

import (
	"flag"
	"fmt"
	"os"
	"path"
	"testing"
)

var (
	enableRecordingFlag *bool
	destFlag            *string
	pkgPathFlag         *string
	pkgNameFlag         *string
)

func init() {
	// a bit of hackish was to check if tests are enabled, but seems to work reliably
	if isInTests() {
		defineTestparrotFlags()
	}
}

func defineTestparrotFlags() {
	// if flags have not been yet define, define them
	if flag.Lookup("testparrot.record") == nil {
		enableRecordingFlag = flag.Bool("testparrot.record", false, "whether to enable testparrot recording")
		destFlag = flag.String("testparrot.dest", "", "override destination path")
		pkgPathFlag = flag.String("testparrot.pkgpath", "", "override package path")
		pkgNameFlag = flag.String("testparrot.pkgname", "", "override package name")
	}
}

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
func BeforeTests(recorder *Recorder) {
	beforeTests(recorder)
}

// AfterTests is method to use in TestMain after running tests
func AfterTests(recorder *Recorder, recorderVar string) {
	afterTests(recorder, recorderVar, 1)
}

func beforeTests(recorder *Recorder) {

	// define additional test flag for enabling testparrot recording and parse it
	defineTestparrotFlags()
	flag.Parse()

	if !recorder.RecordingEnabled() {
		recorder.EnableRecording(*enableRecordingFlag)
	}

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

	var genFilePath string
	if *destFlag == "" {
		genFileName := fmt.Sprintf("%s_recording_test.go", pkgName)
		genFilePath = path.Join(fsPath, genFileName)
	} else {
		genFilePath = *destFlag
	}

	if *pkgPathFlag != "" {
		pkgPath = *pkgPathFlag
	}

	if *pkgNameFlag != "" {
		pkgName = *pkgNameFlag
	}

	generator := NewGenerator(pkgPath, pkgName)
	err = generator.GenerateToFile(recorder, recorderVar, genFilePath)

	if err != nil {
		panic(newErr(err))
	}
}
