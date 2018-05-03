package refactor

type simpleBean struct {
}

func (bean *simpleBean) Factory(fn interface{}, argv ...interface{}) BeanI {
	return bean
}

func (bean *simpleBean) Finalize(fnName string) BeanI {
	return bean
}

func (bean *simpleBean) ID(id string) BeanI {
	return bean
}

func (bean *simpleBean) Init(fnName string) BeanI {
	return bean
}

func (bean *simpleBean) Property(name string, values ...interface{}) BeanI {
	return bean
}

func (bean *simpleBean) GetFactory() interface{} {
	return nil
}

func (bean *simpleBean) GetFactoryArgv() []interface{} {
	return nil
}

func (bean *simpleBean) GetID() *string {
	return nil
}

func (bean *simpleBean) GetInit() *string {
	return nil
}

func (bean *simpleBean) GetProperty(name string) []interface{} {
	return nil
}
