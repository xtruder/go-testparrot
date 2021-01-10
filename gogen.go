package testparrot

import (
	"io"
	"reflect"

	. "github.com/dave/jennifer/jen"
)

const (
	pkgPath   = "github.com/xtruder/go-testparrot"
	pkgName   = "testparrot"
	valToPtrF = "ValToPtr"
	loadF     = "Load"
)

// GoGenerator generates golang code
type GoGenerator struct {
	// Name of generate package
	PkgName string
}

func (g *GoGenerator) Generate(registry *Registry, out io.Writer) error {
	f := NewFilePath(g.PkgName)

	recorders := registry.Recorders()

	statements := []Code{}
	for _, recorder := range recorders {
		loadF := Qual(pkgPath, loadF)

		val, err := valToCode(g, reflect.ValueOf(recorder.Records()))
		if err != nil {
			return err
		}

		loadCall := loadF.Call(Lit(recorder.Name()), val)
		statements = append(statements, loadCall)
	}

	// Create init method
	f.Func().Id("init").Params().Block(statements...)

	// Render code
	return f.Render(out)
}

func typeToCode(g *GoGenerator, typ reflect.Type) *Statement {
	switch typ.Kind() {
	case reflect.Interface:
		return Interface()
	default:
		pkgPath := typ.PkgPath()
		if typ.PkgPath() == "" || g.PkgName == pkgPath {
			return Id(typ.Name())
		}

		return Qual(typ.PkgPath(), typ.Name())
	}
}

func sliceToCode(g *GoGenerator, sliceVal reflect.Value) (Code, error) {
	typ := sliceVal.Type().Elem()

	values := []Code{}
	for i := 0; i < sliceVal.Len(); i++ {
		elem := sliceVal.Index(i)

		code, err := valToCode(g, elem)
		if err != nil {
			return nil, err
		}

		values = append(values, code)
	}

	return Index().Add(typeToCode(g, typ)).Values(values...), nil
}

func mapToCode(g *GoGenerator, mapVal reflect.Value) (Code, error) {
	values := Dict{}

	var keyType reflect.Type
	var valType reflect.Type
	for _, k := range mapVal.MapKeys() {
		v := mapVal.MapIndex(k)

		keyCode, err := valToCode(g, k)
		if err != nil {
			return nil, err
		}

		valCode, err := valToCode(g, v)
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

func structToCode(g *GoGenerator, structVal reflect.Value) (Code, error) {
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

		code, err := valToCode(g, val)
		if err != nil {
			return nil, err
		}

		values[Id(field.Name)] = code
	}

	return typeToCode(g, structType).Values(values), nil
}

func ptrToCode(g *GoGenerator, ptrVal reflect.Value) (Code, error) {
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
		return ptrToCode(g, val)
	case reflect.Ptr:
		code, err := ptrToCode(g, val)
		if err != nil {
			return nil, err
		}

		return valToPtrF.Call(code), nil
	default:
		code, err := valToCode(g, val)
		if err != nil {
			return nil, err
		}

		return Op("&").Add(code), nil
	}
}

func valToCode(g *GoGenerator, value reflect.Value) (Code, error) {
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
		return sliceToCode(g, value)
	case reflect.Interface:
		return valToCode(g, value.Elem())
	case reflect.Map:
		return mapToCode(g, value)
	case reflect.Ptr:
		return ptrToCode(g, value)
	case reflect.Struct:
		return structToCode(g, value)
	default:
		panic("unsupported kind")
	}
}
