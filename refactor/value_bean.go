package refactor

import "reflect"

type valueBean struct {
	value interface{}
}

func (bean *valueBean) GetID() *string {
	return nil
}

func (bean *valueBean) GetValue() interface{} {

	tvpe := reflect.TypeOf(bean.value)
	valuePtr := reflect.New(tvpe)

	value := reflect.ValueOf(bean.value)
	valuePtr.Elem().Set(value)

	return valuePtr.Interface()
}
