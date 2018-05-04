package refactor

import (
	"fmt"
	"reflect"

	"github.com/yarencheng/gospring/refactor/dependency"
)

type applicationContext struct {
	graph      *dependency.Graph
	beans      map[string]BeanI
	singletons map[string]*reflect.Value
}

func NewApplicationContext(beans ...BeanI) (ApplicationContextI, error) {

	ctx := applicationContext{
		graph:      dependency.NewGraph(),
		beans:      make(map[string]BeanI),
		singletons: make(map[string]*reflect.Value),
	}

	for _, bean := range beans {
		if e := ctx.addBean(nil, bean); e != nil {
			return nil, fmt.Errorf("Process bean [%v] failed. Caused by: %v", bean, e)
		}
	}

	return &ctx, nil
}

func (ctx *applicationContext) GetBean(id string) (interface{}, error) {

	bean, present := ctx.beans[id]

	if !present {
		return nil, fmt.Errorf("There is no bean with ID [%v]", id)
	}

	value, e := ctx.getBean(bean)

	if e != nil {
		return nil, e
	}

	return value.Interface(), nil
}

func (ctx *applicationContext) Finalize() error {
	return nil
}

func (ctx *applicationContext) addBean(lastId *string, bean BeanI) error {

	id := bean.GetID()

	switch bean.(type) {
	case ValueBeanI:
	case ReferenceBeanI:

		if lastId != nil {
			if !ctx.graph.AddDependency(*bean.GetID(), *lastId) {
				return fmt.Errorf("Found a cyclic dependency between ID [%v] and ID [%v]", *bean.GetID(), *lastId)
			}
		}

	case StructBeanI:

		if id != nil {
			if ctx.isIdUsed(*id) {
				return fmt.Errorf("ID [%v] of bean [%v] is used", *id, bean)
			}
			ctx.beans[*id] = bean
		}

		if lastId != nil && id != nil {
			if !ctx.graph.AddDependency(*id, *lastId) {
				return fmt.Errorf("Found a cyclic dependency between ID [%v] and ID [%v]", *id, *lastId)
			}
		}

		sBean := bean.(StructBeanI)
		for _, values := range sBean.GetProperties() {
			for _, value := range values {

				nextId := lastId
				if id != nil {
					nextId = id
				}
				if e := ctx.addBean(nextId, value); e != nil {
					return e
				}
			}
		}

	default:
		return fmt.Errorf("Type [%T] of bean [%v] is not support", bean, bean)
	}

	return nil
}

func (ctx *applicationContext) isIdUsed(id string) bool {
	_, present := ctx.beans[id]
	return present
}

func (ctx *applicationContext) getBean(bean BeanI) (*reflect.Value, error) {

	switch bean.(type) {
	case ValueBeanI:
		i := bean.(ValueBeanI).GetValue()
		v := reflect.ValueOf(i)
		return &v, nil
	case ReferenceBeanI:
		id := bean.(ReferenceBeanI).GetID()
		v, e := ctx.getBean(ctx.beans[*id])
		if e != nil {
			return nil, e
		}
		return v, nil
	case StructBeanI:
	default:
		return nil, fmt.Errorf("Type [%T] of bean [%v] is not support", bean, bean)
	}

	sBean := bean.(StructBeanI)
	switch sBean.GetScope() {
	case Singleton:
		return ctx.getSingletonBean(sBean)
	case Prototype:
		return ctx.getPrototypeBean(sBean)
	case Default:
		return ctx.getSingletonBean(sBean)
	default:
		return nil, fmt.Errorf("Scope [%T] of bean [%v] is not support", sBean.GetScope(), sBean)
	}
}

func (ctx *applicationContext) getSingletonBean(bean StructBeanI) (*reflect.Value, error) {

	value, present := ctx.singletons[*bean.GetID()]

	if present {
		return value, nil
	}

	var e error
	value, e = ctx.getPrototypeBean(bean)

	if e != nil {
		return nil, e
	}

	ctx.singletons[*bean.GetID()] = value

	return value, nil
}
func (ctx *applicationContext) getPrototypeBean(bean StructBeanI) (*reflect.Value, error) {

	return nil, fmt.Errorf("TODO")
}
