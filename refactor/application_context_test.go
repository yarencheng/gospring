package refactor

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

type ValueBeanMock struct {
	mock.Mock
}

func (m *ValueBeanMock) GetID() *string {
	args := m.Called()
	s := args.Get(0)
	if s == nil {
		return nil
	} else {
		return s.(*string)
	}
}

func (m *ValueBeanMock) GetValue() interface{} {
	args := m.Called()
	return args.Get(0)
}

type ReferenceBeanMock struct {
	mock.Mock
}

func (m *ReferenceBeanMock) GetID() *string {
	args := m.Called()
	s := args.Get(0)
	if s == nil {
		return nil
	} else {
		return s.(*string)
	}
}

func (m *ReferenceBeanMock) ID(id string) ReferenceBeanI {
	args := m.Called(id)
	return args.Get(0).(ReferenceBeanI)
}

func Test_NewApplicationContext_empty(t *testing.T) {
	// arrange

	// action
	_, e := NewApplicationContext()

	// assert
	assert.Nil(t, e)
}

func Test_NewApplicationContext_ValueBeanI_withoutId(t *testing.T) {
	// arrange
	mock := new(ValueBeanMock)
	mock.On("GetID").Return(nil)

	beans := []BeanI{mock}

	// action
	ctx, e := NewApplicationContext(beans...)
	require.Nil(t, e)

	// assert
	assert.NotContains(t, ctx.(*applicationContext).beanById, mock)
}

func Test_NewApplicationContext_ValueBeanI_withId(t *testing.T) {
	// arrange
	id := "id"

	mock := new(ValueBeanMock)
	mock.On("GetID").Return(&id)

	beans := []BeanI{mock}

	// action
	ctx, e := NewApplicationContext(beans...)
	require.Nil(t, e)

	// assert
	assert.Contains(t, ctx.(*applicationContext).beanById, id)
}

func Test_NewApplicationContext_ValueBeanI_withDuplicatedId(t *testing.T) {
	// arrange
	id := "id"

	mock1 := new(ValueBeanMock)
	mock1.On("GetID").Return(&id)
	mock2 := new(ValueBeanMock)
	mock2.On("GetID").Return(&id)

	beans := []BeanI{mock1, mock2}

	// action
	_, e := NewApplicationContext(beans...)
	require.NotNil(t, e)

	// assert
}

func Test_NewApplicationContext_ReferenceBeanI(t *testing.T) {
	// arrange
	mock := new(ReferenceBeanMock)
	mock.On("GetID").Return(nil)
	beans := []BeanI{mock}

	// action
	_, e := NewApplicationContext(beans...)

	// assert
	assert.Nil(t, e)
}

func Test_NewApplicationContext_StructBeanI_withId(t *testing.T) {
	// arrange
	id := "id"
	type beanStract struct{}
	beans := Beans(
		Bean(beanStract{}).ID(id),
	)

	// action
	_, e := NewApplicationContext(beans...)
	ctx, e := NewApplicationContext(beans...)
	require.Nil(t, e)

	// assert
	assert.Contains(t, ctx.(*applicationContext).beanById, id)
}

func Test_NewApplicationContext_StructBeanI_withDuplicatedId(t *testing.T) {
	// arrange
	type beanStract struct{}
	beans := Beans(
		Bean(beanStract{}).ID("id"),
		Bean(beanStract{}).ID("id"),
	)

	// action
	_, e := NewApplicationContext(beans...)

	// assert
	assert.NotNil(t, e)
}

func Test_NewApplicationContext_StructBeanI_checkFactory_notFunction(t *testing.T) {
	// arrange
	type beanStract struct{}
	beans := Beans(
		Bean(beanStract{}).
			Factory(""),
	)

	// action
	_, e := NewApplicationContext(beans...)

	// assert
	assert.NotNil(t, e)
}

func Test_NewApplicationContext_StructBeanI_checkFactory_lengthOfArgvMismatch_1(t *testing.T) {
	// arrange
	type beanStract struct{}
	beans := Beans(
		Bean(&beanStract{}).
			Factory(func(i int) *beanStract {
				return nil
			}),
	)

	// action
	_, e := NewApplicationContext(beans...)

	// assert
	assert.NotNil(t, e)
}

func Test_NewApplicationContext_StructBeanI_checkFactory_lengthOfArgvMismatch_2(t *testing.T) {
	// arrange
	type beanStract struct{}
	beans := Beans(
		Bean(&beanStract{}).
			Factory(func() *beanStract {
				return nil
			}, 123),
	)

	// action
	_, e := NewApplicationContext(beans...)

	// assert
	assert.NotNil(t, e)
}

func Test_NewApplicationContext_StructBeanI_checkFactory_returnOneValue(t *testing.T) {
	// arrange
	type beanStract struct{}
	beans := Beans(
		Bean(&beanStract{}).
			Factory(func() string {
				return "return something else"
			}),
	)

	// action
	_, e := NewApplicationContext(beans...)

	// assert
	assert.NotNil(t, e)
}

func Test_NewApplicationContext_StructBeanI_checkFactory_returnTwoValue_1(t *testing.T) {
	// arrange
	type beanStract struct{}
	beans := Beans(
		Bean(&beanStract{}).
			Factory(func() (string, error) {
				return "return something else", nil
			}),
	)

	// action
	_, e := NewApplicationContext(beans...)

	// assert
	assert.NotNil(t, e)
}

func Test_NewApplicationContext_StructBeanI_checkFactory_returnTwoValue_2(t *testing.T) {
	// arrange
	type beanStract struct{}
	beans := Beans(
		Bean(&beanStract{}).
			Factory(func() (*beanStract, string) {
				return nil, "return something else"
			}),
	)

	// action
	_, e := NewApplicationContext(beans...)

	// assert
	assert.NotNil(t, e)
}

func Test_NewApplicationContext_StructBeanI_checkFactory_returnMoreThanTwoValue(t *testing.T) {
	// arrange
	type beanStract struct{}
	beans := Beans(
		Bean(&beanStract{}).
			Factory(func() (*beanStract, error, string) {
				return nil, nil, "return something more"
			}),
	)

	// action
	_, e := NewApplicationContext(beans...)

	// assert
	assert.NotNil(t, e)
}

func Test_NewApplicationContext_StructBeanI_prototypeWithFinalizer(t *testing.T) {
	// arrange
	type beanStract struct{}
	beans := Beans(
		Bean(beanStract{}).
			Prototype().
			Finalize("aa"),
	)

	// action
	_, e := NewApplicationContext(beans...)

	// assert
	assert.NotNil(t, e)
}

func Test_NewApplicationContext_StructBeanI_dependencyLoop_1(t *testing.T) {
	// arrange
	type beanStract struct{}
	beans := Beans(
		Bean(beanStract{}).ID("id_1").
			Property("aa",
				Bean(beanStract{}).
					ID("id_2").
					Property("aa",
						Ref("id_1"),
					),
			),
	)

	// action
	_, e := NewApplicationContext(beans...)

	// assert
	assert.NotNil(t, e)
}

func Test_NewApplicationContext_StructBeanI_dependencyLoop_2(t *testing.T) {
	// arrange
	type beanStract struct{}
	beans := Beans(
		Bean(beanStract{}).ID("id_1").
			Property("aa",
				Bean(beanStract{}).
					Property("aa",
						Ref("id_1"),
					),
			),
	)

	// action
	_, e := NewApplicationContext(beans...)

	// assert
	assert.NotNil(t, e)
}
