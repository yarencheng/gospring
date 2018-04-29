package gospring

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
		for _, pb := range b.pros[name] {
			ctx.getBeanByIDRecursive(pb.Interface().(*bean))
		}
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
		return ctx.getSingleTonBean(bean)
	case scopeSingleton:
		return ctx.getSingleTonBean(bean)
	case scopePrototype:
		return ctx.getPrototypeBean(bean)
	default:
		return reflect.Value{}, fmt.Errorf("unsupport scope [%v]", bean.scope)
	}
}

func (ctx *applicationContext) getSingleTonBean(bean *bean) (reflect.Value, error) {

	if bean.id != nil {
		if v, present := ctx.singletonById[*bean.id]; present {
			return v, nil
		}
	}

	v, e := ctx.getPrototypeBean(bean)
	if e != nil {
		return reflect.Value{}, fmt.Errorf("Can't create bean [%v]. Cuased by: %v", *bean.id, e)
	}

	if bean.id != nil {
		ctx.singletonById[*bean.id] = v
	}

	return v, nil
}

func (ctx *applicationContext) getPrototypeBean(b *bean) (reflect.Value, error) {

	i, e := b.new()
	if e != nil {
		return reflect.Value{}, fmt.Errorf("Create bean [%v] failed. Caused by: ", e)
	}

	v := reflect.ValueOf(i)

	for name, type_ := range b.prosType {

		field := v.Elem().FieldByName(name)
		var values []reflect.Value

		switch type_ {
		case propertyTypeBean:
			values = make([]reflect.Value, len(b.pros[name]))
			for i, _ := range b.pros[name] {
				pb := b.pros[name][i].Interface().(*bean)
				vpb, vpbe := ctx.getBean(pb)
				if vpbe != nil {
					return reflect.Value{}, fmt.Errorf("Get bean [%v] failed. Caused by: %v", *pb, vpbe)
				}
				values[i] = vpb
			}

		case propertyTypeRef:
			values = make([]reflect.Value, len(b.pros[name]))
			for i, _ := range b.pros[name] {
				id := b.pros[name][i].Interface().(*refBean).ref
				pb := ctx.beanById[id]
				vpb, vpbe := ctx.getBean(pb)
				if vpbe != nil {
					return reflect.Value{}, fmt.Errorf("Get bean [%v] failed. Caused by: %v", *pb, vpbe)
				}
				values[i] = vpb
			}

		case propertyTypeValue:
			values = make([]reflect.Value, len(b.pros[name]))
			for i, _ := range b.pros[name] {
				vt := reflect.New(b.pros[name][i].Type())
				vt.Elem().Set(b.pros[name][i])
				values[i] = vt
			}

		default:
			return reflect.Value{}, fmt.Errorf("Type of property [%v] is unknown", type_)
		}

		switch field.Type().Kind() {
		case reflect.Slice:
			slice := reflect.MakeSlice(field.Type(), len(values), len(values))

			switch field.Type().Elem().Kind() {
			case reflect.Ptr:
				for i, value := range values {
					slice.Index(i).Set(value)
				}
			default:
				for i, value := range values {
					slice.Index(i).Set(value.Elem())
				}
			}

			field.Set(slice)
		case reflect.Ptr:
			field.Set(values[0])
		default:
			field.Set(values[0].Elem())
		}
	}

	if e := ctx.execBeanInit(v, b); e != nil {
		return reflect.Value{}, fmt.Errorf("Initialize bean [%v] failed. Caused by: %v", *b, e)
	}

	return v, nil
}

func (ctx *applicationContext) execBeanInit(value reflect.Value, bena *bean) error {
	if bena.initFn != nil {
		rv := bena.initFn.Call([]reflect.Value{value})
		switch len(rv) {
		case 0:
			return nil
		case 1:
			if e, ok := rv[0].Interface().(error); ok {
				return fmt.Errorf(
					"init bean [%v] failed. Caused by: %v",
					*bena,
					e,
				)
			} else {
				return fmt.Errorf(
					"init function of bean [%v] returns 1 unexpected value",
					*bena,
				)
			}
		default:
			return fmt.Errorf(
				"init function of bean [%v] returns %d unexpected values",
				*bena,
				len(rv),
			)
		}
	}

	// try to find A.Init(a* A) function in struct
	// TODO

	return nil
}
