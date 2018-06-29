package bean

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yarencheng/gospring/v1"
)

func Test_NewStructBeanV1_checkDefaultValue(t *testing.T) {
	// action
	bean, err := NewStructBeanV1(nil, &v1.Bean{Type: reflect.TypeOf("")})
	require.NoError(t, err)

	// assert
	assert.Equal(t, v1.Default, bean.scope)
}

func Test_GetValue(t *testing.T) {
	// arrange
	type testStruct struct{}
	bean := &StructBean{
		tvpe:  reflect.TypeOf(testStruct{}),
		scope: v1.Default,
	}

	// action
	v, err := bean.GetValue()
	require.NoError(t, err)

	// assert
	assert.Exactly(t, v.Interface(), &testStruct{})
}

func Test_GetValue_prototype(t *testing.T) {
	// arrange
	type testStruct struct{ i int }
	bean := &StructBean{
		tvpe:  reflect.TypeOf(testStruct{}),
		scope: v1.Prototype,
	}

	// action
	v1, err := bean.GetValue()
	require.NoError(t, err)
	v2, err := bean.GetValue()
	require.NoError(t, err)

	// assert
	assert.NotEqual(t, v1.Pointer(), v2.Pointer())
}

func Test_GetValue_singleton(t *testing.T) {
	// arrange
	type testStruct struct{ i int }
	bean := &StructBean{
		tvpe:  reflect.TypeOf(testStruct{}),
		scope: v1.Singleton,
	}

	// action
	v1, err := bean.GetValue()
	require.NoError(t, err)
	v2, err := bean.GetValue()
	require.NoError(t, err)

	// assert
	assert.Equal(t, v1.Pointer(), v2.Pointer())
}

func Test_GetValue_defaultScope(t *testing.T) {
	// arrange
	type testStruct struct{ i int }
	bean := &StructBean{
		tvpe:  reflect.TypeOf(testStruct{}),
		scope: v1.Default,
	}

	// action
	v1, err := bean.GetValue()
	require.NoError(t, err)
	v2, err := bean.GetValue()
	require.NoError(t, err)

	// assert
	assert.Equal(t, v1.Pointer(), v2.Pointer())
}

func Test_GetValue_fromFactory(t *testing.T) {
	// arrange
	type testStruct struct{ i int }
	expected := &testStruct{i: 123}
	config := &v1.Bean{
		Type: reflect.TypeOf(testStruct{}),
		FactoryFn: func() *testStruct {
			return expected
		},
	}
	bean, err := NewStructBeanV1(nil, config)
	require.NoError(t, err)

	// action
	v, err := bean.GetValue()
	require.NoError(t, err)

	// assert
	actual, ok := v.Interface().(*testStruct)
	require.True(t, ok)
	assert.Equal(t, expected, actual)
}

type Test_GetValue_withDefaultStartFn_struct struct {
	i int
}

func (s *Test_GetValue_withDefaultStartFn_struct) Start() {
	s.i = 123
}
func Test_GetValue_withDefaultStartFn(t *testing.T) {
	// arrange
	expected := &Test_GetValue_withDefaultStartFn_struct{i: 123}
	config := &v1.Bean{
		Type: reflect.TypeOf(Test_GetValue_withDefaultStartFn_struct{}),
	}
	bean, err := NewStructBeanV1(nil, config)
	require.NoError(t, err)

	// action
	v, err := bean.GetValue()
	require.NoError(t, err)

	// assert
	actual, ok := v.Interface().(*Test_GetValue_withDefaultStartFn_struct)
	require.True(t, ok, "%v", v.Type())
	assert.Equal(t, expected, actual)
}

type Test_GetValue_withDefaultStopFn_struct struct {
	i int
}

func (s *Test_GetValue_withDefaultStopFn_struct) Stop() {
	s.i = 123
}
func Test_GetValue_withDefaultStopFn(t *testing.T) {
	// arrange
	expected := &Test_GetValue_withDefaultStopFn_struct{i: 123}
	config := &v1.Bean{
		Type:  reflect.TypeOf(Test_GetValue_withDefaultStopFn_struct{}),
		Scope: v1.Singleton,
	}
	bean, err := NewStructBeanV1(nil, config)
	require.NoError(t, err)

	// arrange
	v, err := bean.GetValue()
	require.NoError(t, err)

	// action
	err = bean.Stop(context.Background())
	require.NoError(t, err)

	// assert
	actual, ok := v.Interface().(*Test_GetValue_withDefaultStopFn_struct)
	require.True(t, ok, "%v", v.Type())
	assert.Equal(t, expected, actual)
}
