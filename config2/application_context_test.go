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

func Test_applicationContext_GetBean_withRefbeanProperty_andIsPointer(t *testing.T) {

	// arrange
	type beanStruct1 struct{ S string }
	type beanStruct2 struct{ Bean1 *beanStruct1 }
	ctx, _ := ApplicationContext(Beans(
		Bean(beanStruct1{}).Id("Bean1").PropertyValue("S", "ss"),
		Bean(beanStruct2{}).Id("Bean2").
			PropertyRef("Bean1", "Bean1"),
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

func Test_applicationContext_GetBean_withRefbeanProperty_andIsNotPointer(t *testing.T) {

	// arrange
	type beanStruct1 struct{ S string }
	type beanStruct2 struct{ Bean1 beanStruct1 }
	ctx, _ := ApplicationContext(Beans(
		Bean(beanStruct1{}).Id("Bean1").PropertyValue("S", "ss"),
		Bean(beanStruct2{}).Id("Bean2").
			PropertyRef("Bean1", "Bean1"),
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

func Test_applicationContext_GetBean_withValueProperty_andIsPointer(t *testing.T) {

	// arrange
	type beanStruct struct{ S *string }
	ctx, _ := ApplicationContext(Beans(
		Bean(beanStruct{}).Id("Bean1").PropertyValue("S", "ss"),
	))

	// action
	bean, e := ctx.GetBean("Bean1")
	if e != nil {
		assert.FailNow(t, e.Error())
	}

	// aasert
	s, ok := bean.(*beanStruct)
	assert.True(t, ok)

	assert.Equal(t, "ss", *s.S)
}

func Test_applicationContext_GetBean_withValueProperty_andIsNotPointer(t *testing.T) {

	// arrange
	type beanStruct struct{ S string }
	ctx, _ := ApplicationContext(Beans(
		Bean(beanStruct{}).Id("Bean1").PropertyValue("S", "ss"),
	))

	// action
	bean, e := ctx.GetBean("Bean1")
	if e != nil {
		assert.FailNow(t, e.Error())
	}

	// aasert
	s, ok := bean.(*beanStruct)
	assert.True(t, ok)

	assert.Equal(t, "ss", s.S)
}

func Test_applicationContext_demo(t *testing.T) {

	// arrange

	type Bean2 struct {
		Value string
	}
	type Bean1 struct {
		Ivalue   int
		Svalue   string
		IPointer *int
		SPointer *string
		Bvalue   Bean2
		BPointer *Bean2
		Blocal   *Bean2
	}
	ctx, _ := ApplicationContext(Beans(
		Bean(Bean1{}).Id("id_b"),
		Bean(Bean1{}).Id("id_c").Prototype(),
		Bean(Bean1{}).Id("id_d").Singleton(),
		Bean(Bean1{}).Id("id_e").PropertyValue("Ivalue", 123),
		Bean(Bean1{}).Id("id_f").PropertyValue("Svalue", "hi spring"),
		Bean(Bean1{}).Id("id_g").PropertyValue("IPointer", 123),
		Bean(Bean1{}).Id("id_h").PropertyValue("SPointer", "hi spring"),
		Bean(Bean1{}).Id("id_i").PropertyRef("Bvalue", "id_a"),
		Bean(Bean1{}).Id("id_j").PropertyRef("BPointer", "id_a"),
		Bean(Bean1{}).Id("id_k").PropertyBean(
			"Blocal", Bean(Bean2{}).PropertyValue("Value", "is a local bean without ID"),
		),
		Bean(Bean2{}).Id("id").PropertyValue("Value", "bean2 foo"),
		Bean(Bean2{}).Id("id_a"),
	))

	var bean interface{}
	var beanA interface{}
	var beanB interface{}
	var i int
	var s string

	bean, _ = ctx.GetBean("id")
	assert.Equal(t, Bean2{Value: "bean2 foo"}, *bean.(*Bean2))

	bean, _ = ctx.GetBean("id_a")
	assert.Equal(t, Bean2{}, *bean.(*Bean2))

	bean, _ = ctx.GetBean("id_b")
	assert.Equal(t, Bean1{}, *bean.(*Bean1))

	beanA, _ = ctx.GetBean("id_c")
	beanB, _ = ctx.GetBean("id_c")
	assert.NotEqual(t, unsafe.Pointer(beanA.(*Bean1)), unsafe.Pointer(beanB.(*Bean1)))

	beanA, _ = ctx.GetBean("id_d")
	beanB, _ = ctx.GetBean("id_d")
	assert.Equal(t, unsafe.Pointer(beanA.(*Bean1)), unsafe.Pointer(beanB.(*Bean1)))

	bean, _ = ctx.GetBean("id_e")
	assert.Equal(t, Bean1{Ivalue: 123}, *bean.(*Bean1))

	bean, _ = ctx.GetBean("id_f")
	assert.Equal(t, Bean1{Svalue: "hi spring"}, *bean.(*Bean1))

	bean, _ = ctx.GetBean("id_g")
	i = 123
	assert.Equal(t, Bean1{IPointer: &i}, *bean.(*Bean1))

	bean, _ = ctx.GetBean("id_h")
	s = "hi spring"
	assert.Equal(t, Bean1{SPointer: &s}, *bean.(*Bean1))

	bean, _ = ctx.GetBean("id_i")
	assert.Equal(t, Bean1{Bvalue: Bean2{}}, *bean.(*Bean1))

	bean, _ = ctx.GetBean("id_j")
	assert.Equal(t, Bean1{BPointer: &Bean2{}}, *bean.(*Bean1))

	bean, _ = ctx.GetBean("id_k")
	assert.Equal(t, Bean1{Blocal: &Bean2{
		Value: "is a local bean without ID",
	}}, *bean.(*Bean1))
}
