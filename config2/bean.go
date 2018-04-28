package config2

import (
	"fmt"
	"reflect"
)

type scope string

const (
	scopeDefault   scope = "scopeDefault"
	scopeSingleton scope = "scopeSingleton"
	scopePrototype scope = "scopePrototype"
)

type refBean struct {
	ref string
}

type propertyType string

const (
	propertyTypeBean  propertyType = "propertyTypeBean"
	propertyTypeRef   propertyType = "propertyTypeRef"
	propertyTypeValue propertyType = "propertyTypeValue"
)

type bean struct {
	type_         reflect.Type
	id            *string
	factoryFn     *reflect.Value
	factoryFnArgv []reflect.Value
	scope         scope
	pros          map[string]reflect.Value
	prosType      map[string]propertyType
}

type beans struct {
	beans []*bean
}

func Beans(bs ...*bean) *beans {
	return &beans{
		beans: bs,
	}
}

func Bean(i interface{}) *bean {
	return &bean{
		type_:         reflect.TypeOf(i),
		factoryFn:     getDefaultFactoryFn(reflect.TypeOf(i)),
		factoryFnArgv: make([]reflect.Value, 0),
		scope:         scopeDefault,
		pros:          make(map[string]reflect.Value),
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
	b.scope = scopeSingleton
	return b
}

func (b *bean) Prototype() *bean {
	b.scope = scopePrototype
	return b
}

func (b *bean) PropertyValue(name string, v interface{}) *bean {
	b.pros[name] = reflect.ValueOf(v)
	b.prosType[name] = propertyTypeValue
	return b
}

func (b *bean) PropertyBean(name string, bean *bean) *bean {
	b.pros[name] = reflect.ValueOf(bean)
	b.prosType[name] = propertyTypeBean
	return b
}

func (b *bean) PropertyRef(name string, ref string) *bean {
	b.pros[name] = reflect.ValueOf(&refBean{
		ref: ref,
	})
	b.prosType[name] = propertyTypeRef
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
		fn = func() (*string, error) {
			s := ""
			return &s, nil
		}
	case reflect.Int:
		fn = func() (*int, error) {
			i := int(0)
			return &i, nil
		}
	case reflect.Struct:
		fn = func() (interface{}, error) {
			return reflect.New(t).Interface(), nil
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
