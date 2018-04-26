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
	beansById map[string]reflect.Value
}

func NewAbstractApplicatoinContext(metas []*beans.BeanMetaData) (*AbstractApplicatoinContext, error) {

	var ctx AbstractApplicatoinContext
	ctx.metas = metas
	ctx.metasById = make(map[string]*beans.BeanMetaData)
	ctx.beansById = make(map[string]reflect.Value)

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

	// only one can exist in a property between refrence or value
	for _, meta := range ctx.metas {
		for _, property := range meta.GetProperties() {
			if (len(property.GetValue()) == 0) != (len(property.GetReference()) == 0) {
				continue
			}
			e := fmt.Errorf("only one can exist in a property between refrence or value [%v]", property.GetName())
			return e
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
			if len(property.GetReference()) == 0 {
				continue
			}
			nodes[meta.GetId()].childs[property.GetReference()] = nodes[property.GetReference()]
		}
	}

	for _, v := range nodes {
		m := make(map[string]*node)
		if e := walkAndCheckLoop(m, v); e != nil {
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
	if _, present := walked[cur.id]; present {
		var buffer bytes.Buffer
		for _, w := range walked {
			buffer.WriteString("[")
			buffer.WriteString(w.id)
			buffer.WriteString("]>")
		}
		e := fmt.Errorf("detect a loop %s[%s]", buffer.String(), cur.id)
		return e
	}

	walked[cur.id] = cur

	for _, v := range cur.childs {
		if e := walkAndCheckLoop(walked, v); e != nil {
			return e
		}
	}

	delete(walked, cur.id)

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
		b, e := ctx.getSingletonBean(meta)
		if e != nil {
			return nil, e
		}
		return b.Interface(), nil
	case beans.Prototype:
		b, e := ctx.getPrototypeBean(meta)
		if e != nil {
			return nil, e
		}
		return b.Interface(), nil
	default:
		e := fmt.Errorf("unknown scope [%v]", meta.GetScope())
		return nil, e
	}
}

func (ctx *AbstractApplicatoinContext) getSingletonBean(meta *beans.BeanMetaData) (reflect.Value, error) {

	if bean, present := ctx.beansById[meta.GetId()]; present {
		return bean, nil
	}

	bean, e := ctx.getPrototypeBean(meta)

	if e != nil {
		return reflect.Value{}, e
	}

	ctx.beansById[meta.GetId()] = bean

	return bean, nil
}

func (ctx *AbstractApplicatoinContext) getPrototypeBean(meta *beans.BeanMetaData) (reflect.Value, error) {

	bean := reflect.New(meta.GetStruct())

	for _, p := range meta.GetProperties() {

		if _, present := meta.GetStruct().FieldByName(p.GetName()); !present {
			e := fmt.Errorf("There is no field named [%v] in bean [%v]", p.GetName(), meta.GetId())
			return reflect.Value{}, e
		}

		field := bean.Elem().FieldByName(p.GetName())

		if len(p.GetReference()) > 0 {
			// TODO: inject bean
			continue
		}

		switch field.Type().Kind() {
		case reflect.String:
			field.Set(reflect.ValueOf(p.GetValue()))
		default:
			e := fmt.Errorf("Unsopport type %v", field.Type())
			return reflect.Value{}, e
		}
	}

	return bean, nil
}

func (ctx *AbstractApplicatoinContext) Start() error {
	return errors.New("not implement")
}

func (ctx *AbstractApplicatoinContext) Stop() error {
	return errors.New("not implement")
}
