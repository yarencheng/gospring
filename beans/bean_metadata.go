package beans

import "reflect"

type BeanMetaData interface {
	GetId() string
	GetProperties() []PropertyMetaData
	GetScope() Scope
	GetType() reflect.Type
}

type BeanMetaData_old struct {
	id         string
	scope      Scope
	ztruct     reflect.Type
	properties []PropertyMetaData_old
}

func NewBeanMetaData_old(id string, scope Scope, ztruct reflect.Type, properties []PropertyMetaData_old) *BeanMetaData_old {
	return &BeanMetaData_old{
		id:         id,
		scope:      scope,
		ztruct:     ztruct,
		properties: properties,
	}
}

func (meta *BeanMetaData_old) GetId() string {
	return meta.id
}

func (meta *BeanMetaData_old) GetScope() Scope {
	return meta.scope
}

func (meta *BeanMetaData_old) GetStruct() reflect.Type {
	return meta.ztruct
}

func (meta *BeanMetaData_old) GetProperties() []PropertyMetaData_old {
	tmp := make([]PropertyMetaData_old, len(meta.properties))
	copy(tmp, meta.properties)
	return tmp
}
