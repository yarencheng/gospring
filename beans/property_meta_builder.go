package beans

import "fmt"

type propertyRef struct {
	name string
	ref  string
}

func (p *propertyRef) GetName() string {
	return p.name
}

func (p *propertyRef) GetReference() string {
	return p.ref
}

func (p *propertyRef) GetBean() BeanMetaData {
	return nil
}

func (p *propertyRef) IsReference() bool {
	return true
}

type propertyBean struct {
	name string
	bean BeanMetaData
}

func (p *propertyBean) GetName() string {
	return p.name
}

func (p *propertyBean) GetReference() string {
	return ""
}

func (p *propertyBean) GetBean() BeanMetaData {
	return p.bean
}

func (p *propertyBean) IsReference() bool {
	return false
}

type propertyMetaBuilder struct {
	name  string
	ref   string
	bean  beanMetaBuilder
	isRef bool
}

func ProperRef(name string, ref string) *propertyMetaBuilder {
	var b propertyMetaBuilder
	b.isRef = true
	b.name = name
	b.ref = ref
	return &b
}

func ProperBean(name string, bean beanMetaBuilder) *propertyMetaBuilder {
	var b propertyMetaBuilder
	b.isRef = false
	b.name = name
	b.bean = bean
	return &b
}

func (b *propertyMetaBuilder) Build() (PropertyMetaData, error) {
	if b.isRef {
		var p propertyRef
		p.name = b.name
		p.ref = b.ref
		return &p, nil
	} else {
		var p propertyBean
		p.name = b.name
		var e error
		p.bean, e = b.bean.Build()
		if e != nil {
			return nil, fmt.Errorf("can't build metadata of bean [%v]. Caused by: %v", b.bean, e)
		}
		return &p, nil
	}
}
