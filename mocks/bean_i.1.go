package mocks

import (
	"context"
	"reflect"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/mock"
)

type BeanI interface {
	GetUUID() uuid.UUID
	GetID() string
	GetValue() (reflect.Value, error)
	Stop(ctx context.Context) error
}

type BeanMock struct {
	mock.Mock
}

func (m *BeanMock) GetUUID() uuid.UUID {
	args := m.Called()

	v0 := args.Get(0)
	var r0 uuid.UUID
	if v0 != nil {
		r0 = v0.(uuid.UUID)
	}

	return r0
}

func (m *BeanMock) GetID() string {
	args := m.Called()
	return args.String(0)
}

func (m *BeanMock) GetValue() (reflect.Value, error) {
	args := m.Called()

	v0 := args.Get(0)
	var r0 reflect.Value
	if v0 != nil {
		r0 = v0.(reflect.Value)
	}

	return r0, args.Error(1)
}

func (m *BeanMock) Stop(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}
