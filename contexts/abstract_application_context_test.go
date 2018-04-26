package contexts

import (
	"reflect"
	"strconv"
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

func TestGetBean_getBeanWithProperty_string(t *testing.T) {

	type Bean struct {
		Property_1 string
	}

	expected := "hahahaha"

	ctx, _ := NewAbstractApplicatoinContext([]*beans.BeanMetaData{
		beans.NewBeanMetaData("bean_1_id", beans.Prototype, reflect.TypeOf(Bean{}), []beans.PropertyMetaData{
			*beans.NewPropertyMetaData("Property_1", "", expected),
		}),
	})

	beanP, e := ctx.GetBean("bean_1_id")
	if e != nil {
		t.Fatal(e)
	}
	bean := beanP.(*Bean)

	if expected != bean.Property_1 {
		t.Errorf("expected=[%v] bean.Property_1=[%v]", expected, bean.Property_1)
	}
}

func TestGetBean_getBeanWithProperty_int(t *testing.T) {

	type Bean struct {
		Property_1 int
	}

	expected := 123

	ctx, _ := NewAbstractApplicatoinContext([]*beans.BeanMetaData{
		beans.NewBeanMetaData("bean_1_id", beans.Prototype, reflect.TypeOf(Bean{}), []beans.PropertyMetaData{
			*beans.NewPropertyMetaData("Property_1", "", strconv.Itoa(expected)),
		}),
	})

	beanP, e := ctx.GetBean("bean_1_id")
	if e != nil {
		t.Fatal(e)
	}
	bean := beanP.(*Bean)

	if expected != bean.Property_1 {
		t.Errorf("expected=[%v] bean.Property_1=[%v]", expected, bean.Property_1)
	}
}

func TestGetBean_getBeanWithProperty_singletonBean(t *testing.T) {

	type Bean1 struct{}
	type Bean2 struct {
		Bean1 *Bean1
	}

	ctx, _ := NewAbstractApplicatoinContext([]*beans.BeanMetaData{
		beans.NewBeanMetaData("bean_1_id", beans.Singleton, reflect.TypeOf(Bean1{}), nil),
		beans.NewBeanMetaData("bean_2_id", beans.Singleton, reflect.TypeOf(Bean2{}), []beans.PropertyMetaData{
			*beans.NewPropertyMetaData("Bean1", "bean_1_id", ""),
		}),
	})

	bean1P, e1 := ctx.GetBean("bean_1_id")
	if e1 != nil {
		t.Fatal(e1)
	}

	bean2P, e2 := ctx.GetBean("bean_2_id")
	if e2 != nil {
		t.Fatal(e2)
	}

	bean1 := bean1P.(*Bean1)
	bean2 := bean2P.(*Bean2)

	if bean1 != bean2.Bean1 {
		t.Errorf("bean1=[%p] bean2.Bean1=[%p]", bean1, bean2.Bean1)
	}
}
