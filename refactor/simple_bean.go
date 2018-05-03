package refactor

import "reflect"

func Bean(tvpe interface{}) StructBeanI {
	i := "Init"
	f := "Finalize"
	return &simpleBean{
		tvpe:       reflect.TypeOf(tvpe),
		properties: make(map[string][]interface{}),
		init:       &i,
		finalize:   &f,
	}
}

type simpleBean struct {
	id          *string
	tvpe        reflect.Type
	properties  map[string][]interface{}
	factoryFn   interface{}
	factoryArgv []interface{}
	init        *string
	finalize    *string
}

func (bean *simpleBean) Factory(fn interface{}, argv ...interface{}) StructBeanI {
	bean.factoryFn = fn
	bean.factoryArgv = argv
	return bean
}

func (bean *simpleBean) Finalize(fnName string) StructBeanI {
	bean.finalize = &fnName
	return bean
}

func (bean *simpleBean) ID(id string) StructBeanI {
	bean.id = &id
	return bean
}

func (bean *simpleBean) Init(fnName string) StructBeanI {
	bean.init = &fnName
	return bean
}

func (bean *simpleBean) Property(name string, values ...interface{}) StructBeanI {
	bean.properties[name] = values
	return bean
}

func (bean *simpleBean) TypeOf(i interface{}) StructBeanI {
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
