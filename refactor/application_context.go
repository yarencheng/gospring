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

	if id := bean.GetID(); id != nil {
		if _, ok := bean.(ReferenceBeanI); !ok {
			if _, present := ctx.beanById[*id]; present {
				return fmt.Errorf("ID [%v] already exist", *id)
			}
			ctx.beanById[*id] = bean
		}
	}

	if _, ok := bean.(StructBeanI); ok {
		if tvpe := bean.GetType(); tvpe != nil {
			if tvpe.Kind() == reflect.Ptr {
				return fmt.Errorf("Type of bean [%v] is a pointer instead of struct", bean)
			}
		}
	}

	if fn, argv := bean.GetFactory(); fn != nil {
		tvpe := reflect.TypeOf(fn)

		if tvpe.Kind() != reflect.Func {
			return fmt.Errorf("Factory of bean [%v] is not a function but a [%v]", bean, tvpe.Kind())
		}

		if len(argv) != tvpe.NumIn() {
			return fmt.Errorf("Factory of bean [%v] need [%v] instead of [%v] parameters", bean, len(argv), tvpe.NumIn())
		}

		isPointer := func(ptr, strvct reflect.Type) bool {
			switch ptr.Kind() {
			case reflect.Interface:
				return true
			case reflect.Ptr:
				return ptr.Elem() == strvct
			default:
				return false
			}
		}

		switch tvpe.NumOut() {
		case 1:
			if !isPointer(tvpe.Out(0), bean.GetType()) {
				return fmt.Errorf("The return type from factory function of bean [%v] is [%v] instead of [&%v]",
					bean, tvpe.Out(0), bean.GetType())
			}
		case 2:
			if !isPointer(tvpe.Out(0), bean.GetType()) {
				return fmt.Errorf("The 1st return type from factory function of bean [%v] is [%v] instead of [&%v]",
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

	if rbean, ok := bean.(ReferenceBeanI); ok {
		if pbean, p := ctx.beanById[rbean.GetReference()]; p {
			bean = pbean
		} else {
			return nil, fmt.Errorf("There is no bean with ID [%v]", rbean.GetReference())
		}
	}

	switch bean.GetScope() {
	case Singleton:
		if bean.GetID() != nil {
			return ctx.getSingletonBean(bean)
		} else {
			return ctx.getPrototypeBean(bean)
		}
	case Prototype:
		return ctx.getPrototypeBean(bean)
	case Default:
		return ctx.getSingletonBean(bean)
	default:
		return nil, fmt.Errorf("Scope [%T] of bean [%v] is not support", bean.GetScope(), bean)
	}
}

func (ctx *applicationContext) getSingletonBean(bean BeanI) (*reflect.Value, error) {

	if value, present := ctx.singletons[*bean.GetID()]; present {
		return value, nil
	}

	value, e := ctx.getPrototypeBean(bean)

	if e != nil {
		return nil, e
	}

	return value, nil
}
func (ctx *applicationContext) getPrototypeBean(bean BeanI) (*reflect.Value, error) {

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

	for i, _ := range factoryReturns {
		if factoryReturns[i].Type().Kind() == reflect.Interface {
			factoryReturns[i] = factoryReturns[i].Elem()
		}
	}

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
