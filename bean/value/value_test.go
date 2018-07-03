package value

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
	"github.com/yarencheng/gospring/v1"
)

func Test_GetValue(t *testing.T) {
	datas := []struct {
		value interface{}
	}{
		{value: "aabbcc"},
		{value: 123},
		{value: map[string]string{
			"kkk": "vvv",
		}},
	}
	for _, data := range datas {
		// arrange
		bean, err := V1ValueParser(nil, &v1.Value{
			Value: data.value,
		})
		require.NotNil(t, bean)
		require.NoError(t, err)

		// action
		value, err := bean.GetValue()
		require.NotNil(t, value)
		require.NoError(t, err)

		// assert
		assert.Equal(t, data.value, value.Interface())
	}
}
