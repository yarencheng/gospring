package config2

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func Test_applicationContext_GetBean_singleton(t *testing.T) {

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

func Test_applicationContext_GetBean_singletonsss(t *testing.T) {

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
