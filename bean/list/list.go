package list

import (
	"context"
	"fmt"
	"reflect"

	uuid "github.com/satori/go.uuid"
	"github.com/yarencheng/gospring/interfaces"
	"github.com/yarencheng/gospring/v1"
)

type List struct {
	uuid    uuid.UUID
	id      string
	ctx     interfaces.ApplicationContextI
	configs []uuid.UUID
}

func V1ListParser(ctx interfaces.ApplicationContextI, config interface{}) (interfaces.BeanI, error) {

	c, ok := config.(*v1.List)
	if !ok {
		return nil, fmt.Errorf("[%v] can not convert to [%v]", reflect.TypeOf(config), reflect.TypeOf(&v1.List{}))
	}

	l := &List{
		uuid:    uuid.NewV4(),
		id:      c.ID,
		ctx:     ctx,
		configs: make([]uuid.UUID, len(c.Configs)),
	}

	for i, v := range c.Configs {
		b, err := ctx.AddConfig(v)
		if err != nil {
			return nil, fmt.Errorf("Add config [%v] failed. err: [%v]", v, err)
		}
		l.configs[i] = b.GetUUID()
	}

	return l, nil
}

func (l *List) GetUUID() uuid.UUID {
	return l.uuid
}

func (l *List) GetID() string {
	return l.id
}

func (l *List) GetValue() (reflect.Value, error) {

	v := make([]interface{}, len(l.configs))

	for i, uuid := range l.configs {
		e, err := l.ctx.GetByUUID(uuid)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("Get the value failed. err: [%v]", err)
		}
		v[i] = e
	}

	return reflect.ValueOf(v), nil
}

func (l *List) Stop(ctx context.Context) error {
	return nil
}
