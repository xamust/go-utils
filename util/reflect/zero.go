package reflect

import (
	"errors"
	"reflect"
)

var (
	errInvalidType = errors.New("invalid type struct")
)

func NewZeroValue(src any) (any, error) {
	val := reflect.ValueOf(src)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() == reflect.Invalid {
		return nil, errInvalidType
	}

	result := reflect.New(val.Type())

	return result.Interface(), nil
}
