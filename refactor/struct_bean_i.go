package refactor

import "reflect"

type StructBeanI interface {
	Factory(fn interface{}, argv ...interface{}) StructBeanI
	Finalize(fnName string) StructBeanI
	ID(id string) StructBeanI
	Init(fnName string) StructBeanI
	Property(name string, values ...interface{}) StructBeanI
	TypeOf(i interface{}) StructBeanI
	GetFactory() (interface{}, []interface{})
	GetFinalize() *string
	GetID() *string
	GetInit() *string
	GetProperty(name string) []interface{}
	GetType() reflect.Type
}
