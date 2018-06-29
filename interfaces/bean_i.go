package interfaces

import (
	"context"
	"reflect"

	uuid "github.com/satori/go.uuid"
)

type BeanI interface {
	GetUUID() uuid.UUID
	GetID() string
	GetValue() (reflect.Value, error)
	Stop(ctx context.Context) error
}
