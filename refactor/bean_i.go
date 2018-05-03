package refactor

import "reflect"

type BeanI interface {
	Factory(fn interface{}, argv ...interface{}) BeanI
	Finalize(fnName string) BeanI
	ID(id string) BeanI
	Init(fnName string) BeanI
	Property(name string, values ...interface{}) BeanI
	TypeOf(i interface{}) BeanI
	GetFactory() (interface{}, []interface{})
	GetFinalize() *string
	GetID() *string
	GetInit() *string
	GetProperty(name string) []interface{}
	GetType() reflect.Type
}