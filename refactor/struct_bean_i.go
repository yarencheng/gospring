package refactor

import "reflect"

type Scope string

const (
	Default   Scope = "Default"
	Singleton Scope = "Singleton"
	Prototype Scope = "Prototype"
)

type StructBeanI interface {
	GetFactory() (interface{}, []BeanI)
	GetFinalize() *string
	GetID() *string
	GetInit() *string
	GetProperty(name string) []BeanI
	GetProperties() map[string][]BeanI
	GetScope() Scope
	GetType() reflect.Type
	Factory(fn interface{}, argv ...interface{}) StructBeanI
	Finalize(fnName string) StructBeanI
	ID(id string) StructBeanI
	Init(fnName string) StructBeanI
	Property(name string, values ...interface{}) StructBeanI
	Prototype() StructBeanI
	Singleton() StructBeanI
	TypeOf(i interface{}) StructBeanI
}
