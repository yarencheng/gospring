package interfaces

import (
	uuid "github.com/satori/go.uuid"
)

type ApplicationContextI interface {
	GetByID(id string) (interface{}, error)
	GetByUUID(uuid uuid.UUID) (interface{}, error)
	AddConfig(config interface{}) (BeanI, error)
}
