package testparrot

import (
	"io"
	"os"
	"reflect"
	"sort"

	. "github.com/dave/jennifer/jen"
)

const (
	valToPtrF     = "ValToPtr"
	headerComment = "Code generated by testparrot. DO NOT EDIT."
)

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

func (g *Generator) GenerateToFile(recorder *Recorder, recorderVar string, filePath string) error {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0660)
	if err != nil {
		return err
	}

	if err := g.Generate(recorder, recorderVar, file); err != nil {
		return err
	}

	return file.Close()
}

func (g *Generator) Generate(recorder *Recorder, recorderVar string, out io.Writer) error {
	f := NewFilePathName(g.pkgPath, g.pkgName)

	f.HeaderComment(headerComment)

	keys := make([]string, 0, len(recorder.allRecordings))

	for key := range recorder.allRecordings {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	statements := []Code{}
	for _, key := range keys {
		testKey := key
		recordings := recorder.allRecordings[key]

		// call load on global recorder or on locally defined recorder
		var loadF *Statement
		if recorder == R {
			loadF = Qual(pkgPath, "R.Load")
		} else {
			loadF = Id(recorderVar + ".Load")
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
	typ := sliceVal.Type().Elem()

	values := []Code{}
	for i := 0; i < sliceVal.Len(); i++ {
		elem := sliceVal.Index(i)

		code, err := valToCode(g, elem, sliceVal)
		if err != nil {
			return nil, err
		}

		values = append(values, code)
	}

	return Index().Add(typeToCode(g, typ)).Values(values...), nil
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
		return valToPtrF.Call(Lit(val.Interface())), nil
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

func valToCode(g *Generator, value reflect.Value, parent reflect.Value) (Code, error) {
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
		return Lit(value.Interface()), nil
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