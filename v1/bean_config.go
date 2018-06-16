package v1

import "reflect"

type Bean struct {
	ID   string
	Type reflect.Type
}

func T(i interface{}) reflect.Type {
	return reflect.TypeOf(i)
}
