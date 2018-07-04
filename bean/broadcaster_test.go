package bean

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yarencheng/gospring/mocks"
	"github.com/yarencheng/gospring/v1"
)

func Test_Broadcaster_GetValue_intChannel(t *testing.T) {
	// arrange: mock input source
	sourceCh := make(chan int)
	ctx := new(mocks.ApplicationContextMock)
	ctx.On("GetByID", "from_id").Return(sourceCh, nil)

	// arrange: create broad cast bean
	bean, err := V1BroadcastParser(ctx, &v1.Broadcast{
		SourceID: "from_id",
		Size:     1,
	})
	require.NotNil(t, bean)
	require.NoError(t, err)

	// action: get 2 output channel
	value1, err := bean.GetValue()
	require.True(t, value1.IsValid())
	require.NoError(t, err)
	value2, err := bean.GetValue()
	require.True(t, value2.IsValid())
	require.NoError(t, err)
	// action: send data to input channel
	sourceCh <- 123

	// assert: both channel receive same data
	i1, ok := value1.Recv()
	assert.True(t, ok)
	assert.Equal(t, 123, i1.Interface().(int))
	i2, ok := value2.Recv()
	assert.True(t, ok)
	assert.Equal(t, 123, i2.Interface().(int))
}

func Test_Broadcaster_GetValue_structChannel(t *testing.T) {
	type TestStruct struct {
		I int
	}

	// arrange: mock input source
	sourceCh := make(chan TestStruct)
	ctx := new(mocks.ApplicationContextMock)
	ctx.On("GetByID", "from_id").Return(sourceCh, nil)

	// arrange: create broad cast bean
	bean, err := V1BroadcastParser(ctx, &v1.Broadcast{
		SourceID: "from_id",
		Size:     1,
	})
	require.NotNil(t, bean)
	require.NoError(t, err)

	// action: get 2 output channel
	value1, err := bean.GetValue()
	require.True(t, value1.IsValid())
	require.NoError(t, err)
	value2, err := bean.GetValue()
	require.True(t, value2.IsValid())
	require.NoError(t, err)
	// action: send data to input channel
	sourceCh <- TestStruct{I: 123}

	// assert: both channel receive same data
	i1, ok := value1.Recv()
	assert.True(t, ok)
	assert.Equal(t, TestStruct{I: 123}, i1.Interface().(TestStruct))
	i2, ok := value2.Recv()
	assert.True(t, ok)
	assert.Equal(t, TestStruct{I: 123}, i2.Interface().(TestStruct))
}

func Test_Broadcaster_GetValue_interfaceChannel(t *testing.T) {
	type TestStruct struct {
		I int
	}

	// arrange: mock input source
	sourceCh := make(chan interface{})
	ctx := new(mocks.ApplicationContextMock)
	ctx.On("GetByID", "from_id").Return(sourceCh, nil)

	// arrange: create broad cast bean
	bean, err := V1BroadcastParser(ctx, &v1.Broadcast{
		SourceID: "from_id",
		Size:     1,
	})
	require.NotNil(t, bean)
	require.NoError(t, err)

	// action: get 2 output channel
	value1, err := bean.GetValue()
	require.True(t, value1.IsValid())
	require.NoError(t, err)
	value2, err := bean.GetValue()
	require.True(t, value2.IsValid())
	require.NoError(t, err)
	// action: send data to input channel
	sourceCh <- TestStruct{I: 123}

	// assert: both channel receive same data
	i1, ok := value1.Recv()
	assert.True(t, ok)
	assert.Equal(t, TestStruct{I: 123}, i1.Interface().(TestStruct))
	i2, ok := value2.Recv()
	assert.True(t, ok)
	assert.Equal(t, TestStruct{I: 123}, i2.Interface().(TestStruct))
}
