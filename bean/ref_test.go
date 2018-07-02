package bean

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/yarencheng/gospring/v1"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"

	"github.com/yarencheng/gospring/mocks"
)

func Test_Reference_GetValue(t *testing.T) {
	for i, config := range []interface{}{
		"a_id",
		&v1.Ref{ID: "a_id"},
	} {
		fmt.Printf("i=[%v] config=[%#v]\n", i, config)
		// arrange: value for the property bean
		type ChildStruct struct {
			Value string
		}
		expected := &ChildStruct{
			Value: "hahaha",
		}

		type ParentStruct struct {
			Child ChildStruct
		}

		// arrange: Mock for the property bean
		childBean := new(mocks.BeanMock)
		childBean.On("GetValue").Return(reflect.ValueOf(expected), nil)

		// arrange: mock input source
		ctx := new(mocks.ApplicationContextMock)
		ctx.On("GetBeanByID", "a_id").Return(childBean, true)

		// arrange: create parent bean
		bean, err := V1RefParser(ctx, config)
		require.NoError(t, err)
		require.NotNil(t, bean)

		// action
		actualValue, err := bean.GetValue()
		require.NoError(t, err)

		// assert
		assert.Equal(t, expected, actualValue.Interface())
	}
}
