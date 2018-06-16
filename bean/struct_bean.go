package bean

import "github.com/yarencheng/gospring/v1"

type StructBean struct {
	id string
}

func NewStructBeanV1(b v1.Bean) *StructBean {
	return &StructBean{}
}

func (b *StructBean) GetID() string {
	return b.id
}
