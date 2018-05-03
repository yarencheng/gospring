package refactor

import (
	"fmt"

	"github.com/yarencheng/gospring/refactor/dependency"
)

type applicationContext struct {
	graph           *dependency.Graph
	beans           map[string]BeanI
	singletonValues map[string]interface{}
}

func NewApplicationContext(beans ...BeanI) (ApplicationContextI, error) {

	ctx := applicationContext{
		graph:           dependency.NewGraph(),
		beans:           make(map[string]BeanI),
		singletonValues: make(map[string]interface{}),
	}

	for _, bean := range beans {
		if e := ctx.addBean(nil, bean); e != nil {
			return nil, fmt.Errorf("Process bean [%v] failed. Caused by: %v", bean, e)
		}
	}

	return &ctx, nil
}

func (ctx *applicationContext) GetBean() (interface{}, error) {
	return nil, nil
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
