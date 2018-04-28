package config2

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func Test_applicationContext_GetBean_defaultScope_string(t *testing.T) {

	// arrange
	ctx, _ := ApplicationContext(Beans(
		Bean("").Id("id"),
	))

	// action
	bean1, e1 := ctx.GetBean("id")
	if e1 != nil {
		t.FailNow()
	}
	bean2, e2 := ctx.GetBean("id")
	if e2 != nil {
		t.FailNow()
	}

	// aasert

	s1, ok1 := bean1.(*string)
	assert.True(t, ok1)
	s2, ok2 := bean2.(*string)
	assert.True(t, ok2)

	assert.Equal(t, unsafe.Pointer(s1), unsafe.Pointer(s2))
}

func Test_applicationContext_GetBean_defaultScope_struct(t *testing.T) {

	// arrange
	type beanStruct struct{ s string }
	ctx, _ := ApplicationContext(Beans(
		Bean(beanStruct{}).Id("id"),
	))

	// action
	bean1, e1 := ctx.GetBean("id")
	if e1 != nil {
		t.FailNow()
	}
	bean2, e2 := ctx.GetBean("id")
	if e2 != nil {
		t.FailNow()
	}

	// aasert

	s1, ok1 := bean1.(*beanStruct)
	assert.True(t, ok1)
	s2, ok2 := bean2.(*beanStruct)
	assert.True(t, ok2)

	assert.Equal(t, unsafe.Pointer(s1), unsafe.Pointer(s2))
}
func Test_applicationContext_GetBean_singleton_string(t *testing.T) {

	// arrange
	ctx, _ := ApplicationContext(Beans(
		Bean("").Id("id").Singleton(),
	))

	// action
	bean1, e1 := ctx.GetBean("id")
	if e1 != nil {
		t.FailNow()
	}
	bean2, e2 := ctx.GetBean("id")
	if e2 != nil {
		t.FailNow()
	}

	// aasert

	s1, ok1 := bean1.(*string)
	assert.True(t, ok1)
	s2, ok2 := bean2.(*string)
	assert.True(t, ok2)

	assert.Equal(t, unsafe.Pointer(s1), unsafe.Pointer(s2))
}

func Test_applicationContext_GetBean_singleton_struct(t *testing.T) {

	// arrange
	type beanStruct struct{ s string }
	ctx, _ := ApplicationContext(Beans(
		Bean(beanStruct{}).Id("id").Singleton(),
	))

	// action
	bean1, e1 := ctx.GetBean("id")
	if e1 != nil {
		t.FailNow()
	}
	bean2, e2 := ctx.GetBean("id")
	if e2 != nil {
		t.FailNow()
	}

	// aasert

	s1, ok1 := bean1.(*beanStruct)
	assert.True(t, ok1)
	s2, ok2 := bean2.(*beanStruct)
	assert.True(t, ok2)

	assert.Equal(t, unsafe.Pointer(s1), unsafe.Pointer(s2))
}

func Test_applicationContext_GetBean_prototype_string(t *testing.T) {

	// arrange
	ctx, _ := ApplicationContext(Beans(
		Bean("").Id("id").Prototype(),
	))

	// action
	bean1, e1 := ctx.GetBean("id")
	if e1 != nil {
		t.FailNow()
	}
	bean2, e2 := ctx.GetBean("id")
	if e2 != nil {
		t.FailNow()
	}

	// aasert

	s1, ok1 := bean1.(*string)
	assert.True(t, ok1)
	s2, ok2 := bean2.(*string)
	assert.True(t, ok2)

	assert.NotEqual(t, unsafe.Pointer(s1), unsafe.Pointer(s2))
}

func Test_applicationContext_GetBean_prototype_struct(t *testing.T) {

	// arrange
	type beanStruct struct{ s string }
	ctx, _ := ApplicationContext(Beans(
		Bean(beanStruct{}).Id("id").Prototype(),
	))

	// action
	bean1, e1 := ctx.GetBean("id")
	if e1 != nil {
		t.FailNow()
	}
	bean2, e2 := ctx.GetBean("id")
	if e2 != nil {
		t.FailNow()
	}

	// aasert

	s1, ok1 := bean1.(*beanStruct)
	assert.True(t, ok1)
	s2, ok2 := bean2.(*beanStruct)
	assert.True(t, ok2)

	assert.NotEqual(t, unsafe.Pointer(s1), unsafe.Pointer(s2))
}

func Test_applicationContext_GetBean_withBeanProperty_andIsPointer(t *testing.T) {

	// arrange
	type beanStruct1 struct{ S string }
	type beanStruct2 struct{ Bean1 *beanStruct1 }
	ctx, _ := ApplicationContext(Beans(
		Bean(beanStruct2{}).Id("Bean2").
			PropertyBean("Bean1", Bean(beanStruct1{})),
	))

	// action
	bean, e := ctx.GetBean("Bean2")
	if e != nil {
		assert.FailNow(t, e.Error())
	}

	// aasert
	s, ok := bean.(*beanStruct2)
	assert.True(t, ok)

	assert.NotNil(t, s.Bean1)
}

func Test_applicationContext_GetBean_withBeanProperty_andIsNotPointer(t *testing.T) {

	// arrange
	type beanStruct1 struct{ S string }
	type beanStruct2 struct{ Bean1 beanStruct1 }
	ctx, _ := ApplicationContext(Beans(
		Bean(beanStruct2{}).Id("Bean2").
			PropertyBean("Bean1",
				Bean(beanStruct1{}).PropertyValue("S", "ss"),
			),
	))

	// action
	bean, e := ctx.GetBean("Bean2")
	if e != nil {
		assert.FailNow(t, e.Error())
	}

	// aasert
	s, ok := bean.(*beanStruct2)
	assert.True(t, ok)

	assert.Equal(t, "ss", s.Bean1.S)
}
