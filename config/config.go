package config

import (
	"fmt"
	"reflect"
)

type config struct {
	beans []*bean
}

type bean struct {
	id    string
	type_ reflect.Type
	scope scope
}

type property struct {
	name  string
	isRef bool
	ref   string
	bean  *bean
}

type scope string

const (
	singleton scope = "singleton"
	prototype scope = "prototype"
)

func Config(beans ...*bean) *config {
	return &config{
		beans: beans,
	}
}

func (c *config) AddBean(beans ...*bean) *config {
	c.beans = append(c.beans, beans...)
	return c
}

func (c *config) Validate() bool {
	// TODO
	return true
}

func Bean(id string, type_ reflect.Type) *bean {
	return &bean{
		id:    id,
		type_: type_,
		scope: singleton,
	}
}

func (b *bean) Singleton() *bean {
	b.scope = singleton
	return b
}

func (b *bean) Prototype() *bean {
	b.scope = prototype
	return b
}

func Ref(name, ref string) *property {
	return &property{
		name:  name,
		isRef: true,
		ref:   ref,
	}
}

func PropertyBean(name string, b *bean) *property {
	return &property{
		name:  name,
		isRef: false,
		bean:  b,
	}
}

type applicationContext struct {
	config   *config
	beanById map[string]*bean
}

func ApplicationContext(config *config) (*applicationContext, error) {

	if !config.Validate() {
		return nil, fmt.Errorf("Configuration is not valid")
	}

	return &applicationContext{
		config: config,
		beanById: func() map[string]*bean {

			m := make(map[string]*bean)
			for _, bean := range config.beans {
				m[bean.id] = bean
			}

			return m
		}(),
	}, nil
}

func (ctx *applicationContext) GetBean(id string) (interface{}, error) {
	bean, present := ctx.beanById[id]
	if !present {
		return nil, fmt.Errorf("no bean with id [%v]", id)
	}
	v := reflect.New(bean.type_)
	return v.Interface(), nil
}
