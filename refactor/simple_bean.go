package refactor

type simpleBean struct {
	id          *string
	properties  map[string][]interface{}
	factoryFn   interface{}
	factoryArgv []interface{}
}

func (bean *simpleBean) Factory(fn interface{}, argv ...interface{}) BeanI {
	bean.factoryFn = fn
	bean.factoryArgv = argv
	return bean
}

func (bean *simpleBean) Finalize(fnName string) BeanI {
	return bean
}

func (bean *simpleBean) ID(id string) BeanI {
	bean.id = &id
	return bean
}

func (bean *simpleBean) Init(fnName string) BeanI {
	return bean
}

func (bean *simpleBean) Property(name string, values ...interface{}) BeanI {
	bean.properties[name] = values
	return bean
}

func (bean *simpleBean) GetFactory() (interface{}, []interface{}) {
	return bean.factoryFn, bean.factoryArgv
}

func (bean *simpleBean) GetID() *string {
	return bean.id
}

func (bean *simpleBean) GetInit() *string {
	return nil
}

func (bean *simpleBean) GetProperty(name string) []interface{} {
	value, _ := bean.properties[name]
	return value
}
