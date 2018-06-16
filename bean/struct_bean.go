package bean

import (
	"reflect"

	"github.com/yarencheng/gospring/v1"
)

type StructBean struct {
	id   string
	tvpe reflect.Type
}

func NewStructBeanV1(config v1.Bean) (*StructBean, error) {
	return &StructBean{
		id:   config.ID,
		tvpe: config.Type,
	}, nil
}

func (b *StructBean) GetID() string {
	return b.id
}
