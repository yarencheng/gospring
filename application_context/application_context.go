package application_context

import (
	"fmt"
	"reflect"

	"github.com/yarencheng/gospring/bean"
	"github.com/yarencheng/gospring/v1"
)

type ApplicationContext struct {
	beans []bean.BeanI
}

func New(configs ...interface{}) (*ApplicationContext, error) {
	var err error
	beans := make([]bean.BeanI, len(configs))
	for i, config := range configs {
		switch config.(type) {
		case v1.Bean:
			beans[i], err = bean.NewStructBeanV1(config.(v1.Bean))
			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("[%v] is not a valid config struct", reflect.TypeOf(config).Name())
		}
	}

	return &ApplicationContext{
		beans: beans,
	}, nil
}

func (c *ApplicationContext) GetByID(id string) interface{} {
	return nil
}
