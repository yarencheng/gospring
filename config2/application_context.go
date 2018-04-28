package config2

import (
	"fmt"
	"reflect"
)

type applicationContext struct {
	beans         *beans
	beanById      map[string]*bean
	singletonById map[string]reflect.Value
}

func ApplicationContext(bs *beans) (*applicationContext, error) {

	// TODO: validate beans

	ctx := &applicationContext{
		beans:         bs,
		beanById:      make(map[string]*bean),
		singletonById: make(map[string]reflect.Value),
	}

	for _, b := range bs.beans {
		ctx.getBeanByIDRecursive(b)
	}

	return ctx, nil
}

func (ctx *applicationContext) getBeanByIDRecursive(b *bean) {

	if b.id != nil {
		ctx.beanById[*b.id] = b
	}

	for name, type_ := range b.prosType {
		if type_ != propertyTypeBean {
			continue
		}
		ctx.getBeanByIDRecursive(b.pros[name].Interface().(*bean))
	}
}

func (ctx *applicationContext) GetBean(id string) (interface{}, error) {

	b, present := ctx.beanById[id]
	if !present {
		return nil, fmt.Errorf("There is no bean with ID [%v]", id)
	}

	v, e := ctx.getBean(b)
	if e != nil {
		return nil, e
	}
	return v.Interface(), nil
}

func (ctx *applicationContext) getBean(bean *bean) (reflect.Value, error) {

	switch bean.scope {
	case scopeDefault:
		return ctx.GetSingleTonBean(bean)
	case scopeSingleton:
		return ctx.GetSingleTonBean(bean)
	case scopePrototype:
		return ctx.GetPrototypeBean(bean)
	default:
		return reflect.Value{}, fmt.Errorf("unsupport scope [%v]", bean.scope)
	}
}

func (ctx *applicationContext) GetSingleTonBean(bean *bean) (reflect.Value, error) {

	if bean.id != nil {
		if v, present := ctx.singletonById[*bean.id]; present {
			return v, nil
		}
	}

	v, e := ctx.GetPrototypeBean(bean)
	if e != nil {
		return reflect.Value{}, fmt.Errorf("Can't create bean [%v]. Cuased by: %v", *bean.id, e)
	}

	if bean.id != nil {
		ctx.singletonById[*bean.id] = v
	}

	return v, nil
}

func (ctx *applicationContext) GetPrototypeBean(b *bean) (reflect.Value, error) {

	i, e := b.new()
	if e != nil {
		return reflect.Value{}, fmt.Errorf("Create bean [%v] failed. Caused by: ", e)
	}

	v := reflect.ValueOf(i)

	for name, type_ := range b.prosType {

		field := v.Elem().FieldByName(name)
		var value reflect.Value
		var valueError error

		switch type_ {
		case propertyTypeBean:
			pb := b.pros[name].Interface().(*bean)
			value, valueError = ctx.getBean(pb)
			if valueError != nil {
				return reflect.Value{}, fmt.Errorf("Get bean [%v] failed. Caused by: %v", *pb, valueError)
			}

		case propertyTypeRef:
			id := b.pros[name].Interface().(string)
			pb := ctx.beanById[id]
			value, valueError = ctx.getBean(pb)
			if valueError != nil {
				return reflect.Value{}, fmt.Errorf("Get bean with ID [%v] failed. Caused by: %v", id, valueError)
			}

		case propertyTypeValue:
			value = reflect.New(b.pros[name].Type())
			value.Elem().Set(b.pros[name])

		default:
			return reflect.Value{}, fmt.Errorf("Type of property [%v] is unknown", type_)
		}

		switch field.Type().Kind() {
		case reflect.Ptr:
			field.Set(value)
		default:
			field.Set(value.Elem())
		}
	}

	return v, nil
}
