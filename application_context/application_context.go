package application_context

import (
	"container/list"
	"fmt"
	"reflect"

	uuid "github.com/satori/go.uuid"
	"github.com/yarencheng/gospring/bean"
	"github.com/yarencheng/gospring/v1"
)

type ApplicationContext struct {
	beans       *list.List
	beansByID   map[string]bean.BeanI
	beansByUUID map[uuid.UUID]bean.BeanI
}

func New(configs ...interface{}) (*ApplicationContext, error) {

	beans, err := configsToBeans(configs...)
	if err != nil {
		return nil, err
	}

	return &ApplicationContext{
		beans:       beans,
		beansByID:   createIDMap(beans),
		beansByUUID: createUUIDMap(beans),
	}, nil
}

func configsToBeans(configs ...interface{}) (*list.List, error) {
	beans := list.New()
	for _, config := range configs {
		switch config.(type) {
		case *v1.Bean:
			cb := config.(*v1.Bean)
			bean, err := bean.NewStructBeanV1(cb)
			beans.PushBack(bean)
			if err != nil {
				return nil, err
			}
			for _, p := range config.(*v1.Bean).Properties {
				pBeans, err := configsToBeans(p.Value)
				if err != nil {
					return nil, err
				}
				beans.PushBackList(pBeans)
			}
		default:
			return nil, fmt.Errorf("[%v] is not a valid config struct", reflect.TypeOf(config).Name())
		}
	}
	return beans, nil
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

func createIDMap(list *list.List) map[string]bean.BeanI {
	m := make(map[string]bean.BeanI)
	for e := list.Front(); e != nil; e = e.Next() {
		b := e.Value.(bean.BeanI)
		if b.GetID() == "" {
			continue
		}
		m[b.GetID()] = b
	}
	return m
}

func createUUIDMap(list *list.List) map[uuid.UUID]bean.BeanI {
	m := make(map[uuid.UUID]bean.BeanI)
	for e := list.Front(); e != nil; e = e.Next() {
		b := e.Value.(bean.BeanI)
		if b.GetUUID() == uuid.Nil {
			panic("UUID is empty")
		}
		m[b.GetUUID()] = b
	}
	return m
}
