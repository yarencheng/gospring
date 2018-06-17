package bean

import (
	"fmt"
	"reflect"

	"github.com/yarencheng/gospring/v1"
)

type StructBean struct {
	id             string
	tvpe           reflect.Type
	scope          v1.Scope
	singletonValue reflect.Value
}

var defaultStruct StructBean = StructBean{
	scope: v1.Default,
}

func NewStructBeanV1(config v1.Bean) (*StructBean, error) {

	switch config.Type.Kind() {
	case reflect.Uintptr:
		fallthrough
	case reflect.Array:
		fallthrough
	case reflect.Chan:
		fallthrough
	case reflect.Func:
		fallthrough
	case reflect.Interface:
		fallthrough
	case reflect.Map:
		fallthrough
	case reflect.Slice:
		fallthrough
	case reflect.Ptr:
		fallthrough
	case reflect.UnsafePointer:
		return nil, fmt.Errorf("[%v] is not a valid type for a bean", config.Type.Kind())
	}

	scope := v1.Default
	if config.Scope == "" {
		scope = v1.Default
	}

	return &StructBean{
		id:    config.ID,
		tvpe:  config.Type,
		scope: scope,
	}, nil
}

func (b *StructBean) GetID() string {
	return b.id
}

func (b *StructBean) GetValue() (reflect.Value, error) {
	switch b.scope {
	case v1.Default:
		fallthrough
	case v1.Singleton:
		if b.singletonValue.IsValid() {
			return b.singletonValue, nil
		}
		v, err := b.createValue()
		if err != nil {
			return reflect.Value{}, fmt.Errorf("Can't create the value, err: %v", err)
		}
		b.singletonValue = v
		return b.singletonValue, nil
	case v1.Prototype:
		return b.createValue()
	default:
		return reflect.Value{}, fmt.Errorf("Unknown scope [%v]", b.scope)
	}
}

func (b *StructBean) createValue() (reflect.Value, error) {
	return reflect.New(b.tvpe), nil
}
