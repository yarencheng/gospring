package gospring

import (
	"reflect"
)

func Ref(id string) ReferenceBeanI {
	return &referenceBean{
		id: id,
	}
}

func Bean(value interface{}) StructBeanI {
	return &structBean{
		tvpe:       reflect.TypeOf(value),
		value:      reflect.ValueOf(value),
		properties: make(map[string][]BeanI),
		init:       nil,
		finalize:   nil,
		scope:      Default,
		factoryFn: func() interface{} {
			v := reflect.New(reflect.TypeOf(value))
			v.Elem().Set(reflect.ValueOf(value))
			return v.Interface()
		},
	}
}

func Beans(values ...interface{}) []BeanI {

	beans := make([]BeanI, len(values))

	for i, value := range values {
		switch value.(type) {
		case StructBeanI:
			beans[i] = value.(BeanI)
			continue
		case ReferenceBeanI:
			beans[i] = value.(BeanI)
			continue
		case ValueBeanI:
			beans[i] = value.(BeanI)
			continue
		default:
		}

		bean := &valueBean{
			value: value,
		}
		beans[i] = bean
	}

	return beans
}

func Chan(value interface{}, buffer int) StructBeanI {

	c := reflect.ChanOf(reflect.BothDir, reflect.TypeOf(value))
	dummy := reflect.MakeChan(c, 1)

	return &structBean{
		tvpe:       dummy.Type(),
		properties: make(map[string][]BeanI),
		init:       nil,
		finalize:   nil,
		scope:      Default,
		factoryFn: func() interface{} {
			v := reflect.MakeChan(c, buffer)
			return v.Interface()
		},
	}
}
