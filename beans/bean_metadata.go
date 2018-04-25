package beans

import "reflect"

type BeanMetaData struct {
	id         string
	scope      Scope
	ztruct     reflect.Type
	properties []PropertyMetaData
}

func NewBeanMetaData(id string, scope Scope, ztruct reflect.Type, properties []PropertyMetaData) *BeanMetaData {
	return &BeanMetaData{
		id:         id,
		scope:      scope,
		ztruct:     ztruct,
		properties: properties,
	}
}

func (meta *BeanMetaData) GetId() string {
	return meta.id
}

func (meta *BeanMetaData) GetScope() Scope {
	return meta.scope
}

func (meta *BeanMetaData) GetStruct() reflect.Type {
	return meta.ztruct
}

func (meta *BeanMetaData) GetProperties() []PropertyMetaData {
	tmp := make([]PropertyMetaData, len(meta.properties))
	copy(tmp, meta.properties)
	return tmp
}
