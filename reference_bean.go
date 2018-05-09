package gospring

import "reflect"

type referenceBean struct {
	id        string
	reference BeanI
}

func (bean *referenceBean) GetID() *string {
	if bean.reference == nil {
		return &bean.id
	} else {
		return bean.reference.GetID()
	}
}
func (bean *referenceBean) GetScope() Scope {
	if bean.reference == nil {
		return Default
	} else {
		return bean.reference.GetScope()
	}
}

func (bean *referenceBean) GetFactory() (interface{}, []BeanI) {
	if bean.reference == nil {
		return nil, nil
	} else {
		return bean.reference.GetFactory()
	}
}

func (bean *referenceBean) GetFinalize() *string {
	if bean.reference == nil {
		return nil
	} else {
		return bean.reference.GetFinalize()
	}
}

func (bean *referenceBean) GetReference() BeanI {
	return bean.reference
}

func (bean *referenceBean) GetInit() *string {
	if bean.reference == nil {
		return nil
	} else {
		return bean.reference.GetInit()
	}
}

func (bean *referenceBean) GetProperty(name string) []BeanI {
	if bean.reference == nil {
		return nil
	} else {
		return bean.reference.GetProperty(name)
	}
}

func (bean *referenceBean) GetProperties() map[string][]BeanI {
	if bean.reference == nil {
		return map[string][]BeanI{}
	} else {
		return bean.reference.GetProperties()
	}
}

func (bean *referenceBean) GetType() reflect.Type {
	if bean.reference == nil {
		return nil
	} else {
		return bean.reference.GetType()
	}
}

func (bean *referenceBean) SetReference(b BeanI) {
	bean.reference = b
}
