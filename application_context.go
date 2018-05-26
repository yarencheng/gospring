package gospring

import (
	"container/list"
	"fmt"
	"reflect"
)

type applicationContext struct {
	graph         *graph
	beanById      map[string]BeanI
	parentByChild map[BeanI]BeanI
	singletons    map[BeanI]*reflect.Value
	singletonList *list.List
}

// NewApplicationContext creates an ApplicationContextI object
//
// NewApplicationContext(
//     Bean(...),
//     Bean(...),
// )
//
// It is a creation function to create an instance with the
// interface ApplicationContextI
func NewApplicationContext(beans ...BeanI) (ApplicationContextI, error) {
	list.New()
	ctx := applicationContext{
		graph:         newGraph(),
		beanById:      make(map[string]BeanI),
		parentByChild: make(map[BeanI]BeanI),
		singletons:    make(map[BeanI]*reflect.Value),
		singletonList: list.New(),
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

	for _, bean := range beans {
		if e := ctx.checkDependencyLoop(bean); e != nil {
			return nil, fmt.Errorf("Detect dependency loop. Cuased by: %v", e)
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
	cur := ctx.singletonList.Back()
	for cur != nil {
		bean := cur.Value.(BeanI)
		value := ctx.singletons[bean]
		if e := ctx.callFinalizeFunc(*value, bean); e != nil {
			return fmt.Errorf(
				"Can't call finalize function of bean [%v]. Caused by: [%v]",
				bean, e)
		}
		cur = cur.Prev()
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
		switch bean.(type) {
		case StructBeanI:
			if e := ctx.setRefBean(bean); e != nil {
				return fmt.Errorf("Replace reference beans for %s inside bean [%v] failed. Caused by: %v",
					des, bean, e)
			}
		case ReferenceBeanI:
			if target, present := ctx.beanById[*bean.GetID()]; present {
				bean.(ReferenceBeanI).SetReference(target)
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

	if e := ctx.checkScope(bean); e != nil {
		return fmt.Errorf("Scope is invalid. Caused by: %v", e)
	}

	for _, ps := range bean.GetProperties() {
		for _, p := range ps {
			if e := ctx.addBean(p); e != nil {
				return fmt.Errorf("Can't add property bean [%v]. Cuased by: %v", p, e)
			}
		}
	}

	return nil
}

func (ctx *applicationContext) checkDependencyLoop(bean BeanI) error {
	for _, ps := range bean.GetProperties() {
		for _, p := range ps {

			ctx.parentByChild[p] = bean

			parent := bean
			for parent.GetID() == nil {
				var present bool
				if parent, present = ctx.parentByChild[parent]; present {
					continue
				} else {
					break
				}
			}

			var parentID *string
			if parent != nil {
				parentID = parent.GetID()
			}
			childID := p.GetID()

			if parentID != nil && childID != nil {
				if !ctx.graph.AddDependency(*childID, *parentID) {
					return fmt.Errorf("Found a circle dependency from [%v] tp [%v]", *parentID, *childID)
				}
			}

			if e := ctx.checkDependencyLoop(p); e != nil {
				return fmt.Errorf("Detected dependency loop. Cuased by: %v", e)
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

func (ctx *applicationContext) checkScope(bean BeanI) error {
	switch bean.GetScope() {
	case Default:
	case Singleton:
	case Prototype:
		if bean.GetFinalize() != nil {
			return fmt.Errorf("A prototype bean can't have finalizer. ")
		}
	default:
		return fmt.Errorf("Unkown scope [%v]", bean.GetScope())
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

		isAssignable := func(return_, beant reflect.Type) bool {

			switch return_.Kind() {
			case reflect.Interface:
				return true
			case reflect.Ptr:
				return return_.Elem() == beant
			case reflect.Chan:
				return beant.Kind() == reflect.Chan
			default:
				return false
			}
		}

		switch tvpe.NumOut() {
		case 1:
			if !isAssignable(tvpe.Out(0), bean.GetType()) {
				return fmt.Errorf("The return type from factory function is [%v] instead of [&%v]",
					tvpe.Out(0), bean.GetType())
			}
		case 2:
			if !isAssignable(tvpe.Out(0), bean.GetType()) {
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

	if r, ok := bean.(ReferenceBeanI); ok {
		return ctx.getBean(r.GetReference())
	}

	switch bean.GetScope() {
	case Singleton:
		return ctx.getSingletonBean(bean)
	case Prototype:
		return ctx.getPrototypeBean(bean)
	case Default:
		return ctx.getSingletonBean(bean)
	default:
		return nil, fmt.Errorf("Scope [%T] of bean [%v] is not support", bean.GetScope(), bean)
	}
}

func (ctx *applicationContext) getSingletonBean(bean BeanI) (*reflect.Value, error) {

	if value, present := ctx.singletons[bean]; present {
		return value, nil
	}

	value, e := ctx.getPrototypeBean(bean)

	if e != nil {
		return nil, e
	}

	ctx.singletons[bean] = value
	ctx.singletonList.PushBack(bean)

	return value, nil
}
func (ctx *applicationContext) getPrototypeBean(bean BeanI) (*reflect.Value, error) {

	factory, factoryArgvBeans := bean.GetFactory()
	factoryV := reflect.ValueOf(factory)

	var value *reflect.Value
	var e error
	if value, e = ctx.createBeanByFactory(factoryV, factoryArgvBeans); e != nil {
		return nil, fmt.Errorf("Create bean failed. Cuased by: %v", e)
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

func (ctx *applicationContext) createBeanByFactory(fn reflect.Value, argvs []BeanI) (*reflect.Value, error) {

	values := make([]reflect.Value, len(argvs))

	for i, argv := range argvs {
		value, e := ctx.getBean(argv)
		if e != nil {
			return nil, fmt.Errorf("Can't get the [%d] argument from bean [%v]. Caused by: %v", i, argv, e)
		}

		fromType := value.Type()
		toType := fn.Type().In(i)

		if fromType.ConvertibleTo(toType) {
			values[i] = *value
		} else if fromType.Elem().ConvertibleTo(toType) {
			if Singleton == argv.GetScope() {
				return nil, fmt.Errorf("Can't inject a singleton to the [%d] non-pointer argument. ", i)
			}
			values[i] = value.Elem()
		} else {
			return nil, fmt.Errorf("The type of the [%d] argument isn't [%v] nor [%v]",
				i,
				fromType,
				fromType.Elem(),
			)
		}
	}

	returns := fn.Call(values)

	for i, _ := range returns {
		if returns[i].Type().Kind() == reflect.Interface {
			returns[i] = returns[i].Elem()
		}
	}

	var value *reflect.Value

	switch len(returns) {
	case 0:
		return nil, fmt.Errorf("Factory function returns nothing")
	case 1:
		value = &returns[0]
		if e, ok := value.Interface().(error); ok {
			return nil, fmt.Errorf("Get error from factory. Caused by: %v", e)
		}
	default:
		value = &returns[0]
		if e, ok := value.Interface().(error); ok {
			return nil, fmt.Errorf("Get error from factory. Caused by: %v", e)
		}
		if returns[1].IsValid() {
			if e, ok := returns[1].Interface().(error); ok {
				return nil, fmt.Errorf("Get error from factory. Caused by: %v", e)
			}
		}
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
		return fmt.Errorf("Bean [%v] can't be convert to [%v]",
			bean,
			toType,
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
			return fmt.Errorf("Bean [%v] can't be convert to [%v]",
				bean,
				toType,
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
		if rv[0].Type() == reflect.TypeOf((*error)(nil)).Elem() {
			if rv[0].IsNil() {
				return nil
			} else {
				return fmt.Errorf(
					"Function [%v] return an error. Caused by: %v",
					*initName,
					rv[0].Interface(),
				)
			}
		} else {
			return fmt.Errorf(
				"Function [%v] returns 1 unexpected value [%v] with type [%v]. ",
				*initName,
				rv[0].Interface(),
				rv[0].Type(),
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
		if bean.GetFinalize() != nil {
			return fmt.Errorf("Can't get finalizer [%v]", *finalName)
		}
		return nil // donothing
	}

	rv := finalFn.Func.Call([]reflect.Value{value})
	switch len(rv) {
	case 0:
		return nil
	case 1:
		if rv[0].Type() == reflect.TypeOf((*error)(nil)).Elem() {
			if rv[0].IsNil() {
				return nil
			} else {
				return fmt.Errorf(
					"Function [%v] return an error. Caused by: %v",
					*finalName,
					rv[0].Interface(),
				)
			}
		} else {
			return fmt.Errorf(
				"Function [%v] returns 1 unexpected value [%v] with type [%v]. ",
				*finalName,
				rv[0].Interface(),
				rv[0].Type(),
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
