package bean

import (
	"context"
	"reflect"
)

type BeanI interface {
	GetID() string
	GetValue() (reflect.Value, error)
	Stop(ctx context.Context) error
}
