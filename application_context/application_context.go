package application_context

import "github.com/yarencheng/gospring/v1"

type ApplicationContext struct {
}

func New(beans ...v1.Bean) ApplicationContextI {
	return &ApplicationContext{}
}

func (c *ApplicationContext) GetByID(id string) interface{} {
	return nil
}
