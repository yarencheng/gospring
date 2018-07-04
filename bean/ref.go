package bean

import (
	"context"
	"fmt"
	"reflect"

	uuid "github.com/satori/go.uuid"
	"github.com/yarencheng/gospring/interfaces"
	"github.com/yarencheng/gospring/v1"
)

type Reference struct {
	uuid     uuid.UUID
	targetID string
	ctx      interfaces.ApplicationContextI
}

func V1RefParser(ctx interfaces.ApplicationContextI, config interface{}) (interfaces.BeanI, error) {

	switch config.(type) {
	case (*v1.Ref):
		return &Reference{
			uuid:     uuid.NewV4(),
			targetID: config.(*v1.Ref).ID,
			ctx:      ctx,
		}, nil
	case string:
		return &Reference{
			uuid:     uuid.NewV4(),
			targetID: config.(string),
			ctx:      ctx,
		}, nil
	default:
		return nil, fmt.Errorf("[%v] is not a *v1.Broadcast nor a string", reflect.TypeOf(config))
	}
}

func (r *Reference) GetUUID() uuid.UUID {
	return r.uuid
}

func (r *Reference) GetID() string {
	return ""
}

func (r *Reference) GetValue() (reflect.Value, error) {

	v, err := r.ctx.GetByID(r.targetID)
	if err != nil {
		return reflect.Value{}, fmt.Errorf("Get the value from the bean failed. err: [%v]", err)
	}

	return reflect.ValueOf(v), nil
}

func (r *Reference) Stop(ctx context.Context) error {
	return nil
}
