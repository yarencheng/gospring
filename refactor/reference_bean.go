package refactor

type referenceBean struct {
	targetId string
}

func (bean *referenceBean) Factory(fn interface{}, argv ...interface{}) BeanI {
	return bean
}

func (bean *referenceBean) Finalize(fnName string) BeanI {
	return bean
}

func (bean *referenceBean) ID(id string) BeanI {
	return bean
}

func (bean *referenceBean) Init(fnName string) BeanI {
	return bean
}

func (bean *referenceBean) Property(name string, values ...interface{}) BeanI {
	return bean
}

func (bean *referenceBean) GetFactory() (interface{}, []interface{}) {
	return nil, []interface{}{}
}

func (bean *referenceBean) GetFinalize() *string {
	return nil
}

func (bean *referenceBean) GetID() *string {
	return &bean.targetId
}

func (bean *referenceBean) GetInit() *string {
	return nil
}

func (bean *referenceBean) GetProperty(name string) []interface{} {
	return nil
}
