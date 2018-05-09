package gospring

import "reflect"

type Scope string

const (
	Default   Scope = "Default"
	Singleton Scope = "Singleton"
	Prototype Scope = "Prototype"
)

type BeanI interface {
	GetFactory() (interface{}, []BeanI)
	GetFinalize() *string
	GetID() *string
	GetInit() *string
	GetProperty(name string) []BeanI
	GetProperties() map[string][]BeanI
	GetScope() Scope
	GetType() reflect.Type
}
