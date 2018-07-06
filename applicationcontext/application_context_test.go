package applicationcontext

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
	ctx := Default()
	_, err := ctx.AddConfig(config)

	// assert
	assert.NoError(t, err)
}

func Test_GetByID(t *testing.T) {
	// arrange

	datas := []struct {
		config       interface{}
		expectedType interface{}
	}{
		{
			config: &v1.Bean{
				ID:   "a_id",
				Type: reflect.TypeOf(""),
			},
			expectedType: new(string),
		}, {
			config: &v1.Bean{
				ID:   "a_id",
				Type: reflect.TypeOf(123),
			},
			expectedType: new(int),
		}, {
			config: &v1.Bean{
				Type: reflect.TypeOf(123),
				Properties: []v1.Property{
					{
						Name: "",
						Config: &v1.Bean{
							ID:   "a_id",
							Type: reflect.TypeOf(true),
						},
					},
				},
			},
			expectedType: new(bool),
		},
	}

	for _, data := range datas {
		ctx := Default()
		_, err := ctx.AddConfig(data.config)
		require.NoError(t, err)

		// action
		bean, err := ctx.GetByID("a_id")

		// assert
		assert.NoError(t, err)
		assert.IsType(t, data.expectedType, bean)
	}
}

func Test_GetByID_withProperty(t *testing.T) {
	// arrange
	type StructInt struct {
		Value *int
	}
	type StructString struct {
		Value *string
	}
	datas := []struct {
		config interface{}
		isNil  interface{}
	}{
		{
			config: &v1.Bean{
				ID:   "a_id",
				Type: reflect.TypeOf(StructInt{}),
			},
			isNil: true,
		}, {
			config: &v1.Bean{
				ID:   "a_id",
				Type: reflect.TypeOf(StructInt{}),
				Properties: []v1.Property{
					{
						Name: "Value",
						Config: &v1.Bean{
							Type: reflect.TypeOf(123),
						},
					},
				},
			},
			isNil: false,
		}, {
			config: &v1.Bean{
				ID:   "a_id",
				Type: reflect.TypeOf(StructString{}),
			},
			isNil: true,
		}, {
			config: &v1.Bean{
				ID:   "a_id",
				Type: reflect.TypeOf(StructString{}),
				Properties: []v1.Property{
					{
						Name: "Value",
						Config: &v1.Bean{
							Type: reflect.TypeOf(""),
						},
					},
				},
			},
			isNil: false,
		},
	}

	for _, data := range datas {
		ctx := Default()
		_, err := ctx.AddConfig(data.config)
		require.NoError(t, err)

		// action
		bean, err := ctx.GetByID("a_id")

		// assert
		assert.NoError(t, err)
		actualField := reflect.ValueOf(bean).Elem().FieldByName("Value")
		assert.Equal(t, data.isNil, actualField.IsNil())
	}
}
