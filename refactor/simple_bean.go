package refactor

type simpleBean struct {
	id         *string
	properties map[string][]interface{}
}

func (bean *simpleBean) Factory(fn interface{}, argv ...interface{}) BeanI {
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

func (bean *simpleBean) GetFactory() interface{} {
	return nil
}

func (bean *simpleBean) GetFactoryArgv() []interface{} {
	return nil
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
