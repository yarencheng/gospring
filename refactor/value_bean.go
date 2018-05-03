package refactor

type valueBean struct {
	value interface{}
}

func (bean *valueBean) GetID() *string {
	return nil
}

func (bean *valueBean) GetValue() interface{} {
	return bean.value
}
