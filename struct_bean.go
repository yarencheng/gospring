package gospring

import "reflect"

const (
	DefaultInitFunc     string = "Init"
	DefaultFinalizeFunc string = "Finalize"
)

type structBean struct {
	id          *string
	tvpe        reflect.Type
	value       reflect.Value
	properties  map[string][]BeanI
	factoryFn   interface{}
	factoryArgv []BeanI
	init        *string
	finalize    *string
	scope       Scope
}

func (bean *structBean) Factory(fn interface{}, argv ...interface{}) StructBeanI {
	bean.factoryFn = fn
	bean.factoryArgv = Beans(argv...)
	return bean
}

func (bean *structBean) Finalize(fnName string) StructBeanI {
	bean.finalize = &fnName
	return bean
}

func (bean *structBean) ID(id string) StructBeanI {
	bean.id = &id
	return bean
}

func (bean *structBean) Init(fnName string) StructBeanI {
	bean.init = &fnName
	return bean
}

func (bean *structBean) Property(name string, values ...interface{}) StructBeanI {
	bean.properties[name] = Beans(values...)
	return bean
}

func (bean *structBean) TypeOf(i interface{}) StructBeanI {
	bean.tvpe = reflect.TypeOf(i)
	return bean
}

func (bean *structBean) GetFactory() (interface{}, []BeanI) {
	return bean.factoryFn, bean.factoryArgv
}

func (bean *structBean) GetFinalize() *string {
	return bean.finalize
}

func (bean *structBean) GetID() *string {
	return bean.id
}

func (bean *structBean) GetInit() *string {
	return bean.init
}

func (bean *structBean) GetProperty(name string) []BeanI {
	value, _ := bean.properties[name]
	return value
}

func (bean *structBean) GetProperties() map[string][]BeanI {
	return bean.properties
}

func (bean *structBean) GetType() reflect.Type {
	return bean.tvpe
}

func (bean *structBean) Prototype() StructBeanI {
	bean.scope = Prototype
	return bean
}
func (bean *structBean) Singleton() StructBeanI {
	bean.scope = Singleton
	return bean
}

func (bean *structBean) GetScope() Scope {
	return bean.scope
}
