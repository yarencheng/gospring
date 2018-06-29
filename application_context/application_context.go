package application_context

import (
	"container/list"
	"fmt"
	"reflect"

	"github.com/yarencheng/gospring/bean"
	"github.com/yarencheng/gospring/v1"
)

type ApplicationContext struct {
	beans *list.List
}

func New(configs ...interface{}) (*ApplicationContext, error) {

	beans, err := configsToBeans(configs...)
	if err != nil {
		return nil, err
	}

	return &ApplicationContext{
		beans: beans,
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

	for e := c.beans.Front(); e != nil; e = e.Next() {

		if b, ok := e.Value.(bean.BeanI); !ok {
			return nil, fmt.Errorf("[%v] is not a bean.BeanI", reflect.TypeOf(e.Value))
		} else if id == b.GetID() {
			v, err := b.GetValue()
			if err != nil {
				return nil, err
			}
			return v.Interface(), nil
		}
	}

	return nil, fmt.Errorf("ID [%v] dose not exist", id)
}
