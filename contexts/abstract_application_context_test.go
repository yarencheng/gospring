package contexts

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/yarencheng/gospring/beans"
)

func TestGetBean_getSingleton(t *testing.T) {

	ctx, _ := NewAbstractApplicatoinContext([]*beans.BeanMetaData_old{
		beans.NewBeanMetaData_old("bean_1_id", beans.Singleton, reflect.TypeOf(""), nil),
	})

	bean, _ := ctx.GetBean("bean_1_id")

	if _, ok := bean.(*string); !ok {
		t.Error()
	}
}

func TestGetBean_getPrototype(t *testing.T) {

	ctx, _ := NewAbstractApplicatoinContext([]*beans.BeanMetaData_old{
		beans.NewBeanMetaData_old("bean_1_id", beans.Prototype, reflect.TypeOf(""), nil),
	})

	bean, _ := ctx.GetBean("bean_1_id")

	if _, ok := bean.(*string); !ok {
		t.Error()
	}
}

func TestGetBean_getSingletonTwice(t *testing.T) {

	ctx, _ := NewAbstractApplicatoinContext([]*beans.BeanMetaData_old{
		beans.NewBeanMetaData_old("bean_1_id", beans.Singleton, reflect.TypeOf(""), nil),
	})

	bean1, _ := ctx.GetBean("bean_1_id")
	bean2, _ := ctx.GetBean("bean_1_id")

	if bean1 != bean2 {
		t.Error()
	}
}

func TestGetBean_getPrototypeTwice(t *testing.T) {

	ctx, _ := NewAbstractApplicatoinContext([]*beans.BeanMetaData_old{
		beans.NewBeanMetaData_old("bean_1_id", beans.Prototype, reflect.TypeOf(""), nil),
	})

	bean1, _ := ctx.GetBean("bean_1_id")
	bean2, _ := ctx.GetBean("bean_1_id")

	if bean1 == bean2 {
		t.Error()
	}
}

func TestGetBean_noSuchId(t *testing.T) {

	ctx, _ := NewAbstractApplicatoinContext([]*beans.BeanMetaData_old{})

	if _, e := ctx.GetBean("bean_1_id"); e == nil {
		t.Error()
	}
}

func TestNewAbstractApplicatoinContext_idConfilck(t *testing.T) {

	_, e := NewAbstractApplicatoinContext([]*beans.BeanMetaData_old{
		beans.NewBeanMetaData_old("bean_1_id", beans.Prototype, reflect.TypeOf(""), nil),
		beans.NewBeanMetaData_old("bean_1_id", beans.Prototype, reflect.TypeOf(""), nil),
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

	ctx, _ := NewAbstractApplicatoinContext([]*beans.BeanMetaData_old{
		beans.NewBeanMetaData_old("bean_1_id", beans.Prototype, reflect.TypeOf(Bean{}), []beans.PropertyMetaData_old{
			*beans.NewPropertyMetaData_old("Property_1", "", expected),
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

	ctx, _ := NewAbstractApplicatoinContext([]*beans.BeanMetaData_old{
		beans.NewBeanMetaData_old("bean_1_id", beans.Prototype, reflect.TypeOf(Bean{}), []beans.PropertyMetaData_old{
			*beans.NewPropertyMetaData_old("Property_1", "", strconv.Itoa(expected)),
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

	ctx, _ := NewAbstractApplicatoinContext([]*beans.BeanMetaData_old{
		beans.NewBeanMetaData_old("bean_1_id", beans.Singleton, reflect.TypeOf(Bean1{}), nil),
		beans.NewBeanMetaData_old("bean_2_id", beans.Singleton, reflect.TypeOf(Bean2{}), []beans.PropertyMetaData_old{
			*beans.NewPropertyMetaData_old("Bean1", "bean_1_id", ""),
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

func TestGetBean_getBeanWithProperty_singletonBean_notPointer(t *testing.T) {

	type Bean1 struct{}
	type Bean2 struct {
		Bean1 Bean1 // not pointer
	}

	ctx, _ := NewAbstractApplicatoinContext([]*beans.BeanMetaData_old{
		beans.NewBeanMetaData_old("bean_1_id", beans.Singleton, reflect.TypeOf(Bean1{}), nil),
		beans.NewBeanMetaData_old("bean_2_id", beans.Singleton, reflect.TypeOf(Bean2{}), []beans.PropertyMetaData_old{
			*beans.NewPropertyMetaData_old("Bean1", "bean_1_id", ""),
		}),
	})

	_, e := ctx.GetBean("bean_2_id")

	if e == nil {
		t.Error()
	}
}

func TestGetBean_getBeanWithProperty_prototypeBean_copyPointer(t *testing.T) {

	type Bean1 struct {
		a int
	}
	type Bean2 struct {
		Bean1 *Bean1
	}

	ctx, _ := NewAbstractApplicatoinContext([]*beans.BeanMetaData_old{
		beans.NewBeanMetaData_old("bean_1_id", beans.Prototype, reflect.TypeOf(Bean1{}), nil),
		beans.NewBeanMetaData_old("bean_2_id", beans.Prototype, reflect.TypeOf(Bean2{}), []beans.PropertyMetaData_old{
			*beans.NewPropertyMetaData_old("Bean1", "bean_1_id", ""),
		}),
	})

	bean1P, e1 := ctx.GetBean("bean_2_id")
	bean2P, e2 := ctx.GetBean("bean_2_id")

	if e1 != nil {
		t.Fatal(e1)
	}

	if e2 != nil {
		t.Fatal(e2)
	}

	bean1 := bean1P.(*Bean2)

	bean2 := bean2P.(*Bean2)

	if bean1.Bean1 == nil {
		t.Error()
	}

	if bean2.Bean1 == nil {
		t.Error()
	}

	if bean1.Bean1 == bean2.Bean1 {
		t.Errorf("bean1.Bean1=[%p] bean2.Bean1=[%p]", bean1.Bean1, bean2.Bean1)
	}
}

func TestGetBean_getBeanWithProperty_prototypeBean_copyValue(t *testing.T) {

	type Bean1 struct {
		I int
	}
	type Bean2 struct {
		Bean1 Bean1
	}

	expected := 123

	ctx, _ := NewAbstractApplicatoinContext([]*beans.BeanMetaData_old{
		beans.NewBeanMetaData_old("bean_1_id", beans.Prototype, reflect.TypeOf(Bean1{}), []beans.PropertyMetaData_old{
			*beans.NewPropertyMetaData_old("I", "", strconv.Itoa(expected)),
		}),
		beans.NewBeanMetaData_old("bean_2_id", beans.Prototype, reflect.TypeOf(Bean2{}), []beans.PropertyMetaData_old{
			*beans.NewPropertyMetaData_old("Bean1", "bean_1_id", ""),
		}),
	})

	beanP, e := ctx.GetBean("bean_2_id")

	if e != nil {
		t.Fatal(e)
	}

	bean := beanP.(*Bean2)

	if bean.Bean1.I != expected {
		t.Errorf("bean.Bean1.I=[%v] expected=[%v]", bean.Bean1.I, expected)
	}
}