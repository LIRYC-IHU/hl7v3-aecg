package types

import "reflect"

func BoolPtr(b bool) *bool {
	return &b
}

func IsPointer(v any) bool {
	return reflect.TypeOf(v).Kind() == reflect.Pointer
}
