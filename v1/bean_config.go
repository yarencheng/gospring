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
	Name   string
	Config interface{}
}

type Channel struct {
	ID   string
	Type reflect.Type
	Size int
}

type Broadcast struct {
	ID       string
	SourceID string
	Size     int
}

type Ref struct {
	ID string
}

func T(i interface{}) reflect.Type {
	return reflect.TypeOf(i)
}
