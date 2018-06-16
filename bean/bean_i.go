package bean

import "reflect"

type BeanI interface {
	GetID() string
	GetValue() (reflect.Value, error)
}
