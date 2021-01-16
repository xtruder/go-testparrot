![tests](https://github.com/xtruder/go-testparrot/workflows/test/badge.svg)

# go-testparrot :parrot:

go-testparrot records and replies expected test values, so you don't have
to hardcode complex test values.

## About

Are you **tired of hard coding** values in tests?
Do you **copy paste failed test values** like as monkey?

**What if there would be a record button to record test values and reply them
later, just like a parrot!**

## Example

```go
//go:generate go test ./. -testparrot.record
package example

import (
	"testing"

	"github.com/xtruder/go-testparrot"
)

func TestSomething(t* testing.T) {
    value := doSomething()

    // When running in recording mode return value will equal
    // passed value and file with recordings will be generated.
    // When running without recording, values from recording file
    // will be replied.
    expected := testparrot.RecordNext(t, value)

    if value != expected {
        t.Errorf("doSomething() = %w; want %w", value, expected)
    }
}

func TestKV(t *testing.T) {
    value1 := doSomething()
    expected1 := testparrot.Record(t, "key1", value1)

    if value1 != expected1 {
        t.Errorf("doSomething() = %w; want %w", value1, expected1)
    }

    value2 := doSomethingElse()
    expected2 := testparrot.Record(t, "key2", value2)

    if value2 != expected2 {
        t.Errorf("doSomethingElse() = %w; want %w", value2, expected2)
    }
}

func TestMain(m *testing.M) {
	testparrot.Run(m)
}
```

To record values you just need to run `go generate <package>` and recorded
values will be saved to `<package>_recording_test.go` file in same package as tests.

## Development

This project requires at least `go` version `1.15` installed.

### Testing

```bash
go test ./
```

### VSCode

If you are using [visual studio code](https://code.visualstudio.com/) you can open
project in [vscode remote container](https://code.visualstudio.com/docs/remote/containers).
