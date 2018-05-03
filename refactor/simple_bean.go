package refactor

import "reflect"

type simpleBean struct {
	id          *string
	tvpe        reflect.Type
	properties  map[string][]interface{}
	factoryFn   interface{}
	factoryArgv []interface{}
	init        *string
	finalize    *string
}

func (bean *simpleBean) Factory(fn interface{}, argv ...interface{}) BeanI {
	bean.factoryFn = fn
	bean.factoryArgv = argv
	return bean
}

func (bean *simpleBean) Finalize(fnName string) BeanI {
	bean.finalize = &fnName
	return bean
}

func (bean *simpleBean) ID(id string) BeanI {
	bean.id = &id
	return bean
}

func (bean *simpleBean) Init(fnName string) BeanI {
	bean.init = &fnName
	return bean
}

func (bean *simpleBean) Property(name string, values ...interface{}) BeanI {
	bean.properties[name] = values
	return bean
}

func (bean *simpleBean) TypeOf(i interface{}) BeanI {
	bean.tvpe = reflect.TypeOf(i)
	return bean
}

func (bean *simpleBean) GetFactory() (interface{}, []interface{}) {
	return bean.factoryFn, bean.factoryArgv
}

func (bean *simpleBean) GetFinalize() *string {
	return bean.finalize
}

func (bean *simpleBean) GetID() *string {
	return bean.id
}

func (bean *simpleBean) GetInit() *string {
	return bean.init
}

func (bean *simpleBean) GetProperty(name string) []interface{} {
	value, _ := bean.properties[name]
	return value
}

func (bean *simpleBean) GetType() reflect.Type {
	return bean.tvpe
}
