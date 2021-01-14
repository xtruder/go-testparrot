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
	genPath := path.Join(t.TempDir(), "gen.go")
	flag.Set("testparrot.dest", genPath)
	defer flag.Set("testparrot.dest", "")
	flag.Parse()

	t.Run("recording disabled", func(t *testing.T) {
		recorder := NewRecorder()
		recorder.Load("test", []Recording{{"key", "value"}})

		AfterTests(recorder, "recorder")

		// make sure no file was generated
		_, err := os.Stat(genPath)
		require.True(t, os.IsNotExist(err))
	})

	t.Run("recording enabled", func(t *testing.T) {
		recorder := NewRecorder()
		recorder.Load("test", []Recording{{"key", "value"}})
		recorder.EnableRecording(true)

		flag.Set("testparrot.pkgpath", "my/go-pkg")
		flag.Set("testparrot.pkgname", "pkg")

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
}
