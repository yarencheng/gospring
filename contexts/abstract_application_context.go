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

	for _, meta := range metas {
		if old, present := ctx.metasById[meta.GetId()]; present {
			e := fmt.Errorf("id [%s] is already used by the other bean [%v]", meta.GetId(), old)
			return nil, e
		}
		ctx.metasById[meta.GetId()] = meta
	}

	// check all reference inside beans exist
	for _, meta := range metas {
		for _, property := range meta.GetProperties() {
			if _, present := ctx.metasById[property.GetReference()]; !present {
				e := fmt.Errorf("reference [%s] of property [%s] is dose not exist", property.GetReference(), meta.GetId())
				return nil, e
			}
		}
	}

	// check dependecy info
	if e := ctx.checkDependencyLoop(); e != nil {
		eout := fmt.Errorf("Detect a circle dependency. Cuased by: %v", e)
		return nil, eout
	}

	return &ctx, nil
}

type node struct {
	id     string
	childs map[string]*node
	isWalk bool
}

func recursive(walked []*node, cur *node) error {
	if len(cur.childs) == 0 {
		return nil
	}
	for _, ch := range cur.childs {
		if ch.isWalk {
			var buffer bytes.Buffer
			for _, w := range walked {
				buffer.WriteString("[")
				buffer.WriteString(w.id)
				buffer.WriteString("]>")
			}
			e := fmt.Errorf("detect a loop %s[%s]", buffer.String(), ch.id)
			return e
		}
	}
	for _, ch := range cur.childs {
		if e := recursive(append(walked, ch), ch); e != nil {
			return e
		}
	}

	return nil
}

func (ctx *AbstractApplicatoinContext) checkDependencyLoop() error {

	nodes := make(map[string]*node)

	for _, meta := range ctx.metas {
		nodes[meta.GetId()] = &node{
			id:     meta.GetId(),
			childs: make(map[string]*node),
			isWalk: false,
		}
	}

	for _, meta := range ctx.metas {
		for _, property := range meta.GetProperties() {
			nodes[meta.GetId()].childs[property.GetReference()] = nodes[property.GetReference()]
		}
	}

	for _, v := range nodes {
		if e := recursive(make([]*node, 10), v); e != nil {
			return e
		}
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
