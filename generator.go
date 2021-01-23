package testparrot

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	. "github.com/dave/jennifer/jen"
)

const (
	valToPtrF     = "ValToPtr"
	headerComment = "Code generated by testparrot. DO NOT EDIT."
)

type GenOptions struct {
	RecorderVar string
	Filter      func(map[string][]Recording) map[string][]Recording
}

// Generator generates golang code
type Generator struct {
	// Path of generate package
	pkgPath string

	// Name of generate package
	pkgName string
}

func NewGenerator(pkgPath, pkgName string) *Generator {
	return &Generator{pkgPath: pkgPath, pkgName: pkgName}
}

func (g *Generator) GenerateToFile(recorder *Recorder, opts GenOptions, filePath string) error {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0660)
	if err != nil {
		return err
	}

	if err := g.Generate(recorder, opts, file); err != nil {
		return err
	}

	return file.Close()
}

func (g *Generator) Generate(recorder *Recorder, opts GenOptions, out io.Writer) error {
	f := NewFilePathName(g.pkgPath, g.pkgName)

	f.HeaderComment(headerComment)

	keys := make([]string, 0, len(recorder.allRecordings))

	allRecordings := recorder.allRecordings
	if opts.Filter != nil {
		allRecordings = opts.Filter(allRecordings)
	}

	for key := range allRecordings {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	statements := []Code{}
	for _, key := range keys {
		testKey := key
		recordings := allRecordings[key]

		// call load on global recorder or on locally defined recorder
		var loadF *Statement
		if recorder == R {
			loadF = Qual(pkgPath, "R.Load")
		} else {
			loadF = Id(opts.RecorderVar + ".Load")
		}

		val, err := valToCode(g, reflect.ValueOf(recordings), reflect.Value{})
		if err != nil {
			return err
		}

		key, err := valToCode(g, reflect.ValueOf(testKey), reflect.Value{})
		if err != nil {
			return err
		}

		loadCall := loadF.Call(key, val)
		statements = append(statements, loadCall)
	}

	// Create init method
	f.Func().Id("init").Params().Block(statements...)

	// Render code
	return f.Render(out)
}

func typeToCode(g *Generator, typ reflect.Type) *Statement {
	switch typ.Kind() {
	case reflect.Interface:
		return Interface()
	default:
		pkgPath := typ.PkgPath()

		if pkgPath == "" || g.pkgPath == pkgPath {
			return Id(typ.Name())
		}

		return Qual(pkgPath, typ.Name())
	}
}

func sliceToCode(g *Generator, sliceVal reflect.Value, parent reflect.Value) (Code, error) {
	typ := sliceVal.Type()
	elemType := typ.Elem()

	a := sliceVal.Type().Name()
	fmt.Println(a)

	var litValue Code
	var values []Code
	var err error

	if sliceVal.Kind() == reflect.Slice && elemType.Kind() == reflect.Uint8 {
		s := reflect.New(reflect.SliceOf(elemType)).Elem()
		s = reflect.AppendSlice(s, sliceVal)

		val := s.Interface().([]byte)
		if utf8.Valid(val) {
			litValue, err = litToCode(g, reflect.ValueOf(string(val)))
			if err != nil {
				return nil, err
			}
		}
	}

	if litValue == nil {
		values = []Code{}
		for i := 0; i < sliceVal.Len(); i++ {
			elem := sliceVal.Index(i)

			code, err := valToCode(g, elem, sliceVal)
			if err != nil {
				return nil, err
			}

			values = append(values, code)
		}
	}

	if typ.Name() != "" {
		prefix := Id(typ.Name())

		if litValue != nil {
			return prefix.Call(litValue), nil
		}

		return prefix.Values(values...), nil
	}

	return Index().Add(typeToCode(g, elemType)).Values(values...), nil
}

func mapToCode(g *Generator, mapVal reflect.Value, parent reflect.Value) (Code, error) {
	values := Dict{}

	var keyType reflect.Type
	var valType reflect.Type
	for _, k := range mapVal.MapKeys() {
		v := mapVal.MapIndex(k)

		keyCode, err := valToCode(g, k, mapVal)
		if err != nil {
			return nil, err
		}

		valCode, err := valToCode(g, v, mapVal)
		if err != nil {
			return nil, err
		}

		values[keyCode] = valCode

		keyType = k.Type()
		valType = v.Type()
	}

	typeCode := typeToCode(g, valType)

	return Map(typeToCode(g, keyType)).Add(typeCode).Values(values), nil
}

func structToCode(g *Generator, structVal reflect.Value, parent reflect.Value) (Code, error) {
	structType := structVal.Type()

	if structType.Kind() == reflect.Ptr {
		structType = structType.Elem()
		structVal = structVal.Elem()
	}

	values := Dict{}
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		val := structVal.Field(i)

		// only set value if it is not zero
		if val.IsZero() {
			continue
		}

		if val.Kind() == reflect.Ptr && val.IsNil() {
			continue
		}

		code, err := valToCode(g, val, structVal)
		if err != nil {
			return nil, err
		}

		values[Id(field.Name)] = code
	}

	// if parent is a slice and struct type is same as slice type, we can omit struct type
	if parent.IsValid() && parent.Kind() == reflect.Slice {
		elemType := parent.Type().Elem()
		if elemType.Name() == structType.Name() {
			return Values(values), nil
		}
	}

	return typeToCode(g, structType).Values(values), nil
}

