package bean

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yarencheng/gospring/v1"
)

func Test_NewStructBeanV1_checkDefaultValue(t *testing.T) {
	// action
	bean, err := NewStructBeanV1(v1.Bean{Type: reflect.TypeOf("")})
	require.NoError(t, err)

	// assert
	assert.Equal(t, v1.Default, bean.scope)
}

func Test_GetValue(t *testing.T) {
	// arrange
	type testStruct struct{}
	bean := &StructBean{
		tvpe: reflect.TypeOf(testStruct{}),
	}

	// action
	v, err := bean.GetValue()

	// assert
	assert.Exactly(t, v.Interface(), &testStruct{})
	assert.NoError(t, err)
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
