![tests](https://github.com/xtruder/go-testparrot/workflows/test/badge.svg)

# go-testparrot :parrot:

go-testparrot records and replies expected test values, so you don't have
to hardcode complex test values.

## About

Are you **tired of hard coding** values in tests?
Do you **copy paste failed test values** like as monkey?

**What if there would be a record button to record test values and reply them
later, just like a parrot!**

### Features

- Simple interface for record and reply values based on sequential or key value.
- Generation of recorded values in readable go code.

## Quick start

### Create a package with some tests:

```go
package example

import (
	"testing"

	"github.com/xtruder/go-testparrot"
)

func doSomething() string {
    return "value"
}

func TestSomething(t* testing.T) {
    value := doSomething()

    // When running in recording mode return value will equal
    // passed value and file with recordings will be generated.
    // When running without recording, values from recording file
    // will be replied.
    expected := testparrot.RecordNext(t, value)

    if value != expected {
        t.Errorf("doSomething() = %v; want %v", value, expected)
    }
}

func TestMain(m *testing.M) {
	testparrot.Run(m)
}
```

You must provide `TestMain` method that will run `testparrot.Run` of if you need additional steps after/before running tests, you can also use `testparrot.BeforeTests` and `testparrot.AfterTests` helper methods.

### Record values

To record values run tests with recording enabled:

```bash
go test <package> -testparrot.record
```

This will record values and save them into `<package>_recording_test.go` file in same directory as tests.

You can also use `go:generate` by placing comment like:

```go
//go:generate go test ./ -testparrot.record
```

in package under test and running `go generate <package>`

### Run tests

Run tests like you woul ussually run them, but with recording disabled:

```bash
go test <package>
```

## Developing go-testparrot

See
[CONTRIBUTING.md](https://github.com/xtruder/go-testparrot/blob/master/.github/CONTRIBUTING.md)
for best practices and instructions on setting up your development environment
to work on Packer.