func ptrToCode(g *Generator, ptrVal reflect.Value, parent reflect.Value) (Code, error) {
	val := ptrVal.Elem()
	valType := val.Type()

	valToPtrF := Qual(pkgPath, valToPtrF)

	switch valType.Kind() {
	// for most type need to get pointer to literal value
	case
		reflect.Bool,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uintptr,
		reflect.Float32,
		reflect.Float64,
		reflect.Complex64,
		reflect.Complex128,
		reflect.String:
		return valToPtrF.Call(Lit(val.Interface())).Assert(Id("*" + valType.Name())), nil
	case reflect.Interface:
		return ptrToCode(g, val, ptrVal)
	case reflect.Ptr:
		code, err := ptrToCode(g, val, ptrVal)
		if err != nil {
			return nil, err
		}

		return valToPtrF.Call(code), nil
	default:
		code, err := valToCode(g, val, ptrVal)
		if err != nil {
			return nil, err
		}

		return Op("&").Add(code), nil
	}
}

func litToCode(g *Generator, value reflect.Value) (Code, error) {

	switch value.Type().Kind() {
	case reflect.String:
		val := value.Interface().(string)

		hasNewlines := strings.Count(val, "\n") > 0
		hasOnlyNewlineAtEnd := strings.Count(val, "\n") == 1 && strings.HasSuffix(val, "\n")
		if hasNewlines && !hasOnlyNewlineAtEnd {
			return Id(fmt.Sprintf("`%s`", val)), nil
		}

		fallthrough
	default:
		return Lit(value.Interface()), nil
	}
}

func valToCode(g *Generator, value reflect.Value, parent reflect.Value) (Code, error) {
	if value == (reflect.Value{}) {
		return Nil(), nil
	}

	// check if value is of special concrete types that cannot be easily generated
	switch v := value.Interface().(type) {
	case time.Time:
		return Qual("time", "Date").Call(
			Lit(v.Year()),
			Lit(int(v.Month())),
			Lit(v.Day()),
			Lit(v.Hour()),
			Lit(v.Minute()),
			Lit(v.Second()),
			Lit(v.Nanosecond()),
			Qual("time", "FixedZone").Call(Lit(v.Location().String()), Lit(0)),
		), nil
	}

	valType := value.Type()

	switch valType.Kind() {
	// for most values construct literal values
	case
		reflect.Bool,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uintptr,
		reflect.Float32,
		reflect.Float64,
		reflect.Complex64,
		reflect.Complex128,
		reflect.String:
		return litToCode(g, value)
	case reflect.Array, reflect.Slice:
		return sliceToCode(g, value, parent)
	case reflect.Interface:
		return valToCode(g, value.Elem(), value)
	case reflect.Map:
		return mapToCode(g, value, parent)
	case reflect.Ptr:
		return ptrToCode(g, value, parent)
	case reflect.Struct:
		return structToCode(g, value, parent)
	default:
		panic("unsupported kind")
	}
}
