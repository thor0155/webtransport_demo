package utils

import (
	"reflect"

	"go.uber.org/zap"
)

func ObjectToZapFields(obj any) []zap.Field {
	v := reflect.ValueOf(obj)
	t := reflect.TypeOf(obj)

	// 若是指標，先 dereference
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
		t = t.Elem()
	}

	// 必須是 struct
	if v.Kind() != reflect.Struct {
		panic("ObjectToZapFields: obj must be a struct or pointer to struct")
	}

	fields := make([]zap.Field, 0, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		fieldVal := v.Field(i)
		fieldType := t.Field(i)

		tag := fieldType.Tag.Get("mapstructure")
		if tag == "" {
			tag = fieldType.Name
		}

		// 跳過零值
		if IsZero(fieldVal) {
			continue
		}

		fields = append(fields, zap.Any(tag, fieldVal.Interface()))
	}

	return fields
}
