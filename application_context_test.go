package gospring

import (
	"fmt"
	"runtime"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func Test_applicationContext_GetBean_withValueProperty_andIsSliceOfString(t *testing.T) {

	// arrange
	type beanStruct struct{ S []string }
	ctx, _ := ApplicationContext(Beans(
		Bean(beanStruct{}).Id("Bean1").PropertyValue("S", "AAA", "BBB"),
	))

	// action
	bean, e := ctx.GetBean("Bean1")
	if e != nil {
		assert.FailNow(t, e.Error())
	}

	// aasert
	s, ok := bean.(*beanStruct)
	assert.True(t, ok)

	assert.Equal(t, beanStruct{
		S: []string{"AAA", "BBB"},
	}, *s)
}

func Test_applicationContext_GetBean_withValueProperty_andIsSliceOfStringPointer(t *testing.T) {

	// arrange
	type beanStruct struct{ S []*string }
	ctx, _ := ApplicationContext(Beans(
		Bean(beanStruct{}).Id("Bean1").PropertyValue("S", "AAA", "BBB"),
	))

	// action
	bean, e := ctx.GetBean("Bean1")
	if e != nil {
		assert.FailNow(t, e.Error())
	}

	// aasert
	s, ok := bean.(*beanStruct)
	assert.True(t, ok)

	aaa := "AAA"
	bbb := "BBB"
	assert.Equal(t, beanStruct{
		S: []*string{&aaa, &bbb},
	}, *s)
}

type Test_applicationContext_GetBean_withCostumeInitFn_struct struct{ S string }

var Test_applicationContext_GetBean_withCostumeInitFn_isRun = false

func (*Test_applicationContext_GetBean_withCostumeInitFn_struct) Test_applicationContext_GetBean_withCostumeInitFn_fn() {
	Test_applicationContext_GetBean_withCostumeInitFn_isRun = true
}

func Test_applicationContext_GetBean_withCostumeInitFn(t *testing.T) {

	// arrange
	ctx, _ := ApplicationContext(Beans(
		Bean(Test_applicationContext_GetBean_withCostumeInitFn_struct{}).Id("Bean1").Init("Test_applicationContext_GetBean_withCostumeInitFn_fn"),
	))

	// action
	_, e := ctx.GetBean("Bean1")
	require.Nil(t, e)

	// aasert
	assert.True(t, Test_applicationContext_GetBean_withCostumeInitFn_isRun)
}

type Test_applicationContext_GetBean_withCostumeInitFn_andReturnError_struct struct{ S string }

var Test_applicationContext_GetBean_withCostumeInitFn_andReturnError_isRun = false

func (*Test_applicationContext_GetBean_withCostumeInitFn_andReturnError_struct) Test_applicationContext_GetBean_withCostumeInitFn_andReturnError_fn() error {
	Test_applicationContext_GetBean_withCostumeInitFn_andReturnError_isRun = true
	return fmt.Errorf("")
}

func Test_applicationContext_GetBean_withCostumeInitFn_andReturnError(t *testing.T) {

	// arrange
	ctx, _ := ApplicationContext(Beans(
		Bean(Test_applicationContext_GetBean_withCostumeInitFn_andReturnError_struct{}).Id("Bean1").Init("Test_applicationContext_GetBean_withCostumeInitFn_andReturnError_fn"),
	))

	// action
	_, e := ctx.GetBean("Bean1")
	require.NotNil(t, e)

	// aasert
	assert.True(t, Test_applicationContext_GetBean_withCostumeInitFn_andReturnError_isRun)
}

type Test_applicationContext_GetBean_withStructInitFn_struct struct{ S string }

var Test_applicationContext_GetBean_withStructInitFn_isRun = false

func (b *Test_applicationContext_GetBean_withStructInitFn_struct) Init() {
	Test_applicationContext_GetBean_withStructInitFn_isRun = true
}

func Test_applicationContext_GetBean_withStructInitFn(t *testing.T) {

	// arrange
	ctx, _ := ApplicationContext(Beans(
		Bean(Test_applicationContext_GetBean_withStructInitFn_struct{}).Id("Bean1"),
	))

	// action
	_, e := ctx.GetBean("Bean1")
	require.Nil(t, e)

	// aasert
	assert.True(t, Test_applicationContext_GetBean_withStructInitFn_isRun)
}

type Test_applicationContext_GetBean_withCostumeFinalizeFn_struct struct{ S string }

var Test_applicationContext_GetBean_withCostumeFinalizeFn_isRun = false

func (b *Test_applicationContext_GetBean_withCostumeFinalizeFn_struct) Aaa() {
	Test_applicationContext_GetBean_withCostumeFinalizeFn_isRun = true
}

func Test_applicationContext_GetBean_withCostumeFinalizeFn(t *testing.T) {

	// arrange
	ctx, _ := ApplicationContext(Beans(
		Bean(Test_applicationContext_GetBean_withCostumeFinalizeFn_struct{}).Id("Bean1").Finalize("Aaa"),
	))

	// action
	_, e := ctx.GetBean("Bean1")
	require.Nil(t, e)
	ctx.Finalize()

	// aasert
	assert.True(t, Test_applicationContext_GetBean_withCostumeFinalizeFn_isRun)
}

type beanStruct_g2 struct{ S string }

func (b *beanStruct_g2) Finalize() {
	Test_applicationContext_GetBean_withStructFinalizeFn_isCalled = true
}

var Test_applicationContext_GetBean_withStructFinalizeFn_isCalled = false

func Test_applicationContext_GetBean_withStructFinalizeFn(t *testing.T) {

	// arrange
	ctx, _ := ApplicationContext(Beans(
		Bean(beanStruct_g2{}).Id("Bean1"),
	))

	// action
	_, e := ctx.GetBean("Bean1")
	require.Nil(t, e)
	ctx.Finalize()

	// aasert
	runtime.GC()
	assert.True(t, Test_applicationContext_GetBean_withStructFinalizeFn_isCalled)
}
