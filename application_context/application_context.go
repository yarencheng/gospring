package application_context

import (
	"container/list"
	"fmt"
	"reflect"

	uuid "github.com/satori/go.uuid"
	"github.com/yarencheng/gospring/bean"
	"github.com/yarencheng/gospring/interfaces"
	"github.com/yarencheng/gospring/v1"
)

type ApplicationContext struct {
	beans       *list.List
	beansByID   map[string]interfaces.BeanI
	beansByUUID map[uuid.UUID]interfaces.BeanI
}

func Default() *ApplicationContext {
	ctx := &ApplicationContext{
		beans:       list.New(),
		beansByID:   make(map[string]interfaces.BeanI),
		beansByUUID: make(map[uuid.UUID]interfaces.BeanI),
	}

	return ctx
}

func New() *ApplicationContext {

	return &ApplicationContext{
		beans:       list.New(),
		beansByID:   make(map[string]interfaces.BeanI),
		beansByUUID: make(map[uuid.UUID]interfaces.BeanI),
	}
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
	switch config.(type) {
	case *v1.Bean:
		cb := config.(*v1.Bean)

		bean, err := bean.NewStructBeanV1(c, cb)
		if err != nil {
			return nil, fmt.Errorf("Create bean failed. err: [%v]", err)
		}

		if _, exist := c.beansByID[bean.GetID()]; exist {
			return nil, fmt.Errorf("ID [%v] allready exists", bean.GetID())
		}

		if _, exist := c.beansByUUID[bean.GetUUID()]; exist {
			return nil, fmt.Errorf("UUID [%v] allready exists", bean.GetUUID())
		}

		c.beans.PushBack(bean)
		c.beansByID[bean.GetID()] = bean
		c.beansByUUID[bean.GetUUID()] = bean

		return bean, nil

	default:
		return nil, fmt.Errorf("[%v] is not a valid config struct", reflect.TypeOf(config).Name())
	}
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
