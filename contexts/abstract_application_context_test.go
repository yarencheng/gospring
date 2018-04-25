package contexts

import (
	"reflect"
	"testing"

	"github.com/yarencheng/gospring/beans"
)

func TestGetBean(t *testing.T) {

	ctx := NewAbstractApplicatoinContext([]*beans.BeanMetaData{
		beans.NewBeanMetaData("bean_1_id", beans.Singleton, reflect.TypeOf("")),
	})

	bean, _ := ctx.GetBean("bean_1_id")

	if _, ok := bean.(*string); !ok {
		t.Error("Failed to get a string")
	}
}

func TestGetBean_getTwice(t *testing.T) {

	ctx := NewAbstractApplicatoinContext([]*beans.BeanMetaData{
		beans.NewBeanMetaData("bean_1_id", beans.Singleton, reflect.TypeOf("")),
	})

	bean1, _ := ctx.GetBean("bean_1_id")
	bean2, _ := ctx.GetBean("bean_1_id")

	if bean1 != bean2 {
		t.Error("Failed to get a string")
	}
}

func TestGetBean_noSuchId(t *testing.T) {

	ctx := NewAbstractApplicatoinContext([]*beans.BeanMetaData{})

	if _, e := ctx.GetBean("bean_1_id"); e == nil {
		t.Error()
	}
}
