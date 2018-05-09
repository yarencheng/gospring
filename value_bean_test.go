package gospring

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func Test_GetID(t *testing.T) {
	// arrange
	bean := valueBean{}

	// action
	id := bean.GetID()

	// assert
	assert.Nil(t, id)
}

func Test_GetValue(t *testing.T) {
	// arrange
	expected := 123
	bean := valueBean{
		value: expected,
	}

	// action
	value := bean.GetValue()

	// assert
	actual, ok := value.(*int)
	require.True(t, ok)
	assert.Equal(t, expected, *actual)
}

func Test_GetScope(t *testing.T) {
	// arrange
	bean := valueBean{}

	// action
	actual := bean.GetScope()

	// assert
	assert.Equal(t, Prototype, actual)
}

func Test_GetFactory(t *testing.T) {
	// arrange
	bean := valueBean{}

	// action
	_, actualArgv := bean.GetFactory()

	// reflect.ValueOf(actualFn).

	// assert
	// TODO: compare returned function. there seems no way to do it now
	assert.Equal(t, []BeanI{}, actualArgv)
}

func Test_GetFinalize(t *testing.T) {
	// arrange
	bean := valueBean{}

	// action
	actual := bean.GetFinalize()

	// assert
	assert.Nil(t, actual)
}

func Test_GetInit(t *testing.T) {
	// arrange
	bean := valueBean{}

	// action
	actual := bean.GetInit()

	// assert
	assert.Nil(t, actual)
}

func Test_GetProperty(t *testing.T) {
	// arrange
	bean := valueBean{}

	// action
	actual := bean.GetProperty("")

	// assert
	assert.Nil(t, actual)
}

func Test_GetProperties(t *testing.T) {
	// arrange
	bean := valueBean{}

	// action
	actual := bean.GetProperties()

	// assert
	assert.Empty(t, actual)
}

func Test_GetType(t *testing.T) {
	// arrange
	bean := valueBean{
		value: 123,
	}

	// action
	actual := bean.GetType()

	// assert
	assert.Equal(t, reflect.TypeOf(123), actual)
}
