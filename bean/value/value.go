package value

import (
	"context"
	"fmt"
	"reflect"

	uuid "github.com/satori/go.uuid"
	"github.com/yarencheng/gospring/interfaces"
	"github.com/yarencheng/gospring/v1"
)

type Value struct {
	uuid  uuid.UUID
	id    string
	value interface{}
}

func V1ValueParser(ctx interfaces.ApplicationContextI, config interface{}) (interfaces.BeanI, error) {
	c, ok := config.(*v1.Value)
	if !ok {
		return nil, fmt.Errorf("[%v] can not convert to *v1.Value", reflect.TypeOf(config))
	}

	b := &Value{
		uuid:  uuid.NewV4(),
		id:    c.ID,
		value: c.Value,
	}

	return b, nil
}

func (v *Value) GetUUID() uuid.UUID {
	return v.uuid
}

func (v *Value) GetID() string {
	return v.id
}

func (v *Value) GetValue() (reflect.Value, error) {
	copy := v.value
	return reflect.ValueOf(copy), nil
}

func (v *Value) Stop(ctx context.Context) error {
	return nil
}
