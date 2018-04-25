package contexts

import (
	"reflect"
	"testing"

	"github.com/yarencheng/gospring/beans"
)

func TestGetBean_getSingleton(t *testing.T) {

	ctx, _ := NewAbstractApplicatoinContext([]*beans.BeanMetaData{
		beans.NewBeanMetaData("bean_1_id", beans.Singleton, reflect.TypeOf(""), nil),
	})

	bean, _ := ctx.GetBean("bean_1_id")

	if _, ok := bean.(*string); !ok {
		t.Error()
	}
}

func TestGetBean_getPrototype(t *testing.T) {

	ctx, _ := NewAbstractApplicatoinContext([]*beans.BeanMetaData{
		beans.NewBeanMetaData("bean_1_id", beans.Prototype, reflect.TypeOf(""), nil),
	})

	bean, _ := ctx.GetBean("bean_1_id")

	if _, ok := bean.(*string); !ok {
		t.Error()
	}
}

func TestGetBean_getSingletonTwice(t *testing.T) {

	ctx, _ := NewAbstractApplicatoinContext([]*beans.BeanMetaData{
		beans.NewBeanMetaData("bean_1_id", beans.Singleton, reflect.TypeOf(""), nil),
	})

	bean1, _ := ctx.GetBean("bean_1_id")
	bean2, _ := ctx.GetBean("bean_1_id")

	if bean1 != bean2 {
		t.Error()
	}
}

func TestGetBean_getPrototypeTwice(t *testing.T) {

	ctx, _ := NewAbstractApplicatoinContext([]*beans.BeanMetaData{
		beans.NewBeanMetaData("bean_1_id", beans.Prototype, reflect.TypeOf(""), nil),
	})

	bean1, _ := ctx.GetBean("bean_1_id")
	bean2, _ := ctx.GetBean("bean_1_id")

	if bean1 == bean2 {
		t.Error()
	}
}

func TestGetBean_noSuchId(t *testing.T) {

	ctx, _ := NewAbstractApplicatoinContext([]*beans.BeanMetaData{})

	if _, e := ctx.GetBean("bean_1_id"); e == nil {
		t.Error()
	}
}

func TestNewAbstractApplicatoinContext_idConfilck(t *testing.T) {

	_, e := NewAbstractApplicatoinContext([]*beans.BeanMetaData{
		beans.NewBeanMetaData("bean_1_id", beans.Prototype, reflect.TypeOf(""), nil),
		beans.NewBeanMetaData("bean_1_id", beans.Prototype, reflect.TypeOf(""), nil),
	})

	if e == nil {
		t.Error()
	}
}
