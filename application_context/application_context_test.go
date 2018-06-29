package application_context

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
	"github.com/yarencheng/gospring/v1"
)

func Test_New_v1Bean(t *testing.T) {
	// arrange
	config := &v1.Bean{
		ID:   "aa",
		Type: reflect.TypeOf(""),
	}

	// action
	_, err := New(config)

	// assert
	assert.NoError(t, err)
}

func Test_GetByID(t *testing.T) {
	// arrange

	datas := []struct {
		configs      []interface{}
		expectedType interface{}
	}{
		{
			configs: []interface{}{&v1.Bean{
				ID:   "a_id",
				Type: reflect.TypeOf(""),
			}},
			expectedType: new(string),
		}, {
			configs: []interface{}{&v1.Bean{
				ID:   "a_id",
				Type: reflect.TypeOf(123),
			}},
			expectedType: new(int),
		}, {
			configs: []interface{}{&v1.Bean{
				Type: reflect.TypeOf(123),
				Properties: []v1.Property{
					{
						Name: "",
						Value: &v1.Bean{
							ID:   "a_id",
							Type: reflect.TypeOf(true),
						},
					},
				},
			}},
			expectedType: new(bool),
		},
	}

	for _, data := range datas {
		ctx, err := New(data.configs...)
		require.NoError(t, err)

		// action
		bean, err := ctx.GetByID("a_id")

		// assert
		assert.NoError(t, err)
		assert.IsType(t, data.expectedType, bean)
	}
}
