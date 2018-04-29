package gospring

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getDefaultFactoryFn_string(t *testing.T) {
	// arrange

	// action
	v := getDefaultFactoryFn(reflect.TypeOf(string("")))
	rvs := (*v).Call(nil)
	if !rvs[1].IsNil() {
		t.FailNow()
	}

	// assert
	actual, ok := rvs[0].Interface().(*string)
	assert.True(t, ok)
	assert.Equal(t, "", *actual)
}

func Test_getDefaultFactoryFn_int(t *testing.T) {
	// arrange

	// action
	v := getDefaultFactoryFn(reflect.TypeOf(int(999)))
	rvs := (*v).Call(nil)
	if !rvs[1].IsNil() {
		t.FailNow()
	}

	// assert
	actual, ok := rvs[0].Interface().(*int)
	assert.True(t, ok)
	assert.Equal(t, 0, *actual)
}

func Test_getDefaultFactoryFn_struct(t *testing.T) {
	// arrange
	type beanStruct struct{ S string }

	// action
	v := getDefaultFactoryFn(reflect.TypeOf(beanStruct{}))
	rvs := (*v).Call(nil)
	if !rvs[1].IsNil() {
		t.FailNow()
	}

	// assert
	actual := rvs[0].Interface().(*beanStruct)
	// assert.True(t, ok)
	assert.Equal(t, beanStruct{}, *actual)
}

func Test_bean_new_withDefaultFunction(t *testing.T) {
	// arrange
	b := Bean(int(9999))

	// action
	r, e := b.new()
	if e != nil {
		t.FailNow()
	}

	// assert
	actual, ok := r.(*int)
	assert.True(t, ok)
	assert.Equal(t, 0, *actual)
}

func Test_bean_new_withCostumeFunction(t *testing.T) {
	// arrange
	isExecute := false
	b := Bean(int(9999)).
		Factory(func() (*int, error) {
			isExecute = true
			i := int(777)
			return &i, nil
		})

	// action
	r, e := b.new()
	if e != nil {
		t.FailNow()
	}

	// assert
	assert.True(t, isExecute)
	actual, ok := r.(*int)
	assert.True(t, ok)
	assert.Equal(t, 777, *actual)
}

func Test_bean_new_withCostumeFunction_withArgv(t *testing.T) {
	// arrange
	isExecute := false
	b := Bean(int(9999)).
		Factory(func(i int) (*int, error) {
			isExecute = true
			ii := int(i + 1)
			return &ii, nil
		}, 777)

	// action
	r, e := b.new()
	if e != nil {
		t.FailNow()
	}

	// assert
	assert.True(t, isExecute)
	actual, ok := r.(*int)
	assert.True(t, ok)
	assert.Equal(t, 778, *actual)
}

func Test_bean_new_withCostumeFunction_returnError(t *testing.T) {
	// arrange
	isExecute := false
	b := Bean(int(9999)).Factory(func() (int, error) {
		isExecute = true
		return int(777), fmt.Errorf("")
	})

	// action
	_, e := b.new()

	// assert
	assert.True(t, isExecute)
	assert.NotNil(t, e)
}

func Test_bean_Singleton(t *testing.T) {
	// arrange

	// action
	b := Bean("").Singleton()

	// assert
	assert.Equal(t, scopeSingleton, b.scope)
}

func Test_bean_Prototype(t *testing.T) {
	// arrange

	// action
	b := Bean("").Prototype()

	// assert
	assert.Equal(t, scopePrototype, b.scope)
}

func Test_bean_Init(t *testing.T) {
	// arrange
	isCall := false
	b := Bean("").Init(func(s *string) {
		isCall = true
	})

	// action
	s := ""
	argv := []reflect.Value{
		reflect.ValueOf(&s),
	}
	b.initFn.Call(argv)

	// assert
	assert.True(t, isCall)
}

func Test_bean_Finalize(t *testing.T) {
	// arrange
	isCall := false
	b := Bean("").Finalize(func(s *string) {
		isCall = true
	})

	// action
	s := ""
	argv := []reflect.Value{
		reflect.ValueOf(&s),
	}
	b.finalizeFn.Call(argv)

	// assert
	assert.True(t, isCall)
}
