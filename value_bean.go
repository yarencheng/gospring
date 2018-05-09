package gospring

import (
	"reflect"
)

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

func (bean *valueBean) GetScope() Scope {
	return Prototype
}

func (bean *valueBean) GetFactory() (interface{}, []BeanI) {

	m := reflect.ValueOf(bean).MethodByName("GetValue")

	return m.Interface(), []BeanI{}
}

func (bean *valueBean) GetFinalize() *string {
	return nil
}

func (bean *valueBean) GetInit() *string {
	return nil
}

func (bean *valueBean) GetProperty(name string) []BeanI {
	return nil
}

func (bean *valueBean) GetProperties() map[string][]BeanI {
	return map[string][]BeanI{}
}

func (bean *valueBean) GetType() reflect.Type {
	return reflect.TypeOf(bean.value)
}
