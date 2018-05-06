package refactor

import (
	"reflect"
)

func Ref(id string) ReferenceBeanI {
	return &referenceBean{
		id: id,
	}
}

func Bean(tvpe interface{}) StructBeanI {
	i := "Init"
	f := "Finalize"
	return &structBean{
		tvpe:       reflect.TypeOf(tvpe),
		properties: make(map[string][]BeanI),
		init:       &i,
		finalize:   &f,
		scope:      Default,
		factoryFn: func() interface{} {
			return reflect.New(reflect.TypeOf(tvpe))
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