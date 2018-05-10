package gospring

import (
	"fmt"
	"reflect"
)

type applicationContext struct {
	graph         *graph
	beanById      map[string]BeanI
	parentByChild map[BeanI]BeanI
	singletons    map[string]*reflect.Value
}

func NewApplicationContext(beans ...BeanI) (ApplicationContextI, error) {

	ctx := applicationContext{
		graph:         newGraph(),
		beanById:      make(map[string]BeanI),
		parentByChild: make(map[BeanI]BeanI),
		singletons:    make(map[string]*reflect.Value),
	}

	for _, bean := range beans {
		if e := ctx.addBean(bean); e != nil {
			return nil, fmt.Errorf("Can't add bean [%v]. Cuased by: %v", bean, e)
		}
	}

	for _, bean := range beans {
		if e := ctx.setRefBean(bean); e != nil {
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
	for id, value := range ctx.singletons {
		bean := ctx.beanById[id]
		if e := ctx.callFinalizeFunc(*value, bean); e != nil {
			return fmt.Errorf(
				"Can't call finalize function of bean [%v]. Caused by: [%v]",
				bean, e)
		}
	}
	return nil
}

func (ctx *applicationContext) setRefBean(parent BeanI) error {

	// (key,value) = (bean,description)
	bs := make(map[BeanI]string)

	_, argvs := parent.GetFactory()
	for i, argv := range argvs {
		bs[argv] = fmt.Sprintf("the number [%d] argument of factory function", i)
	}

	if sbean, ok := parent.(*structBean); ok {
		for name, ps := range sbean.GetProperties() {
			for _, p := range ps {
				bs[p] = fmt.Sprintf("the field [%s]", name)
			}
		}
	}

	for bean, des := range bs {
		if _, ok := bean.(StructBeanI); ok {
			if e := ctx.setRefBean(bean); e != nil {
				return fmt.Errorf("Replace reference beans for %s inside bean [%v] failed. Caused by: %v",
					des, bean, e)
			}
			continue
		}
		if rbean, ok := bean.(ReferenceBeanI); ok {
			if target, present := ctx.beanById[*bean.GetID()]; present {
				rbean.SetReference(target)
			} else {
				return fmt.Errorf("Can find ID [%v] of [%v] inside bean [%v]",
					*bean.GetID(), des, bean)
			}
		}
	}

	return nil
}

func (ctx *applicationContext) addBean(bean BeanI) error {

	if _, ok := bean.(ReferenceBeanI); ok {
		return nil
	}

	if e := ctx.addBeanById(bean); e != nil {
		return e
	}

	if e := ctx.checkType(bean); e != nil {
		return fmt.Errorf("Type is invalid. Caused by: %v", e)
	}

	if e := ctx.checkFactory(bean); e != nil {
		return fmt.Errorf("Factory is invalid. Caused by: %v", e)
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

func (ctx *applicationContext) addBeanById(bean BeanI) error {
	if id := bean.GetID(); id != nil {
		if _, present := ctx.beanById[*id]; present {
			return fmt.Errorf("ID [%v] already exist", *id)
		}
		ctx.beanById[*id] = bean
	}
	return nil
}

func (ctx *applicationContext) checkType(bean BeanI) error {
	if tvpe := bean.GetType(); tvpe != nil {
		if tvpe.Kind() == reflect.Ptr {
			return fmt.Errorf("It can't be a pinter")
		}
	}
	return nil
}

func (ctx *applicationContext) checkFactory(bean BeanI) error {
	if fn, argv := bean.GetFactory(); fn != nil {
		tvpe := reflect.TypeOf(fn)

		if tvpe.Kind() != reflect.Func {
			return fmt.Errorf("Factory is not a function but a [%v]", tvpe.Kind())
		}

		if len(argv) != tvpe.NumIn() {
			return fmt.Errorf("Factory need [%v] instead of [%v] parameters", len(argv), tvpe.NumIn())
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
				return fmt.Errorf("The return type from factory function is [%v] instead of [&%v]",
					tvpe.Out(0), bean.GetType())
			}
		case 2:
			if !isPointer(tvpe.Out(0), bean.GetType()) {
				return fmt.Errorf("The 1st return type from factory function is [%v] instead of [&%v]",
					tvpe.Out(0), bean.GetType())
			}
			if tvpe.Out(1) != reflect.TypeOf((*error)(nil)).Elem() {
				return fmt.Errorf("The 2nd return type from factory function is [%v] instead of error",
					tvpe.Out(1))
			}
		default:
			return fmt.Errorf("Factory should return (interface{}) or (interface{},error)")
		}
	}
	return nil
}

func (ctx *applicationContext) getBean(bean BeanI) (*reflect.Value, error) {

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

	if bean.GetID() != nil {
		if value, present := ctx.singletons[*bean.GetID()]; present {
			return value, nil
		}
	}

	value, e := ctx.getPrototypeBean(bean)

	if e != nil {
		return nil, e
	}

	if bean.GetID() != nil {
		ctx.singletons[*bean.GetID()] = value
	}

	return value, nil
}
func (ctx *applicationContext) getPrototypeBean(bean BeanI) (*reflect.Value, error) {

	factory, factoryArgvBeans := bean.GetFactory()
	factoryV := reflect.ValueOf(factory)
	factoryArgvValues := make([]reflect.Value, len(factoryArgvBeans))

	for i, argvBean := range factoryArgvBeans {
		argvValue, e := ctx.getBean(argvBean)
		if e != nil {
			return nil, fmt.Errorf("Create input bean [%v] for factory [%v] failed", argvBean, factory)
		}
		factoryArgvValues[i] = *argvValue

		fromType := factoryArgvValues[i].Type()
		toType := factoryV.Type().In(i)

		if fromType.ConvertibleTo(toType) {
			factoryArgvValues[i] = *argvValue
		} else if fromType.Elem().ConvertibleTo(toType) {
			if Singleton == argvBean.GetScope() {
				return nil, fmt.Errorf("Can't inject a singleton to non-pointer filed")
			}
			factoryArgvValues[i] = argvValue.Elem()
		} else {
			return nil, fmt.Errorf("Parameter type of factory isn't [%v] nor [%v]",
				fromType,
				fromType.Elem(),
			)
		}

	}

	factoryReturns := factoryV.Call(factoryArgvValues)

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

	for name, ps := range bean.GetProperties() {

		field := value.Elem().FieldByName(name)

		switch field.Type().Kind() {
		case reflect.Slice:
			if e := ctx.injectSlice(field, ps...); e != nil {
				return nil, fmt.Errorf("Can't inject field [%v] into bean [%v]. Caused by: %v", name, bean, e)
			}
		case reflect.Map:
			return nil, fmt.Errorf("TODO")
		default:
			if e := ctx.inject(field, ps[0]); e != nil {
				return nil, fmt.Errorf("Can't inject field [%v] into bean [%v]. Caused by: %v", name, bean, e)
			}
		}

	}

	if e := ctx.callInitFunc(*value, bean); e != nil {
		return nil, fmt.Errorf("Can't call initial function of bean [%v]. Caused by: [%v]", bean, e)
	}

	return value, nil
}

func (ctx *applicationContext) inject(field reflect.Value, bean BeanI) error {

	pv, e := ctx.getBean(bean)
	if e != nil {
		return fmt.Errorf("Can't get bean [%v]. Caused by: %v", bean, e)
	}

	fromType := pv.Type()
	toType := field.Type()

	if fromType.ConvertibleTo(toType) {
		field.Set(*pv)
	} else if fromType.Elem().ConvertibleTo(toType) {
		if Singleton == bean.GetScope() {
			return fmt.Errorf("Can't inject a singleton to non-pointer filed")
		}
		field.Set(pv.Elem())
	} else {
		return fmt.Errorf("Bean [%v] for field isn't [%v] nor [%v]",
			bean,
			fromType,
			fromType.Elem(),
		)
	}

	return nil
}

func (ctx *applicationContext) injectSlice(field reflect.Value, beans ...BeanI) error {

	slice := reflect.MakeSlice(field.Type(), len(beans), len(beans))

	for i, bean := range beans {
		pv, e := ctx.getBean(bean)
		if e != nil {
			return fmt.Errorf("Can't get bean [%v]. Caused by: %v", bean, e)
		}

		fromType := pv.Type()
		toType := field.Type().Elem()

		if fromType.ConvertibleTo(toType) {
			slice.Index(i).Set(*pv)
		} else if fromType.Elem().ConvertibleTo(toType) {
			if Singleton == bean.GetScope() {
				return fmt.Errorf("Can't inject a singleton to non-pointer filed")
			}
			slice.Index(i).Set(pv.Elem())
		} else {
			return fmt.Errorf("Bean [%v] for field isn't [%v] nor [%v]",
				bean,
				fromType,
				fromType.Elem(),
			)
		}
	}

	field.Set(slice)

	return nil
}

func (ctx *applicationContext) callInitFunc(value reflect.Value, bean BeanI) error {

	initName := bean.GetInit()
	if initName == nil {
		s := DefaultInitFunc
		initName = &s
	}

	initFn, ok := value.Type().MethodByName(*initName)
	if !ok {
		if bean.GetInit() != nil {
			return fmt.Errorf("Can't get initializer [%v]", *initName)
		}
		return nil // donothing
	}

	rv := initFn.Func.Call([]reflect.Value{value})
	switch len(rv) {
	case 0:
		return nil
	case 1:
		if e, ok := rv[0].Interface().(error); ok {
			return fmt.Errorf(
				"Function [%v] return an error. Caused by: %v",
				*initName,
				e,
			)
		} else {
			return fmt.Errorf(
				"Function [%v] returns 1 unexpected value [%v]",
				*initName,
				rv[0].Interface(),
			)
		}
	default:
		return fmt.Errorf(
			"Function [%v] returns %d unexpected value",
			*initName,
			len(rv),
		)
	}
}

func (ctx *applicationContext) callFinalizeFunc(value reflect.Value, bean BeanI) error {

	finalName := bean.GetFinalize()
	if finalName == nil {
		s := DefaultFinalizeFunc
		finalName = &s
	}

	finalFn, ok := value.Type().MethodByName(*finalName)
	if !ok {
		if bean.GetInit() != nil {
			return fmt.Errorf("Can't get finalizer [%v]", *finalName)
		}
		return nil // donothing
	}

	rv := finalFn.Func.Call([]reflect.Value{value})
	switch len(rv) {
	case 0:
		return nil
	case 1:
		if e, ok := rv[0].Interface().(error); ok {
			return fmt.Errorf(
				"Function [%v] return an error. Caused by: %v",
				*finalName,
				e,
			)
		} else {
			return fmt.Errorf(
				"Function [%v] returns 1 unexpected value [%v]",
				*finalName,
				rv[0].Interface(),
			)
		}
	default:
		return fmt.Errorf(
			"Function [%v] returns %d unexpected value",
			*finalName,
			len(rv),
		)
	}
}
