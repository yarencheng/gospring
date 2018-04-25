package beans

import "reflect"

type PropertyMetaData struct {
	name string
}

func NewPropertyMetaData(name string, scope Scope, ztruct reflect.Type) *PropertyMetaData {
	return &PropertyMetaData{
		name: name,
	}
}

func (meta *PropertyMetaData) GetName() string {
	return meta.name
}
