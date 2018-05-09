package gospring

type StructBeanI interface {
	Factory(fn interface{}, argv ...interface{}) StructBeanI
	Finalize(fnName string) StructBeanI
	ID(id string) StructBeanI
	Init(fnName string) StructBeanI
	Property(name string, values ...interface{}) StructBeanI
	Prototype() StructBeanI
	Singleton() StructBeanI
	TypeOf(i interface{}) StructBeanI
}
