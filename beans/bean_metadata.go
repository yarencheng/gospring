package beans

import "reflect"

type BeanMetaData struct {
	id     string
	scope  Scope
	ztruct reflect.Type
}

func NewBeanMetaData(id string, scope Scope, ztruct reflect.Type) *BeanMetaData {
	return &BeanMetaData{
		id:     id,
		scope:  scope,
		ztruct: ztruct,
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
