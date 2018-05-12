package gospring

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_TypeOf(t *testing.T) {
	// arrange
	bean := structBean{}

	// action
	bean.TypeOf("")

	// assert
	assert.Equal(t, reflect.TypeOf(""), bean.GetType())
}

func Test_structBeanGetProperty(t *testing.T) {
	// arrange
	bean := structBean{}

	// action
	actual := bean.GetProperty("")

	// assert
	assert.Nil(t, actual)
}
