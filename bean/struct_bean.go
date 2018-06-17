package bean

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/yarencheng/gospring/v1"
)

type StructBean struct {
	id             string
	tvpe           reflect.Type
	scope          v1.Scope
	singletonValue reflect.Value
	factoryFn      reflect.Value
	factoryArgs    []reflect.Value
	startFn        reflect.Value
}

var defaultStruct StructBean = StructBean{
	scope: v1.Default,
}

func NewStructBeanV1(config v1.Bean) (*StructBean, error) {

	if err := checkType(config.Type); err != nil {
		return nil, err
	}

	scope := v1.Default
	if config.Scope == "" {
		scope = v1.Default
	}

	bean := &StructBean{
		id:    config.ID,
		tvpe:  config.Type,
		scope: scope,
	}

	if err := bean.initFactoryFn(config); err != nil {
		return nil, err
	}

	if err := bean.initStartFn(config); err != nil {
		return nil, err
	}

	return bean, nil
}

func checkType(tvpe reflect.Type) error {
	switch tvpe.Kind() {
	case reflect.Uintptr:
		fallthrough
	case reflect.Array:
		fallthrough
	case reflect.Chan:
		fallthrough
	case reflect.Func:
		fallthrough
	case reflect.Interface:
		fallthrough
	case reflect.Map:
		fallthrough
	case reflect.Slice:
		fallthrough
	case reflect.Ptr:
		fallthrough
	case reflect.UnsafePointer:
		return fmt.Errorf("[%v] is not a valid type for a bean", tvpe.Kind())
	}
	return nil
}

func (b *StructBean) GetID() string {
	return b.id
}

func (b *StructBean) GetValue() (reflect.Value, error) {
	switch b.scope {
	case v1.Default:
		fallthrough
	case v1.Singleton:
		if b.singletonValue.IsValid() {
			return b.singletonValue, nil
		}
		v, err := b.createValue()
		if err != nil {
			return reflect.Value{}, fmt.Errorf("Can't create the value, err: %v", err)
		}
		b.singletonValue = v
		return b.singletonValue, nil
	case v1.Prototype:
		return b.createValue()
	default:
		return reflect.Value{}, fmt.Errorf("Unknown scope [%v]", b.scope)
	}
}

func (b *StructBean) createValue() (reflect.Value, error) {
	var v reflect.Value

	if !b.factoryFn.IsValid() {
		v = reflect.New(b.tvpe)
	} else {
		vs := b.factoryFn.Call(b.factoryArgs)

		if len(vs) == 2 && !vs[1].IsNil() {
			return reflect.Value{}, fmt.Errorf("Can't create instance from the factory function. err: %v",
				vs[1].Interface())
		}
		v = vs[0]
	}

	if b.startFn.IsValid() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		in := make([]reflect.Value, 1)
		in[0] = v
		if b.startFn.Type().NumIn() > 1 {
			in = append(in, reflect.ValueOf(ctx))
		}

		done := make(chan int)
		var err []reflect.Value

		go func() {
			err = b.startFn.Call(in)
			done <- 1
		}()

		select {
		case <-ctx.Done():
		case <-done:
		}

		if err := ctx.Err(); err != nil {
			return reflect.Value{}, fmt.Errorf("StartFn [%v] timeout. err: %v",
				b.startFn.Type().Name(), err)
		}

		if len(err) > 0 && !err[0].IsNil() {
			return reflect.Value{}, fmt.Errorf("StartFn [%v] returned an error. err: %v",
				b.startFn.Type().Name(), err)
		}
	}

	return v, nil
}

func (b *StructBean) initFactoryFn(c v1.Bean) error {

	if nil == c.FactoryFn {
		return nil
	}

	b.factoryFn = reflect.ValueOf(c.FactoryFn)
	if b.factoryFn.Kind() != reflect.Func {
		return fmt.Errorf("FactoryFn should be a func instead of [%v]", b.factoryFn.Kind())
	}

	if b.factoryFn.Type().NumIn() != len(c.FactoryArgs) {
		return fmt.Errorf("FactoryFn [%v] need [%v] parameters",
			b.factoryFn.Type().Name(),
			b.factoryFn.Type().NumIn(),
		)
	}

	switch b.factoryFn.Type().NumOut() {
	case 2:
		if !b.factoryFn.Type().Out(1).AssignableTo(reflect.TypeOf((*error)(nil)).Elem()) {
			return fmt.Errorf("The second return value [%v] of the factory [%v] is not assignable to error",
				b.factoryFn.Type().Out(0),
				b.factoryFn.Type(),
			)
		}
		fallthrough
	case 1:
		if b.factoryFn.Type().Out(0).Kind() != reflect.Ptr {
			return fmt.Errorf("The first return value [%v] of the factory [%v] is not a pointer for [%v]",
				b.factoryFn.Type().Out(0),
				b.factoryFn.Type(),
				c.Type,
			)
		}
		if !b.factoryFn.Type().Out(0).Elem().AssignableTo(c.Type) {
			return fmt.Errorf("The first return value [%v] of the factory [%v] is not assignable to [%v]",
				b.factoryFn.Type().Out(0),
				b.factoryFn.Type(),
				c.Type,
			)
		}
	default:
		return fmt.Errorf("FactoryFn [%v] should return with 1 or 2 values",
			b.factoryFn.Type().Name(),
		)
	}

	b.factoryArgs = make([]reflect.Value, len(c.FactoryArgs))
	for i, a := range c.FactoryArgs {
		b.factoryArgs[i] = reflect.ValueOf(a)
	}

	for i := 0; i < b.factoryFn.Type().NumIn(); i++ {
		if !b.factoryFn.Type().In(i).AssignableTo(b.factoryArgs[i].Type()) {
			return fmt.Errorf("The [%v] argument [%v] of the factory [%v] is not assignable to [%v]",
				i,
				b.factoryFn.Type().In(i),
				b.factoryFn.Type(),
				b.factoryArgs[i].Type())
		}
	}

	return nil
}

func (b *StructBean) initStartFn(c v1.Bean) error {

	if c.StartFn == nil {
		ptrType := reflect.PtrTo(b.tvpe)
		m, exist := ptrType.MethodByName("Start")

		if !exist {
			return nil
		}
		b.startFn = m.Func
	} else if reflect.TypeOf(c.StartFn).Kind() == reflect.String {

	} else if reflect.TypeOf(c.StartFn).Kind() == reflect.Func {

	} else {
		return fmt.Errorf("StartFn should be a func or a name of the func for [%v]", b.tvpe)
	}

	switch b.startFn.Type().NumIn() {
	case 2:
		if b.startFn.Type().In(1).AssignableTo(reflect.TypeOf((*context.Context)(nil)).Elem()) {
			return fmt.Errorf("The 2nd argument of the start func [%v] should be a context.Context", b.startFn)
		}
		fallthrough
	case 1:
		if b.startFn.Type().In(0).Kind() != reflect.Ptr ||
			!b.startFn.Type().In(0).Elem().AssignableTo(c.Type) {
			return fmt.Errorf("The 1st argument of the start func [%v] should be a context.Context", b.startFn)
		}
	default:
		return fmt.Errorf("Ony 1 or 2 argument for a start function")
	}

	switch b.startFn.Type().NumOut() {
	case 1:
		if b.startFn.Type().In(0).AssignableTo(reflect.TypeOf((*error)(nil)).Elem()) {
			return fmt.Errorf("The 1st return value of the start func [%v] should be an error", b.startFn)
		}
		fallthrough
	case 0:
	default:
		return fmt.Errorf("Ony 0 or 1 return value for a start function")
	}

	return nil
}
