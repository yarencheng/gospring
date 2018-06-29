package interfaces

import (
	uuid "github.com/satori/go.uuid"
)

type ApplicationContextI interface {
	GetByID(id string) (interface{}, error)
	GetByUUID(uuid uuid.UUID) (interface{}, error)
	GetBeanByID(id string) (BeanI, bool)
	GetBeanByUUID(uuid uuid.UUID) (BeanI, bool)
	AddConfig(config interface{}) (BeanI, error)
}
