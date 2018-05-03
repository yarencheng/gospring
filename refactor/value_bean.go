package refactor

type valueBean struct {
	id    *string
	value interface{}
}

func (bean *valueBean) GetID() *string {
	return bean.id
}

func (bean *valueBean) GetValue() interface{} {
	return bean.value
}
