package refactor

import (
	"fmt"
	"reflect"
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
	return m.Called().Get(0)
}

func (m *ValueBeanMock) GetScope() Scope {
	return m.Called().Get(0).(Scope)
}

func (m *ValueBeanMock) GetFactory() (interface{}, []BeanI) {
	args := m.Called()
	return args.Get(0), args.Get(1).([]BeanI)
}

func (m *ValueBeanMock) GetFinalize() *string {
	return m.Called().Get(0).(*string)
}

func (m *ValueBeanMock) GetInit() *string {
	return m.Called().Get(0).(*string)
}

func (m *ValueBeanMock) GetProperty(name string) []BeanI {
	return m.Called().Get(0).([]BeanI)
}

func (m *ValueBeanMock) GetProperties() map[string][]BeanI {
	return m.Called().Get(0).(map[string][]BeanI)
}

func (m *ValueBeanMock) GetType() reflect.Type {
	args := m.Called()
	s := args.Get(0)
	if s == nil {
		return nil
	} else {
		return s.(reflect.Type)
	}
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

func (m *ReferenceBeanMock) GetScope() Scope {
	return m.Called().Get(0).(Scope)
}

func (m *ReferenceBeanMock) GetFactory() (interface{}, []BeanI) {
	args := m.Called()
	return args.Get(0), args.Get(1).([]BeanI)
}

func (m *ReferenceBeanMock) GetFinalize() *string {
	return m.Called().Get(0).(*string)
}

func (m *ReferenceBeanMock) GetInit() *string {
	return m.Called().Get(0).(*string)
}

func (m *ReferenceBeanMock) GetProperty(name string) []BeanI {
	return m.Called().Get(0).([]BeanI)
}

func (m *ReferenceBeanMock) GetProperties() map[string][]BeanI {
	return m.Called().Get(0).(map[string][]BeanI)
}

func (m *ReferenceBeanMock) GetType() reflect.Type {
	return m.Called().Get(0).(reflect.Type)
}

func (m *ReferenceBeanMock) GetReference() BeanI {
	return m.Called().Get(0).(BeanI)
}

func (m *ReferenceBeanMock) SetReference(b BeanI) {
	m.Called(b)
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
	mock.On("GetFactory").Return(nil, []BeanI{})
	mock.On("GetScope").Return(Singleton)
	mock.On("GetProperties").Return(map[string][]BeanI{})
	mock.On("GetType").Return(nil)

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
	mock.On("GetFactory").Return(nil, []BeanI{})
	mock.On("GetScope").Return(Singleton)
	mock.On("GetProperties").Return(map[string][]BeanI{})
	mock.On("GetType").Return(nil)

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
	mock1.On("GetFactory").Return(nil, []BeanI{})
	mock1.On("GetScope").Return(Singleton)
	mock1.On("GetProperties").Return(map[string][]BeanI{})
	mock1.On("GetType").Return(nil)

	mock2 := new(ValueBeanMock)
	mock2.On("GetID").Return(&id)
	mock2.On("GetFactory").Return(nil, []BeanI{})
	mock2.On("GetScope").Return(Singleton)
	mock2.On("GetProperties").Return(map[string][]BeanI{})
	mock1.On("GetType").Return(nil)

	beans := []BeanI{mock1, mock2}

	// action
	_, e := NewApplicationContext(beans...)
	require.NotNil(t, e)

	// assert
}

func Test_NewApplicationContext_ReferenceBeanI(t *testing.T) {
	// arrange
	id := "id"
	beans := Beans(
		Bean("").ID(id),
		Ref(id),
	)

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

func Test_NewApplicationContext_typeOfbeanCantBePointer(t *testing.T) {
	// arrange
	type beanStract struct{}
	beans := Beans(
		Bean(&beanStract{}),
	)

	// action
	_, e := NewApplicationContext(beans...)

	// assert
	assert.NotNil(t, e)
}

func Test_applicationContext_GetBean_idExist(t *testing.T) {
	// arrange
	type beanStract struct{}
	beans := Beans(
		Bean(beanStract{}).ID("id_1"),
	)

	// action
	ctx, ctxe := NewApplicationContext(beans...)
	require.Nil(t, ctxe)
	bean, beane := ctx.GetBean("id_1")
	require.Nil(t, beane)

	// assert
	assert.NotNil(t, bean)
}

func Test_applicationContext_GetBean_idNotExist(t *testing.T) {
	// arrange
	type beanStract struct{}
	beans := Beans(
		Bean(beanStract{}).ID("id_1"),
	)

	// action
	ctx, ctxe := NewApplicationContext(beans...)
	require.Nil(t, ctxe)
	_, beane := ctx.GetBean("id_2")

	// assert
	require.NotNil(t, beane)
}

func Test_applicationContext_GetBean_fromFactory(t *testing.T) {
	// arrange
	expected := "content text"
	beans := Beans(
		Bean(string("")).ID("id_1").Factory(func() *string {
			return &expected
		}),
	)

	// action
	ctx, ctxe := NewApplicationContext(beans...)
	require.Nil(t, ctxe)
	actual, beane := ctx.GetBean("id_1")
	require.Nil(t, beane)

	// assert
	assert.Equal(t, expected, *actual.(*string))
}

func Test_applicationContext_GetBean_fromFactory_withParameter(t *testing.T) {
	// arrange
	beans := Beans(
		Bean(string("")).ID("id_1").Factory(func(in *string) *string {
			s := "Hi " + *in
			return &s
		}, "gospring"),
	)

	// action
	ctx, ctxe := NewApplicationContext(beans...)
	require.Nil(t, ctxe)
	actual, beane := ctx.GetBean("id_1")
	require.Nil(t, beane)

	// assert
	assert.Equal(t, "Hi gospring", *actual.(*string))
}

func Test_applicationContext_GetBean_fromFactory_withBeanParameter(t *testing.T) {
	// arrange
	beans := Beans(
		Bean(string("")).ID("id_1").Factory(func(in *string) *string {
			s := "Hi " + *in
			return &s
		}, Ref("id_2")),
		Bean(string("")).ID("id_2").Factory(func() *string {
			s := "id_2"
			return &s
		}),
	)

	// action
	ctx, ctxe := NewApplicationContext(beans...)
	require.Nil(t, ctxe)
	actual, beane := ctx.GetBean("id_1")
	require.Nil(t, beane)

	// assert
	assert.Equal(t, "Hi id_2", *actual.(*string))
}

func Test_applicationContext_GetBean_fromFactory_returnError_1(t *testing.T) {
	// arrange
	beans := Beans(
		Bean(string("")).ID("id_1").Factory(func() interface{} {
			return fmt.Errorf("")
		}),
	)

	// action
	ctx, ctxe := NewApplicationContext(beans...)
	require.Nil(t, ctxe)
	_, beane := ctx.GetBean("id_1")

	// assert
	assert.NotNil(t, beane)
}

func Test_applicationContext_GetBean_fromFactory_returnError_2(t *testing.T) {
	// arrange
	beans := Beans(
		Bean(string("")).ID("id_1").Factory(func() (*string, error) {
			return nil, fmt.Errorf("")
		}),
	)

	// action
	ctx, ctxe := NewApplicationContext(beans...)
	require.Nil(t, ctxe)
	_, beane := ctx.GetBean("id_1")

	// assert
	assert.NotNil(t, beane)
}
