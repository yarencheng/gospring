package config

import (
	"fmt"
	"reflect"
	"strconv"
)

type config struct {
	beans []*bean
}

type bean struct {
	id    string
	type_ reflect.Type
	scope scope
	pros  []*property
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
		pros:  make([]*property, 0),
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

func (b *bean) With(p ...*property) *bean {
	b.pros = append(b.pros, p...)
	return b
}

type propertyType string

const (
	pValue     propertyType = "value"
	pReference propertyType = "reference"
	pBean      propertyType = "bean"
)

type property struct {
	name  string
	type_ propertyType
	ref   string
	bean  *bean
	value string
}

func Ref(name, ref string) *property {
	return &property{
		name:  name,
		type_: pReference,
		ref:   ref,
	}
}

func PropertyBean(name string, b *bean) *property {
	return &property{
		name:  name,
		type_: pBean,
		bean:  b,
	}
}

func Value(name string, v string) *property {
	return &property{
		name:  name,
		type_: pValue,
		value: v,
	}
}

type applicationContext struct {
	config        *config
	beanById      map[string]*bean
	singletonById map[string]reflect.Value
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
		singletonById: make(map[string]reflect.Value),
	}, nil
}

func (ctx *applicationContext) GetBean(id string) (interface{}, error) {
	v, e := ctx.getBean(id)
	if e != nil {
		return nil, e
	}
	return v.Interface(), nil
}

func (ctx *applicationContext) getBean(id string) (reflect.Value, error) {
	bean, present := ctx.beanById[id]
	if !present {
		return reflect.Value{}, fmt.Errorf("no bean with id [%v]", id)
	}

	switch bean.scope {
	case singleton:
		return ctx.GetSingleTonBean(bean)
	case prototype:
		return ctx.GetPrototypeBean(bean)
	default:
		return reflect.Value{}, fmt.Errorf("unsupport scope [%v]", bean.scope)
	}
}

func (ctx *applicationContext) GetSingleTonBean(bean *bean) (reflect.Value, error) {

	if v, present := ctx.singletonById[bean.id]; present {
		return v, nil
	}

	v, e := ctx.GetPrototypeBean(bean)
	if e != nil {
		return reflect.Value{}, fmt.Errorf("Can't create bean [%v]. Cuased by: %v", bean.id, e)
	}

	ctx.singletonById[bean.id] = v

	return v, nil
}

func (ctx *applicationContext) GetPrototypeBean(bean *bean) (reflect.Value, error) {

	v := reflect.New(bean.type_)

	for _, p := range bean.pros {

		if _, present := bean.type_.FieldByName(p.name); !present {
			return reflect.Value{}, fmt.Errorf("there is no member named [%v] in struct [%v]", p.name, bean.type_.Name())
		}

		field := v.Elem().FieldByName(p.name)

		if !field.CanSet() {
			return reflect.Value{}, fmt.Errorf("member named [%v] in struct [%v] is not setable", p.name, bean.type_.Name())
		}

		var e error
		switch p.type_ {
		case pValue:
			e = ctx.setNativeField(field, p.value)
		case pBean:
			e = fmt.Errorf("TODO")
		case pReference:
			e = fmt.Errorf("TODO")
		default:
			return reflect.Value{}, fmt.Errorf("type of member named [%v] in struct [%v] is unknown", p.name, bean.type_.Name())
		}

		if e != nil {
			return reflect.Value{}, e
		}
	}

	return v, nil
}

func (ctx *applicationContext) setNativeField(field reflect.Value, value string) error {

	switch field.Type().Kind() {
	case reflect.String:
		field.Set(reflect.ValueOf(value))
	case reflect.Int:
		i, e := strconv.ParseInt(value, 10, 32)
		if e != nil {
			return fmt.Errorf("[%v] can't convert to int. Caused by: %v", value, e)
		}
		field.Set(reflect.ValueOf(int(i)))
	default:
		return fmt.Errorf("Unsopport type %v", field.Type())
	}

	return nil
}
