package refactor

type BeanI interface {
	Factory(fn interface{}, argv ...interface{}) BeanI
	Finalize(fnName string) BeanI
	Id(id string) BeanI
	Init(fnName string) BeanI
	Property(name string, values ...interface{}) BeanI
	GetFactory() interface{}
	GetFactoryArgv() []interface{}
	GetId() *string
	GetInit() *string
	GetProperty(name string) []interface{}
}
