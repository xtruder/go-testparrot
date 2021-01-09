package recorder

import (
	"fmt"
	"reflect"
)

func getFieldName(structPtr interface{}, fieldPtr interface{}) (name string) {
	val := reflect.ValueOf(structPtr).Elem()
	val2 := reflect.ValueOf(fieldPtr).Elem()

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		if valueField.Addr().Interface() == val2.Addr().Interface() {
			return val.Type().Field(i).Name
		}
	}
	return
}

func getStructName(val interface{}) string {
	if t := reflect.TypeOf(val); t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	} else {
		return t.Name()
	}
}

func setStructValue(structPtr interface{}, fieldName string, val interface{}) error {
	structElem := reflect.ValueOf(structPtr).Elem()
	field := structElem.FieldByName(fieldName)

	if field == (reflect.Value{}) {
		return fmt.Errorf("invalid field: %s", fieldName)
	}

	field.Set(reflect.ValueOf(val))

	return nil
}
