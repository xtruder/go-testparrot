package testparrot

import (
	"encoding"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	. "github.com/dave/jennifer/jen"
)

const (
	headerComment = "Code generated by testparrot. DO NOT EDIT."
	ptrF          = "Ptr"
	decodeF       = "Decode"
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

	if sliceVal.Kind() == reflect.Slice && elemType.Kind() == reflect.Uint8 && typ.Name() != "" {
		s := reflect.New(reflect.SliceOf(elemType)).Elem()
		s = reflect.AppendSlice(s, sliceVal)

		val := s.Interface().([]byte)
		if utf8.Valid(val) {
			litValue, err = strToCode(g, string(val))
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
		prefix := typeToCode(g, typ)

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

func decodeValueToCode(g *Generator, value Code, structType reflect.Type) Code {
	var valueTypeCode *Statement
	var assertCode *Statement

	if structType.Kind() == reflect.Ptr {
		structType = structType.Elem()
		valueTypeCode = Op("&").Add(typeToCode(g, structType))
		assertCode = Op("*").Add(typeToCode(g, structType))
	} else {
		valueTypeCode = typeToCode(g, structType)
		assertCode = typeToCode(g, structType)
	}

	return Qual(pkgPath, decodeF).
		Call(
			Index().Byte().Parens(value), valueTypeCode.Values()).
		Assert(assertCode)

}

func structToCode(g *Generator, structVal reflect.Value, parent reflect.Value) (Code, error) {
	structType := structVal.Type()

	if structType.Kind() == reflect.Ptr {
		structType = structType.Elem()
		structVal = structVal.Elem()
	}

	values := Dict{}
	for i := 0; i < structType.NumField(); i++ {
		fieldType := structType.Field(i)
		fieldVal := structVal.Field(i)

		// if field is private, we cannot set it
		if unicode.IsLower(rune(fieldType.Name[0])) {
			continue
		}

		// only set value if it is not zero
		if fieldVal.IsZero() {
			continue
		}

		if fieldVal.Kind() == reflect.Ptr && fieldVal.IsNil() {
			continue
		}

		code, err := valToCode(g, fieldVal, structVal)
		if err != nil {
			return nil, err
		}

		values[Id(fieldType.Name)] = code
	}

	// if struct is empty, we can try if we can apply marshalers to get actual value
	if len(values) == 0 {
		if code, err := marshalersToCode(g, structVal, parent); code != nil {
			return code, err
		}
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

	ptrF := Qual(pkgPath, ptrF)

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
		return ptrF.Call(Lit(val.Interface())).Assert(Id("*" + valType.Name())), nil
	case reflect.Struct:
		code, err := valToCode(g, val, ptrVal)
		if err != nil {
			return nil, err
		}

		return Op("&").Add(code), nil
	case reflect.Interface:
		return ptrToCode(g, val, ptrVal)
	case reflect.Ptr:
		code, err := ptrToCode(g, val, ptrVal)
		if err != nil {
			return nil, err
		}

		return ptrF.Call(code), nil
	case reflect.Slice:
		code, err := sliceToCode(g, val, ptrVal)
		if err != nil {
			return nil, err
		}

		var elemType Code
		if valType.Name() == "" {
			elemType = Id("[]" + valType.Elem().Name())
		} else {
			elemType = Qual(valType.PkgPath(), valType.Name())
		}

		return ptrF.Call(code).Assert(Id("*").Add(elemType)), nil
	default:
		code, err := valToCode(g, val, ptrVal)
		if err != nil {
			return nil, err
		}

		return Op("&").Add(code), nil
	}
}

func strToCode(g *Generator, val string) (Code, error) {
	// split longer strings as multiline strings
	hasNewlines := strings.Count(val, "\n") > 0
	hasOnlyNewlineAtEnd := strings.Count(val, "\n") == 1 && strings.HasSuffix(val, "\n")

	if hasNewlines && !hasOnlyNewlineAtEnd {
		return Id(fmt.Sprintf("`%s`", val)), nil
	}

	return Lit(val), nil
}

func litToCode(g *Generator, value reflect.Value) (Code, error) {
	val := value.Interface()

	switch val.(type) {
	case string:
		return strToCode(g, val.(string))
	default:
		return Lit(val), nil
	}
}

func marshalersToCode(g *Generator, value reflect.Value, parent reflect.Value) (Code, error) {
	switch v := value.Interface().(type) {
	case encoding.TextMarshaler:
		data, err := v.MarshalText()
		if err != nil {
			return nil, err
		}

		litValue, err := strToCode(g, string(data))
		if err != nil {
			return nil, err
		}

		return decodeValueToCode(g, litValue, value.Type()), nil
	case json.Marshaler:
		data, err := v.MarshalJSON()
		if err != nil {
			return nil, err
		}

		litValue, err := strToCode(g, string(data))
		if err != nil {
			return nil, err
		}

		return decodeValueToCode(g, litValue, value.Type()), nil
	default:
		return nil, nil
	}
}

func isEmptyStructSkipPrivateFields(structValue reflect.Value) bool {
	if structValue.Kind() == reflect.Ptr {
		structValue = structValue.Elem()
	}

	if structValue.Kind() != reflect.Struct {
		return false
	}

	structType := structValue.Type()

	for i := 0; i < structType.NumField(); i++ {
		fieldType := structType.Field(i)
		// if field is private, we cannot set it
		if unicode.IsLower(rune(fieldType.Name[0])) {
			continue
		}

		return false
	}

	return true
}

func valToCode(g *Generator, value reflect.Value, parent reflect.Value) (Code, error) {
	if value == (reflect.Value{}) {
		return Nil(), nil
	}

	switch value.Type().Kind() {
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
		// check if is empty struct, if we ignore private fields
		// and in such case apply marshallers
		if isEmptyStructSkipPrivateFields(value) {
			code, err := marshalersToCode(g, value, parent)
			if code != nil || err != nil {
				return code, err
			}
		}

		return ptrToCode(g, value, parent)
	case reflect.Struct:
		// check if is empty struct, if we ignore private fields
		// and in such case apply marshallers
		if isEmptyStructSkipPrivateFields(value) {
			code, err := marshalersToCode(g, value, parent)
			if code != nil || err != nil {
				return code, err
			}
		}

		return structToCode(g, value, parent)
	default:
		panic("unsupported kind")
	}
}
