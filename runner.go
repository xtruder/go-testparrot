package testparrot

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
)

var (
	enableRecordingFlag *bool
	splitFilesFlag      *bool
	destFlag            *string
	filenameFlag        *string
	pkgPathFlag         *string
	pkgNameFlag         *string
)

func defineTestparrotFlags() {
	// if flags have not been yet define, define them
	if flag.Lookup("testparrot.record") == nil {
		enableRecordingFlag = flag.Bool("testparrot.record", false, "whether to enable testparrot recording")
		splitFilesFlag = flag.Bool("testparrot.splitfiles", false, "whether to split tests into multiple files")
		destFlag = flag.String("testparrot.dest", "", "override destination path")
		filenameFlag = flag.String("testparrot.filename", "", "override destination filename")
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
	pkgPath, pkgName, pkgFsPath, err := getPkgInfo(skip+1, true)
	if err != nil {
		panic(newErr(err))
	}

	dest := pkgFsPath
	if *destFlag != "" {
		dest = *destFlag
	}

	if *pkgPathFlag != "" {
		pkgPath = *pkgPathFlag
	}

	if *pkgNameFlag != "" {
		pkgName = *pkgNameFlag
	}

	generator := NewGenerator(pkgPath, pkgName)
	if *splitFilesFlag {
		// group test names by filename
		testNamesByFilename := map[string][]string{}
		for testName, testFilename := range recorder.testFilenames {
			testNamesByFilename[testFilename] = append(testNamesByFilename[testFilename], testName)
		}

		// for every filename generate recordings
		for testFilename, fileTestNames := range testNamesByFilename {
			genFilename := strings.TrimSuffix(testFilename, filepath.Ext(testFilename))
			genFilename = strings.TrimSuffix(genFilename, "_test") + "_recording_test.go"
			genFilePath := path.Join(dest, genFilename)

			filter := func(testRecordings map[string][]Recording) map[string][]Recording {
				result := map[string][]Recording{}

				for testName, recordings := range testRecordings {
					for _, fileTestName := range fileTestNames {
						if fileTestName == testName {
							result[testName] = recordings
						}
					}
				}

				return result
			}

			opts := GenOptions{
				RecorderVar: recorderVar,
				Filter:      filter,
			}
			err = generator.GenerateToFile(recorder, opts, genFilePath)
			if err != nil {
				panic(newErr(err))
			}
		}
	} else {
		var genFilePath string
		if *filenameFlag == "" {
			genFileName := fmt.Sprintf("%s_recording_test.go", pkgName)
			genFilePath = path.Join(dest, genFileName)
		} else {
			genFilePath = path.Join(dest, *filenameFlag)
		}

		opts := GenOptions{RecorderVar: recorderVar}
		err = generator.GenerateToFile(recorder, opts, genFilePath)

		if err != nil {
			panic(newErr(err))
		}
	}
}
