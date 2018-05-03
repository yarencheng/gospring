package refactor

import "reflect"

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
		properties: make(map[string][]interface{}),
		init:       &i,
		finalize:   &f,
	}
}
