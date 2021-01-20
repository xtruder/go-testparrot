package testparrot

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

// current package name and path, we need those when generating, so we can
// import testparrot package
var (
	pkgName = "testparrot"
	pkgPath = reflect.TypeOf(Generator{}).PkgPath()
)

func valToPtr(val interface{}) interface{} {
	switch v := val.(type) {
	case bool:
		return &v
	case int:
		return &v
	case int8:
		return &v
	case int16:
		return &v
	case int32:
		return &v
	case int64:
		return &v
	case uint:
		return &v
	case uint8:
		return &v
	case uint16:
		return &v
	case uint32:
		return &v
	case uint64:
		return &v
	case float32:
		return &v
	case float64:
		return &v
	case complex64:
		return &v
	case complex128:
		return &v
	default:
		return &v
	}
}

// getTestPath method walks up the stack and tries to get path of file where test is in
func getTestPath(t *testing.T) (string, error) {
	skip := 0

	// use only what is before slash in test name
	testName := strings.Split(t.Name(), "/")[0]

	for pc, path, _, ok := runtime.Caller(skip); ok; pc, path, _, ok = runtime.Caller(skip) {
		funcName := runtime.FuncForPC(pc).Name()

		lastSlash := strings.LastIndexByte(funcName, '/')
		if lastSlash < 0 {
			lastSlash = 0
		}
		firstDot := strings.IndexByte(funcName[lastSlash:], '.') + lastSlash

		funcName = funcName[(firstDot + 1):]

		// we assume name of the test name prefixes function name and that
		// file has _test.go suffix
		if strings.HasPrefix(funcName, testName) && strings.HasSuffix(path, "_test.go") {
			return path, nil
		}

		skip++
	}

	return "", fmt.Errorf("test filename not found for: %s", t.Name())
}

// getPkgInfo gets package path, name and fs location of current package
func getPkgInfo(skip int, pkgNameFromSource bool) (pkgPath string, pkgName string, fsPath string, err error) {
	pc, filename, _, ok := runtime.Caller(skip + 1)
	if !ok {
		err = errors.New("could not find package path")
		return
	}

	fsPath = path.Dir(filename)

	// example: github.com/xtruder/go-testparrot.TestValToCode.func1
	funcName := runtime.FuncForPC(pc).Name()
	lastSlash := strings.LastIndexByte(funcName, '/')
	if lastSlash < 0 {
		lastSlash = 0
	}
	firstDot := strings.IndexByte(funcName[lastSlash:], '.') + lastSlash

	// everything until first for after last slash is package path
	pkgPath = funcName[:firstDot]

	// getting package name from package path is unreliable
	pkgName = funcName[(lastSlash + 1):firstDot]

	// retrive package name by parsing source. This is only usable if
	// source code is avalible, like when testing.
	if pkgNameFromSource {
		var astFile *ast.File

		// package name cannot be reliably retrived from runtime info, so it needs
		// to be read from filesyste. This can only be used if source is provided.
		fset := token.NewFileSet()
		astFile, err = parser.ParseFile(fset, filename, nil, parser.PackageClauseOnly)
		if err != nil {
			return
		}

		if astFile.Name == nil {
			err = fmt.Errorf("package name not found")
			return
		}

		pkgName = astFile.Name.Name
	}

	return
}

func newErr(err error) error {
	return fmt.Errorf("testparrot: %v", err)
}
