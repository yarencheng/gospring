package bean

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yarencheng/gospring/v1"
)

func Test_ChannelBean_GetValue(t *testing.T) {
	// arrange
	bean, err := V1ChannelParser(nil, &v1.Channel{
		Type: reflect.TypeOf(""),
	})
	require.NotNil(t, bean)
	require.NoError(t, err)

	// action
	value, err := bean.GetValue()
	require.NotNil(t, value)
	require.NoError(t, err)

	// assert
	_, ok := value.Interface().(chan string)
	require.True(t, ok)
}
