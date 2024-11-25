package utils

import "reflect"

func ValueIsZero(value any) bool {
	return reflect.ValueOf(&value).Elem().IsZero()
}
