package bean

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
