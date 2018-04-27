package beans

import (
	"fmt"
	"reflect"
)

type bean struct {
	id         string
	properties []PropertyMetaData
	scope      Scope
	type_      reflect.Type
}

func (b *bean) GetId() string {
	return b.id
}

func (b *bean) GetProperties() []PropertyMetaData {
	tmp := make([]PropertyMetaData, len(b.properties))
	copy(tmp, b.properties)
	return tmp
}

func (b *bean) GetScope() Scope {
	return b.scope
}

func (b *bean) GetType() reflect.Type {
	return b.type_
}

type beanMetasBuilder struct {
	beans []*beanMetaBuilder
}

func Beans(bs ...*beanMetaBuilder) *beanMetasBuilder {
	var builder beanMetasBuilder
	builder.beans = bs
	return &builder
}

func (b *beanMetasBuilder) Build() ([]BeanMetaData, error) {
	beans := make([]BeanMetaData, len(b.beans))
	for i, bean := range b.beans {
		meta, e := bean.Build()
		if e != nil {
			return nil, e
		}
		beans[i] = meta
	}

	// TODO: sanity check

	return beans, nil
}

type beanMetaBuilder struct {
	id         string
	name       string
	properties []propertyMetaBuilder
	scope      Scope
	type_      reflect.Type
}

func Bean() *beanMetaBuilder {
	var b beanMetaBuilder

	b.id = ""
	b.name = ""
	b.properties = make([]propertyMetaBuilder, 0)
	b.scope = Singleton

	return &b
}

func (b *beanMetaBuilder) ID(id string) *beanMetaBuilder {
	b.id = id
	return b
}

func (b *beanMetaBuilder) Name(name string) *beanMetaBuilder {
	b.name = name
	return b
}

func (b *beanMetaBuilder) Property(ps ...propertyMetaBuilder) *beanMetaBuilder {
	b.properties = append(b.properties, ps...)
	return b
}

func (b *beanMetaBuilder) Singleton() *beanMetaBuilder {
	b.scope = Singleton
	return b
}

func (b *beanMetaBuilder) Prototype() *beanMetaBuilder {
	b.scope = Prototype
	return b
}

func (b *beanMetaBuilder) Type(type_ reflect.Type) *beanMetaBuilder {
	b.type_ = type_
	return b
}

func (b *beanMetaBuilder) Build() (BeanMetaData, error) {

	var bean bean

	bean.id = b.id
	bean.scope = b.scope
	bean.type_ = b.type_

	tmp := make([]PropertyMetaData, len(b.properties))
	for i, pb := range b.properties {
		var e error
		tmp[i], e = pb.Build()
		if e != nil {
			return nil, fmt.Errorf("can't build metadata of property [%v]. Caused by: %v", pb, e)
		}
	}

	// TODO: sanity check

	return &bean, nil
}
