package list

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/satori/go.uuid"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"

	"github.com/yarencheng/gospring/v1"

	"github.com/yarencheng/gospring/mocks"
)

func Test_List_GetValue(t *testing.T) {

	type mystruct struct {
		Value int
	}

	tests := []struct {
		name     string
		expected interface{}
	}{
		{name: "int list", expected: 123},
		{name: "string list", expected: "aabb"},
		{name: "struct list", expected: mystruct{Value: 4455}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// arrange: mock element
			uuid := uuid.NewV4()
			element := new(mocks.BeanMock)
			element.On("GetUUID").Return(uuid)

			ctx := new(mocks.ApplicationContextMock)
			ctx.On("AddConfig", mock.Anything).Return(element, nil)
			ctx.On("GetByUUID", uuid).Return(test.expected, nil)

			bean, err := V1ListParser(ctx, &v1.List{
				Type: reflect.TypeOf(test.expected),
				Configs: []interface{}{
					nil,
				},
			})
			require.NoError(t, err)

			v, err := bean.GetValue()
			require.NoError(t, err)

			// assert
			assert.Equal(t, test.expected, v.Index(0).Interface())
		})
	}
}
