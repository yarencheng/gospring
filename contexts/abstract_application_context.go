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

func NewAbstractApplicatoinContext(metas []*beans.BeanMetaData) *AbstractApplicatoinContext {

	var ctx AbstractApplicatoinContext
	ctx.metas = metas
	ctx.metasById = make(map[string]*beans.BeanMetaData)
	ctx.beansById = make(map[string]interface{})

	for _, meta := range metas {
		ctx.metasById[meta.GetId()] = meta
	}

	return &ctx
}

func (ctx *AbstractApplicatoinContext) GetBean(id string) (interface{}, error) {

	if bean := ctx.beansById[id]; bean != nil {
		return bean, nil
	}

	meta := ctx.metasById[id]

	if meta == nil {
		e := fmt.Sprintf("bean with id [%s] is not defined", id)
		return nil, errors.New(e)
	}

	bean := reflect.New(meta.GetStruct()).Interface()

	ctx.beansById[id] = bean

	return bean, nil
}

func (ctx *AbstractApplicatoinContext) Start() error {
	return errors.New("not implement")
}

func (ctx *AbstractApplicatoinContext) Stop() error {
	return errors.New("not implement")
}
