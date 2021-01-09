package generator

import (
	"reflect"

	. "github.com/dave/jennifer/jen"
)

type generator struct {
	// Name of generate package
	PkgName string
}

func typeToCode(g *generator, typ reflect.Type) *Statement {
	switch typ.Kind() {
	case reflect.Interface:
		return Id("interface{}")
	default:
		if typ.PkgPath() == "" || g.PkgName == typ.PkgPath() {
			return Id(typ.Name())
		}

		return Qual(typ.PkgPath(), typ.Name())
	}
}

func sliceToCode(g *generator, valValue reflect.Value) (Code, error) {
	sliceType := valValue.Type().Elem()
	typeName := sliceType.Name()

	values := []Code{}
	for i := 0; i < valValue.Len(); i++ {
		elem := valValue.Index(i)

		code, err := valToCode(g, elem)
		if err != nil {
			return nil, err
		}

		values = append(values, code)
	}

	return Index().Id(typeName).Values(values...), nil
}

func mapToCode(g *generator, valValue reflect.Value) (Code, error) {
	values := Dict{}

	var keyType reflect.Type
	var valType reflect.Type
	for _, k := range valValue.MapKeys() {
		v := valValue.MapIndex(k)

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

func structToCode(g *generator, structVal reflect.Value) (Code, error) {
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

		values[Lit(field.Name)] = code
	}

	return typeToCode(g, structType).Values(values), nil
}

func ptrToCode(g *generator, ptrValue reflect.Value) (Code, error) {
	valValue := ptrValue.Elem()
	valType := valValue.Type()

	valToPtrF := Qual("github.com/xtruder/go-test-recorder/generator", "ValToPtr")

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
		return valToPtrF.Call(Lit(valValue.Interface())), nil
	case reflect.Interface:
		return ptrToCode(g, valValue)
	case reflect.Ptr:
		code, err := ptrToCode(g, valValue)
		if err != nil {
			return nil, err
		}

		return valToPtrF.Call(code), nil
	default:
		code, err := valToCode(g, valValue)
		if err != nil {
			return nil, err
		}

		return Op("&").Add(code), nil
	}
}

func valToCode(g *generator, valValue reflect.Value) (Code, error) {
	valType := valValue.Type()

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
		return Lit(valValue.Interface()), nil
	case reflect.Array, reflect.Slice:
		return sliceToCode(g, valValue)
	case reflect.Interface:
		return valToCode(g, valValue.Elem())
	case reflect.Map:
		return mapToCode(g, valValue)
	case reflect.Ptr:
		return ptrToCode(g, valValue)
	case reflect.Struct:
		return structToCode(g, valValue)
	default:
		panic("unsupported kind")
	}
}
