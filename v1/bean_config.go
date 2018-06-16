package v1

import (
	"reflect"
)

type Bean struct {
	ID    string
	Type  reflect.Type
	Scope Scope
}

func T(i interface{}) reflect.Type {
	return reflect.TypeOf(i)
}
