package contexts

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/yarencheng/gospring/beans"
)

type AbstractApplicatoinContext struct {
	metas     []*beans.BeanMetaData
	metasById map[string]*beans.BeanMetaData
	beansById map[string]interface{}
}

func NewAbstractApplicatoinContext(metas []*beans.BeanMetaData) (*AbstractApplicatoinContext, error) {

	var ctx AbstractApplicatoinContext
	ctx.metas = metas
	ctx.metasById = make(map[string]*beans.BeanMetaData)
	ctx.beansById = make(map[string]interface{})

	for _, meta := range metas {
		ctx.metasById[meta.GetId()] = meta
	}

	return &ctx, nil
}

func (ctx *AbstractApplicatoinContext) GetBean(id string) (interface{}, error) {

	meta := ctx.metasById[id]

	if meta == nil {
		e := fmt.Sprintf("bean with id [%s] is not defined", id)
		return nil, errors.New(e)
	}

	switch meta.GetScope() {
	case beans.Singleton:
		return ctx.getSingletonBean(meta)
	case beans.Prototype:
		return ctx.getPrototypeBean(meta)
	default:
		e := fmt.Errorf("unknown scope [%v]", meta.GetScope())
		return nil, e
	}
}

func (ctx *AbstractApplicatoinContext) getSingletonBean(meta *beans.BeanMetaData) (interface{}, error) {

	if bean := ctx.beansById[meta.GetId()]; bean != nil {
		return bean, nil
	}

	bean := reflect.New(meta.GetStruct()).Interface()

	ctx.beansById[meta.GetId()] = bean

	return bean, nil
}

func (ctx *AbstractApplicatoinContext) getPrototypeBean(meta *beans.BeanMetaData) (interface{}, error) {

	bean := reflect.New(meta.GetStruct()).Interface()

	return bean, nil
}

func (ctx *AbstractApplicatoinContext) Start() error {
	return errors.New("not implement")
}

func (ctx *AbstractApplicatoinContext) Stop() error {
	return errors.New("not implement")
}
