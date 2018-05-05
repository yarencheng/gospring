package refactor

import "reflect"

type referenceBean struct {
	id        string
	reference BeanI
}

func (bean *referenceBean) GetID() *string {
	return &bean.id
}
func (bean *referenceBean) GetScope() Scope {
	return Default
}

func (bean *referenceBean) GetFactory() (interface{}, []BeanI) {
	return nil, nil
}

func (bean *referenceBean) GetFinalize() *string {
	return nil
}

func (bean *referenceBean) GetReference() BeanI {
	return bean.reference
}

func (bean *referenceBean) SetReference(b BeanI) {
	bean.reference = b
}

func (bean *referenceBean) GetInit() *string {
	return nil
}

func (bean *referenceBean) GetProperty(name string) []BeanI {
	return nil
}

func (bean *referenceBean) GetProperties() map[string][]BeanI {
	return map[string][]BeanI{}
}

func (bean *referenceBean) GetType() reflect.Type {
	return nil
}
