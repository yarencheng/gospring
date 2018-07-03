package application_context

import (
	"bytes"
	"container/list"
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/yarencheng/gospring/bean"
	"github.com/yarencheng/gospring/bean/value"
	"github.com/yarencheng/gospring/v1"

	uuid "github.com/satori/go.uuid"
	"github.com/yarencheng/gospring/interfaces"
)

type ApplicationContext struct {
	beans        *list.List
	beansByID    map[string]interfaces.BeanI
	beansByUUID  map[uuid.UUID]interfaces.BeanI
	configParser map[reflect.Type]interfaces.ConfigParser
}

func New() *ApplicationContext {

	return &ApplicationContext{
		beans:        list.New(),
		beansByID:    make(map[string]interfaces.BeanI),
		beansByUUID:  make(map[uuid.UUID]interfaces.BeanI),
		configParser: make(map[reflect.Type]interfaces.ConfigParser),
	}
}

func Default() *ApplicationContext {
	ctx := New()

	ctx.UseConfigParser(reflect.TypeOf(&v1.Bean{}), bean.V1BeanParser)
	ctx.UseConfigParser(reflect.TypeOf(&v1.Channel{}), bean.V1ChannelParser)
	ctx.UseConfigParser(reflect.TypeOf(&v1.Ref{}), bean.V1RefParser)
	ctx.UseConfigParser(reflect.TypeOf(""), bean.V1RefParser)
	ctx.UseConfigParser(reflect.TypeOf(&v1.Broadcast{}), bean.V1BroadcastParser)
	ctx.UseConfigParser(reflect.TypeOf(&v1.Value{}), value.V1ValueParser)

	return ctx
}

func (c *ApplicationContext) UseConfigParser(configType reflect.Type, parser interfaces.ConfigParser) error {
	c.configParser[configType] = parser
	return nil
}

func (c *ApplicationContext) AddConfigs(configs ...interface{}) error {
	for _, config := range configs {
		if _, err := c.AddConfig(config); err != nil {
			return fmt.Errorf("Add config [%#v] failed. err: [%v]", config, err)
		}
	}
	return nil
}

func (c *ApplicationContext) AddConfig(config interface{}) (interfaces.BeanI, error) {
	tvpe := reflect.TypeOf(config)

	parser, exist := c.configParser[tvpe]
	if !exist {
		return nil, fmt.Errorf("Can't find a parser for [%v]", tvpe)
	}

	bean, err := parser(c, config)
	if err != nil {
		return nil, fmt.Errorf("Create bean failed. err: [%v]", err)
	}

	if bean.GetID() != "" {
		if _, exist := c.beansByID[bean.GetID()]; exist {
			return nil, fmt.Errorf("ID [%v] allready exists", bean.GetID())
		} else {
			c.beansByID[bean.GetID()] = bean
		}
	}

	if _, exist := c.beansByUUID[bean.GetUUID()]; exist {
		return nil, fmt.Errorf("UUID [%v] allready exists", bean.GetUUID())
	} else {
		c.beansByUUID[bean.GetUUID()] = bean
	}

	c.beans.PushBack(bean)

	return bean, nil
}

func (c *ApplicationContext) GetByID(id string) (interface{}, error) {

	if len(id) == 0 {
		return nil, fmt.Errorf("id is empty")
	}

	bean, ok := c.beansByID[id]

	if !ok {
		return nil, fmt.Errorf("ID [%v] dose not exist", id)
	}

	v, err := bean.GetValue()
	if err != nil {
		return nil, err
	}

	return v.Interface(), nil
}

func (c *ApplicationContext) GetByUUID(id uuid.UUID) (interface{}, error) {
	if id == uuid.Nil {
		return nil, fmt.Errorf("id is empty")
	}

	bean, ok := c.beansByUUID[id]

	if !ok {
		return nil, fmt.Errorf("ID [%v] dose not exist", id)
	}

	v, err := bean.GetValue()
	if err != nil {
		return nil, err
	}

	return v.Interface(), nil
}

func (c *ApplicationContext) GetBeanByID(id string) (interfaces.BeanI, bool) {

	if len(id) == 0 {
		return nil, false
	}

	bean, ok := c.beansByID[id]

	if !ok {
		return nil, false
	}

	return bean, true
}

func (c *ApplicationContext) GetBeanByUUID(id uuid.UUID) (interfaces.BeanI, bool) {
	if id == uuid.Nil {
		return nil, false
	}

	bean, ok := c.beansByUUID[id]

	if !ok {
		return nil, false
	}

	return bean, true
}

func (c *ApplicationContext) Stop(ctx context.Context) error {

	errsLock := sync.Mutex{}
	errs := list.New()

	wg := sync.WaitGroup{}
	wg.Add(c.beans.Len())

	for e := c.beans.Front(); e != nil; e = e.Next() {
		v, ok := e.Value.(interfaces.BeanI)
		if !ok {
			return fmt.Errorf("Can't convert [%v] to type interfaces.BeanI", reflect.TypeOf(e.Value))
		}
		go func() {
			err := v.Stop(ctx)
			if err != nil {
				errsLock.Lock()
				defer errsLock.Unlock()
				errs.PushBack(err)
			}
			wg.Done()
		}()
	}

	wait := make(chan int)
	go func() {
		wg.Wait()
		wait <- 1
	}()

	select {
	case <-ctx.Done():
	case <-wait:
	}

	errsLock.Lock()
	defer errsLock.Unlock()

	if err := ctx.Done(); err != nil {
		errs.PushBack(err)
	}

	if errs.Len() > 0 {
		var buffer bytes.Buffer
		buffer.WriteString("[")
		for err := errs.Front(); err != nil; err = err.Next() {
			if err != errs.Front() {
				buffer.WriteString(", ")
			}
			buffer.WriteString("[")
			buffer.WriteString(err.Value.(error).Error())
			buffer.WriteString("]")
		}
		buffer.WriteString("]")
		return errors.New(buffer.String())
	}

	return nil
}
