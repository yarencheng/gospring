package bean

import (
	"context"
	"fmt"
	"reflect"

	uuid "github.com/satori/go.uuid"
	"github.com/yarencheng/gospring/interfaces"
	"github.com/yarencheng/gospring/v1"
)

type ChannelBean struct {
	uuid    uuid.UUID
	id      string
	size    int
	tvpe    reflect.Type
	channel reflect.Value
}

func V1ChannelParser(ctx interfaces.ApplicationContextI, config interface{}) (interfaces.BeanI, error) {
	c, ok := config.(*v1.Channel)
	if !ok {
		return nil, fmt.Errorf("[%v] can not convert to [%v]", reflect.TypeOf(config), reflect.TypeOf(&v1.Channel{}))
	}

	b := &ChannelBean{
		uuid: uuid.NewV4(),
		id:   c.ID,
		tvpe: c.Type,
	}

	return b, nil
}

func (b *ChannelBean) GetUUID() uuid.UUID {
	return b.uuid
}

func (b *ChannelBean) GetID() string {
	return b.id
}

func (b *ChannelBean) GetValue() (reflect.Value, error) {
	if b.channel.IsValid() {
		return b.channel, nil
	}
	tvpe := reflect.ChanOf(reflect.BothDir, b.tvpe)
	b.channel = reflect.MakeChan(tvpe, b.size)
	return b.channel, nil
}

func (b *ChannelBean) Stop(ctx context.Context) error {
	if b.channel.IsValid() {
		b.channel.Close()
	}
	return nil
}
