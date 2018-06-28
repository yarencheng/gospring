package v1

import (
	"reflect"
)

type Bean struct {
	ID          string
	Type        reflect.Type
	Scope       Scope
	FactoryFn   interface{}
	FactoryArgs []interface{}
	StartFn     interface{}
	StopFn      interface{}
	Properties  []Property
}

type Property struct {
	Name  string
	Value interface{}
}

func T(i interface{}) reflect.Type {
	return reflect.TypeOf(i)
}
