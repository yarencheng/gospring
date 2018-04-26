package contexts

import (
	"bytes"
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

	// create map with (key,value)=(id,beanmeta)
	for _, meta := range metas {
		if old, present := ctx.metasById[meta.GetId()]; present {
			e := fmt.Errorf("id [%s] is already used by the other bean [%v]", meta.GetId(), old)
			return nil, e
		}
		ctx.metasById[meta.GetId()] = meta
	}

	if e := ctx.checkProperty(); e != nil {
		return nil, e
	}

	return &ctx, nil
}

func (ctx *AbstractApplicatoinContext) checkProperty() error {
	// chack each name in a bean is unique
	for _, meta := range ctx.metas {
		names := make(map[string]string)
		for _, property := range meta.GetProperties() {
			if _, present := names[property.GetName()]; present {
				e := fmt.Errorf("Found duplicated name [%v] in bean [%v]", property.GetName(), meta.GetId())
				return e
			}
		}
	}

	// check that all references exist
	for _, meta := range ctx.metas {
		for _, property := range meta.GetProperties() {
			if len(property.GetReference()) == 0 {
				continue
			}
			if _, present := ctx.metasById[property.GetReference()]; !present {
				e := fmt.Errorf("reference [%s] of property [%s] is dose not exist", property.GetReference(), meta.GetId())
				return e
			}
		}
	}

	// check loop

	nodes := make(map[string]*node)

	for _, meta := range ctx.metas {
		nodes[meta.GetId()] = &node{
			id:     meta.GetId(),
			childs: make(map[string]*node),
		}
	}

	for _, meta := range ctx.metas {
		for _, property := range meta.GetProperties() {
			nodes[meta.GetId()].childs[property.GetReference()] = nodes[property.GetReference()]
		}
	}

	for _, v := range nodes {
		if e := walkAndCheckLoop(make(map[string]*node), v); e != nil {
			return e
		}
	}

	return nil
}

type node struct {
	id     string
	childs map[string]*node
}

func walkAndCheckLoop(walked map[string]*node, cur *node) error {
	if len(cur.childs) == 0 {
		return nil
	}
	for k, v := range cur.childs {
		if _, present := walked[k]; present {
			var buffer bytes.Buffer
			for _, w := range walked {
				buffer.WriteString("[")
				buffer.WriteString(w.id)
				buffer.WriteString("]>")
			}
			e := fmt.Errorf("detect a loop %s[%s]", buffer.String(), v.id)
			return e
		}

		walked[k] = v
		if e := walkAndCheckLoop(walked, v); e != nil {
			return e
		}
		delete(walked, k)
	}

	return nil
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
