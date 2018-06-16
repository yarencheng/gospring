package bean

import (
	"fmt"
	"reflect"

	"github.com/yarencheng/gospring/v1"
)

type StructBean struct {
	id   string
	tvpe reflect.Type
}

func NewStructBeanV1(config v1.Bean) (*StructBean, error) {

	switch config.Type.Kind() {
	case reflect.Uintptr |
		reflect.Array |
		reflect.Chan |
		reflect.Func |
		reflect.Interface |
		reflect.Map |
		reflect.Slice |
		reflect.Ptr |
		reflect.UnsafePointer:
		return nil, fmt.Errorf("[%v] is not a valid type for a bean", config.Type.Kind())
	}
	return &StructBean{
		id:   config.ID,
		tvpe: config.Type,
	}, nil
}

func (b *StructBean) GetID() string {
	return b.id
}

func (b *StructBean) GetValue() (reflect.Value, error) {
	return reflect.New(b.tvpe), nil
}
