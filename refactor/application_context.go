package refactor

import (
	"fmt"
	"reflect"

	"github.com/yarencheng/gospring/refactor/dependency"
)

type applicationContext struct {
	graph         *dependency.Graph
	beanById      map[string]BeanI
	parentByChild map[BeanI]BeanI
	singletons    map[string]*reflect.Value
}

func NewApplicationContext(beans ...BeanI) (ApplicationContextI, error) {

	ctx := applicationContext{
		graph:         dependency.NewGraph(),
		beanById:      make(map[string]BeanI),
		parentByChild: make(map[BeanI]BeanI),
		singletons:    make(map[string]*reflect.Value),
	}

	for _, bean := range beans {
		if e := ctx.addBean(bean); e != nil {
			return nil, fmt.Errorf("Can't add bean [%v]. Cuased by: %v", bean, e)
		}
	}

	return &ctx, nil
}

func (ctx *applicationContext) GetBean(id string) (interface{}, error) {

	bean, present := ctx.beanById[id]

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

func (ctx *applicationContext) addBean(bean BeanI) error {

	switch bean.(type) {
	case ValueBeanI:
		if e := ctx.addValueBean(bean.(ValueBeanI)); e != nil {
			return fmt.Errorf("Can't add bean [%v]. Cuased by: %v", bean, e)
		}
	case ReferenceBeanI:
		if e := ctx.addReferenceBean(bean.(ReferenceBeanI)); e != nil {
			return fmt.Errorf("Can't add bean [%v]. Cuased by: %v", bean, e)
		}
	case StructBeanI:
		if e := ctx.addStructBean(bean.(StructBeanI)); e != nil {
			return fmt.Errorf("Can't add bean [%v]. Cuased by: %v", bean, e)
		}
	default:
		return fmt.Errorf("bean type [%T] is unknown", bean)
	}
	return nil
}

func (ctx *applicationContext) addValueBean(bean ValueBeanI) error {

	id := bean.GetID()

	if id != nil {
		if _, present := ctx.beanById[*id]; present {
			return fmt.Errorf("ID [%v] already exist", *id)
		}
		ctx.beanById[*id] = bean
	}

	return nil
}

func (ctx *applicationContext) addReferenceBean(bean ReferenceBeanI) error {
	return nil
}

func (ctx *applicationContext) addStructBean(bean StructBeanI) error {

	if id := bean.GetID(); id != nil {
		if _, present := ctx.beanById[*id]; present {
			return fmt.Errorf("ID [%v] already exist", *id)
		}
		ctx.beanById[*id] = bean
	}

	if fn, argv := bean.GetFactory(); fn != nil {
		tvpe := reflect.TypeOf(fn)

		if tvpe.Kind() != reflect.Func {
			return fmt.Errorf("Factory of bean [%v] is not a function but a [%v]", bean, tvpe.Kind())
		}

		if len(argv) != tvpe.NumIn() {
			return fmt.Errorf("Factory of bean [%v] need [%v] instead of [%v] parameters", bean, len(argv), tvpe.NumIn())
		}

		switch tvpe.NumOut() {
		case 1:
			if tvpe.Out(0) != bean.GetType() && tvpe.Out(0).Kind() != reflect.Interface {
				return fmt.Errorf("The return type from factory function of bean [%v] is [%v] instead of [%v]",
					bean, tvpe.Out(0), bean.GetType())
			}
		case 2:
			if tvpe.Out(0) != bean.GetType() && tvpe.Out(0).Kind() != reflect.Interface {
				return fmt.Errorf("The 1st return type from factory function of bean [%v] is [%v] instead of [%v]",
					bean, tvpe.Out(0), bean.GetType())
			}
			if tvpe.Out(1) != reflect.TypeOf((*error)(nil)).Elem() {
				return fmt.Errorf("The 2nd return type from factory function of bean [%v] is [%v] instead of error",
					bean, tvpe.Out(1))
			}
		default:
			return fmt.Errorf("Factory of bean [%v] should return (interface{}) or (interface{},error)", bean)
		}
	}

	switch bean.GetScope() {
	case Default:
	case Singleton:
	case Prototype:
		if bean.GetFinalize() != nil {
			return fmt.Errorf("A prototype bean can't have finalizer [%v].", *bean.GetFinalize())
		}
	default:
		return fmt.Errorf("Unkown scope [%v]", bean.GetScope())
	}

	for _, ps := range bean.GetProperties() {
		for _, p := range ps {

			ctx.parentByChild[p] = bean

			parentID := bean.GetID()
			for parentID == nil {
				if parent, present := ctx.parentByChild[bean]; present {
					parentID = parent.GetID()
					break
				}
			}

			childID := p.GetID()

			if parentID != nil && childID != nil {
				if !ctx.graph.AddDependency(*childID, *parentID) {
					return fmt.Errorf("Found a circle dependency from [%v] tp [%v]", *parentID, *childID)
				}
			}

			if e := ctx.addBean(p); e != nil {
				return fmt.Errorf("Can't add property bean [%v]. Cuased by: %v", p, e)
			}
		}
	}

	return nil
}

func (ctx *applicationContext) getBean(bean BeanI) (*reflect.Value, error) {

	switch bean.(type) {
	case ValueBeanI:
		i := bean.(ValueBeanI).GetValue()
		v := reflect.ValueOf(i)
		return &v, nil
	case ReferenceBeanI:
		id := bean.(ReferenceBeanI).GetID()
		v, e := ctx.getBean(ctx.beanById[*id])
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

	factory, factoryArgvBeans := bean.GetFactory()
	factoryArgvValues := make([]reflect.Value, len(factoryArgvBeans))

	for i, argvBean := range factoryArgvBeans {
		argvValue, e := ctx.getBean(argvBean)
		if e != nil {
			return nil, fmt.Errorf("Create input bean [%v] for factory [%v] failed", argvBean, factory)
		}
		factoryArgvValues[i] = *argvValue
	}

	factoryReturns := reflect.ValueOf(factory).Call(factoryArgvValues)

	var value *reflect.Value

	switch len(factoryReturns) {
	case 0:
		return nil, fmt.Errorf("Factory function of bean [%v] returns nothing", bean)
	case 1:
		value = &factoryReturns[0]
		if e, ok := value.Interface().(error); ok {
			return nil, fmt.Errorf("Create bean [%v] failed. Caused by: %v", bean, e)
		}
	default:
		value = &factoryReturns[0]
		if e, ok := value.Interface().(error); ok {
			return nil, fmt.Errorf("Create bean [%v] failed. Caused by: %v", bean, e)
		}
		if e, ok := factoryReturns[1].Interface().(error); ok {
			return nil, fmt.Errorf("Create bean [%v] failed. Caused by: %v", bean, e)
		}
	}

	return value, nil
}
