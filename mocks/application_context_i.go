package mocks

import (
	"context"
	"reflect"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/mock"
	"github.com/yarencheng/gospring/interfaces"
)

type ApplicationContextMock struct {
	mock.Mock
}

func (m *ApplicationContextMock) GetByID(id string) (interface{}, error) {
	args := m.Called(id)
	return args.Get(0), args.Error(1)
}

func (m *ApplicationContextMock) GetByUUID(uuid uuid.UUID) (interface{}, error) {
	args := m.Called(uuid)
	return args.Get(0), args.Error(1)
}

func (m *ApplicationContextMock) GetBeanByID(id string) (interfaces.BeanI, bool) {
	args := m.Called(id)

	v0 := args.Get(0)
	var r0 interfaces.BeanI
	if v0 != nil {
		r0 = v0.(interfaces.BeanI)
	}

	return r0, args.Bool(1)
}

func (m *ApplicationContextMock) GetBeanByUUID(uuid uuid.UUID) (interfaces.BeanI, bool) {
	args := m.Called(uuid)

	v0 := args.Get(0)
	var r0 interfaces.BeanI
	if v0 != nil {
		r0 = v0.(interfaces.BeanI)
	}

	return r0, args.Bool(1)
}

func (m *ApplicationContextMock) AddConfig(config interface{}) (interfaces.BeanI, error) {
	args := m.Called(config)

	v0 := args.Get(0)
	var r0 interfaces.BeanI
	if v0 != nil {
		r0 = v0.(interfaces.BeanI)
	}

	return r0, args.Error(1)
}

func (m *ApplicationContextMock) UseConfigParser(configType reflect.Type, parser interfaces.ConfigParser) error {
	args := m.Called(configType, parser)

	return args.Error(0)
}

func (m *ApplicationContextMock) Stop(ctx context.Context) error {
	args := m.Called(ctx)

	return args.Error(0)
}
