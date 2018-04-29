package gospring

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
	pros          map[string][]reflect.Value
	prosType      map[string]propertyType
	initFn        *reflect.Value
	finalizeFn    *reflect.Value
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
		pros:          make(map[string][]reflect.Value),
		prosType:      make(map[string]propertyType),
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

func (b *bean) PropertyValue(name string, vs ...interface{}) *bean {

	rvs := make([]reflect.Value, len(vs))
	for i, v := range vs {
		rvs[i] = reflect.ValueOf(v)
	}

	b.pros[name] = rvs
	b.prosType[name] = propertyTypeValue
	return b
}

func (b *bean) PropertyBean(name string, beans ...*bean) *bean {

	rvs := make([]reflect.Value, len(beans))
	for i, bean := range beans {
		rvs[i] = reflect.ValueOf(bean)
	}

	b.pros[name] = rvs
	b.prosType[name] = propertyTypeBean
	return b
}

func (b *bean) PropertyRef(name string, refs ...string) *bean {

	rvs := make([]reflect.Value, len(refs))
	for i, ref := range refs {
		rvs[i] = reflect.ValueOf(&refBean{
			ref: ref,
		})
	}

	b.pros[name] = rvs
	b.prosType[name] = propertyTypeRef
	return b
}

func (b *bean) Init(fn interface{}) *bean {
	fnv := reflect.ValueOf(fn)
	b.initFn = &fnv
	return b
}

func (b *bean) Finalize(fn interface{}) *bean {
	fnv := reflect.ValueOf(fn)
	b.finalizeFn = &fnv
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
