package interfaces

import (
	"context"
	"reflect"

	uuid "github.com/satori/go.uuid"
)

type ApplicationContextI interface {
	GetByID(id string) (interface{}, error)
	GetByUUID(uuid uuid.UUID) (interface{}, error)
	AddConfig(config interface{}) (BeanI, error)
	UseConfigParser(configType reflect.Type, parser ConfigParser) error
	Stop(ctx context.Context) error
}

type ConfigParser func(ctx ApplicationContextI, config interface{}) (BeanI, error)
