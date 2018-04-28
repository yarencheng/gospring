package config2

import (
	"fmt"
	"reflect"
)

type scope string

const (
	singleton scope = "singleton"
	prototype scope = "prototype"
)

type bean struct {
	type_         reflect.Type
	id            *string
	factoryFn     *reflect.Value
	factoryFnArgv []reflect.Value
	scope         scope
}

func Bean(i interface{}) *bean {
	return &bean{
		type_:         reflect.TypeOf(i),
		factoryFn:     getDefaultFactoryFn(reflect.TypeOf(i)),
		factoryFnArgv: make([]reflect.Value, 0),
		scope:         singleton,
	}
}

func (b *bean) Id(id string) *bean {
	b.id = &id
	return b
}

func (b *bean) Factory(fn interface{}, argv ...interface{}) *bean {

	fnv := reflect.ValueOf(fn)
	b.factoryFn = &fnv
	b.factoryFnArgv = make([]reflect.Value, len(argv))

	for i, v := range argv {
		b.factoryFnArgv[i] = reflect.ValueOf(v)
	}
	return b
}

func (b *bean) Singleton() *bean {
	b.scope = singleton
	return b
}

func (b *bean) Prototype() *bean {
	b.scope = prototype
	return b
}

func (b *bean) new() (interface{}, error) {
	rs := b.factoryFn.Call(b.factoryFnArgv)
	if !rs[1].IsNil() {
		return nil, rs[1].Interface().(error)
	}
	return rs[0].Interface(), nil
}

func getDefaultFactoryFn(t reflect.Type) *reflect.Value {

	var fn interface{}

	switch t.Kind() {
	case reflect.String:
		fn = func() (interface{}, error) {
			return "", nil
		}
	case reflect.Int:
		fn = func() (interface{}, error) {
			return int(0), nil
		}
	default:
		fn = func() (interface{}, error) {
			return nil, fmt.Errorf(
				"There is no pre-defined factory function for the type [%v]",
				t.Name(),
			)
		}

	}

	v := reflect.ValueOf(fn)

	return &v
}
