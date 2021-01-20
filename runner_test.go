package testparrot

import (
	"flag"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBeforeTests(t *testing.T) {
	defineTestparrotFlags()

	t.Run("recording disabled", func(t *testing.T) {
		recorder := NewRecorder()
		recorder.Load("test", []Recording{{"key", "value"}})

		BeforeTests(recorder)

		require.Contains(t, recorder.allRecordings, "test")

		// checks wheter testparrot.flag exists
		recordFlagExists := false
		flag.VisitAll(func(f *flag.Flag) {
			recordFlagExists = recordFlagExists || f.Name == "testparrot.record"
		})

		require.True(t, recordFlagExists)
	})

	t.Run("recording enabled", func(t *testing.T) {
		recorder := NewRecorder()
		recorder.Load("test", []Recording{{"key", "value"}})

		flag.Set("testparrot.record", "true")
		defer flag.Set("testparrot.record", "false") // reset flag value

		BeforeTests(recorder)

		require.Empty(t, recorder.allRecordings)
	})
}

func TestAfterTests(t *testing.T) {
	tmpDir := t.TempDir()

	defineTestparrotFlags()

	flag.Set("testparrot.dest", tmpDir)
	flag.Set("testparrot.filename", "gen.go")
	defer flag.Set("testparrot.dest", "")
	defer flag.Set("testparrot.filename", "")
	flag.Parse()

	genPath := path.Join(tmpDir, "gen.go")

	t.Run("recording disabled", func(t *testing.T) {
		recorder := NewRecorder()
		recorder.Load("test", []Recording{{"key", "value"}})

		AfterTests(recorder, "recorder")

		// make sure no file was generated
		_, err := os.Stat(genPath)
		require.True(t, os.IsNotExist(err))
	})

	t.Run("single file", func(t *testing.T) {
		recorder := NewRecorder()
		recorder.Load("test", []Recording{{"key", "value"}})
		recorder.EnableRecording(true)

		flag.Set("testparrot.pkgpath", "my/go-pkg")
		defer flag.Set("testparrot.pkgpath", "")
		flag.Set("testparrot.pkgname", "pkg")
		defer flag.Set("testparrot.pkgname", "")

		AfterTests(recorder, "recorder")

		// make sure file was generated
		_, err := os.Stat(genPath)
		require.NoError(t, err)

		contents, err := ioutil.ReadFile(genPath)
		if err != nil {
			panic(err)
		}

		require.Contains(t, string(contents), "package pkg")
	})

	t.Run("split files", func(t *testing.T) {
		recorder := NewRecorder()
		recorder.Load("test1", []Recording{{"key1", "value1"}})
		recorder.Load("test2", []Recording{{"key1", "value1"}})
		recorder.Load("test3", []Recording{{"key2", "value2"}})
		recorder.Load("test4", []Recording{{"key2", "value2"}})
		recorder.testFilenames["test1"] = "file1_test.go"
		recorder.testFilenames["test2"] = "file1_test.go"
		recorder.testFilenames["test3"] = "file2_test.go"
		recorder.testFilenames["test4"] = "file2_test.go"
		recorder.EnableRecording(true)

		flag.Set("testparrot.pkgpath", "my/go-pkg")
		defer flag.Set("testparrot.pkgpath", "")
		flag.Set("testparrot.pkgname", "pkg")
		defer flag.Set("testparrot.pkgname", "")
		flag.Set("testparrot.splitfiles", "true")
		defer flag.Set("testparrot.splitfiles", "")

		AfterTests(recorder, "recorder")

		// make sure files were generated
		require.FileExists(t, path.Join(tmpDir, "file1_recording_test.go"))
		require.FileExists(t, path.Join(tmpDir, "file2_recording_test.go"))

		contents, err := ioutil.ReadFile(path.Join(tmpDir, "file1_recording_test.go"))
		if err != nil {
			panic(err)
		}

		require.Contains(t, string(contents), "test1")
		require.Contains(t, string(contents), "test2")
		require.NotContains(t, string(contents), "test3")
		require.NotContains(t, string(contents), "test4")

		contents, err = ioutil.ReadFile(path.Join(tmpDir, "file2_recording_test.go"))
		if err != nil {
			panic(err)
		}

		require.Contains(t, string(contents), "test3")
		require.Contains(t, string(contents), "test4")
		require.NotContains(t, string(contents), "test1")
		require.NotContains(t, string(contents), "test2")
	})
}
