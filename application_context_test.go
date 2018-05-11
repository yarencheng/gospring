package gospring

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

func Test_GetBean(t *testing.T) {
	// arrange
	type beanStruct struct{}
	beans := Beans(
		Bean(beanStruct{}).ID("id"),
	)

	ctx, e := NewApplicationContext(beans...)
	require.Nil(t, e)

	// action
	bean, e := ctx.GetBean("id")

	// assert
	assert.NotNil(t, bean)
}

func Test_GetBean_invalidID(t *testing.T) {
	// arrange
	type beanStruct struct{}
	beans := Beans(
		Bean(beanStruct{}).ID("id"),
	)

	ctx, e := NewApplicationContext(beans...)
	require.Nil(t, e)

	// action
	bean, e := ctx.GetBean("invalidID")

	// assert
	assert.Nil(t, bean)
	assert.NotNil(t, e)
}

type Test_GetBean_error_struct struct {
}

func (*Test_GetBean_error_struct) Init() error {
	return fmt.Errorf("")
}

func Test_GetBean_failedToCreateBean(t *testing.T) {
	// arrange
	beans := Beans(
		Bean(Test_GetBean_error_struct{}).ID("id"),
	)

	ctx, e := NewApplicationContext(beans...)
	require.Nil(t, e)

	// action
	bean, e := ctx.GetBean("id")

	// assert
	assert.Nil(t, bean)
	assert.NotNil(t, e)
}

type Test_Finalize_struct struct {
	b bool
}

func (s *Test_Finalize_struct) Finalize() {
	s.b = false
}

func Test_Finalize(t *testing.T) {
	// arrange
	beans := Beans(
		Bean(Test_Finalize_struct{}).ID("id"),
	)

	ctx, e := NewApplicationContext(beans...)
	require.Nil(t, e)
	bean, e := ctx.GetBean("id")
	require.Nil(t, e)

	bean.(*Test_Finalize_struct).b = true

	// action
	ef := ctx.Finalize()
	require.Nil(t, ef)

	// assert
	assert.False(t, bean.(*Test_Finalize_struct).b)
}

type Test_Finalize_error_struct struct {
}

func (s *Test_Finalize_error_struct) Finalize() error {
	return fmt.Errorf("")
}

func Test_Finalize_error(t *testing.T) {
	// arrange
	beans := Beans(
		Bean(Test_Finalize_error_struct{}).ID("id"),
	)

	ctx, e := NewApplicationContext(beans...)
	require.Nil(t, e)
	_, eb := ctx.GetBean("id")
	require.Nil(t, eb)

	// action
	ef := ctx.Finalize()

	// assert
	assert.NotNil(t, ef)
}

func Test_setRefBean_withProperty(t *testing.T) {
	// arrange
	type beanStruct1 struct {
		I int
	}
	type beanStruct2 struct {
		B beanStruct1
	}
	beans := Beans(
		Bean(beanStruct1{}).ID("id_1").Property("I", 123),
		Bean(beanStruct2{}).ID("id_2").Property("B", Ref("id_1")),
	)

	// action
	ctx, e := NewApplicationContext(beans...)
	require.Nil(t, e)
	bean, e := ctx.GetBean("id_2")
	require.Nil(t, e)

	// assert
	bean2, ok := bean.(*beanStruct2)
	require.True(t, ok)
	assert.Equal(t, 123, bean2.B.I)
}

func Test_setRefBean_withNotExistProperty(t *testing.T) {
	// arrange
	type beanStruct1 struct {
		I int
	}
	type beanStruct2 struct {
		B beanStruct1
	}
	beans := Beans(
		Bean(beanStruct1{}).ID("id_1").Property("I", 123),
		Bean(beanStruct2{}).ID("id_2").Property("B", Ref("id_aaa")),
	)

	// action
	_, e := NewApplicationContext(beans...)

	// assert
	assert.NotNil(t, e)
}

func Test_setRefBean_withFactory(t *testing.T) {
	// arrange
	type beanStruct1 struct {
		I int
	}
	type beanStruct2 struct {
		B beanStruct1
	}
	beans := Beans(
		Bean(beanStruct1{}).ID("id_1").Property("I", 123),
		Bean(beanStruct2{}).ID("id_2").Factory(func(s1 *beanStruct1) *beanStruct2 {
			return &beanStruct2{
				B: *s1,
			}
		}, Ref("id_1")),
	)

	// action
	ctx, e := NewApplicationContext(beans...)
	require.Nil(t, e)
	bean, e := ctx.GetBean("id_2")
	require.Nil(t, e)

	// assert
	bean2, ok := bean.(*beanStruct2)
	require.True(t, ok)
	assert.Equal(t, 123, bean2.B.I)
}

func Test_setRefBean_idNotExist(t *testing.T) {
	// arrange
	type beanStruct1 struct {
		I int
	}
	type beanStruct2 struct {
		B beanStruct1
	}
	beans := Beans(
		Bean(beanStruct1{}).ID("id_1").Property("I", 123),
		Bean(beanStruct2{}).ID("id_2").Factory(func(s1 *beanStruct1) *beanStruct2 {
			return &beanStruct2{
				B: *s1,
			}
		}, Ref("id_aaaa")),
	)

	// action
	_, e := NewApplicationContext(beans...)

	// assert
	assert.NotNil(t, e)
}

func Test_setRefBean_recursiveError(t *testing.T) {
	// arrange
	type beanStruct1 struct {
		I int
	}
	type beanStruct2 struct {
		B beanStruct1
	}
	type beanStruct3 struct {
		B beanStruct2
	}
	beans := Beans(
		Bean(beanStruct3{}).ID("id_3").
			Property("B",
				Bean(beanStruct2{}).
					ID("id_2").
					Property("B", Ref("wrong_id")),
			),
	)

	// action
	_, e := NewApplicationContext(beans...)

	// assert
	assert.NotNil(t, e)
}

func Test_addBeanById(t *testing.T) {
	// arrange
	type beanStruct struct{}
	beans := Beans(
		Bean(beanStruct{}).ID("id"),
	)

	// action
	ctx, e := NewApplicationContext(beans...)
	require.Nil(t, e)

	// assert
	actx, ok := ctx.(*applicationContext)
	require.True(t, ok)
	assert.Contains(t, actx.beanById, "id")
}

func Test_addBeanById_conflict(t *testing.T) {
	// arrange
	type beanStruct struct{}
	beans := Beans(
		Bean(beanStruct{}).ID("id"),
		Bean(beanStruct{}).ID("id"),
	)

	// action
	_, e := NewApplicationContext(beans...)

	// assert
	assert.NotNil(t, e)
}

func Test_checkType(t *testing.T) {
	// arrange
	type beanStruct struct{}
	beans := Beans(
		Bean(beanStruct{}),
	)

	// action
	_, e := NewApplicationContext(beans...)

	// assert
	require.Nil(t, e)
}

func Test_checkType_pointer(t *testing.T) {
	// arrange
	type beanStruct struct{}
	beans := Beans(
		Bean(&beanStruct{}),
	)

	// action
	_, e := NewApplicationContext(beans...)

	// assert
	require.NotNil(t, e)
}

func Test_checkScope_prototypeCantHaveFinalizer(t *testing.T) {
	// arrange
	type beanStruct struct{}
	beans := Beans(
		Bean(beanStruct{}).Prototype().Finalize("aa"),
	)

	// action
	_, e := NewApplicationContext(beans...)

	// assert
	require.NotNil(t, e)
}

func Test_checkScope_unknownScope(t *testing.T) {
	// arrange
	type beanStruct struct{}
	var scope Scope = "sss"

	mock := new(BeanMock)
	mock.On("GetScope").Return(scope)
	mock.On("GetID").Return(nil)
	mock.On("GetType").Return(reflect.TypeOf(beanStruct{}))
	mock.On("GetFactory").Return(
		func() *beanStruct { return nil },
		[]BeanI{},
	)

	beans := Beans(
		mock,
	)

	// action
	_, e := NewApplicationContext(beans...)

	// assert
	require.NotNil(t, e)
}

func Test_checkFactory(t *testing.T) {
	// arrange
	type beanStruct struct{}
	beans := Beans(
		Bean(beanStruct{}).Factory(
			func() *beanStruct {
				return nil
			},
		),
	)

	// action
	_, e := NewApplicationContext(beans...)

	// assert
	require.Nil(t, e)
}

func Test_checkFactory_notFunction(t *testing.T) {
	// arrange
	type beanStruct struct{}
	beans := Beans(
		Bean(beanStruct{}).Factory(
			"something else",
		),
	)

	// action
	_, e := NewApplicationContext(beans...)

	// assert
	require.NotNil(t, e)
}

func Test_checkFactory_tooManyArgument(t *testing.T) {
	// arrange
	type beanStruct struct{}
	beans := Beans(
		Bean(beanStruct{}).Factory(
			func() *beanStruct {
				return nil
			},
			"too many",
		),
	)

	// action
	_, e := NewApplicationContext(beans...)

	// assert
	require.NotNil(t, e)
}

func Test_checkFactory_returnNinPointer(t *testing.T) {
	// arrange
	type beanStruct struct{}
	beans := Beans(
		Bean(beanStruct{}).Factory(
			func() beanStruct {
				return beanStruct{}
			},
		),
	)

	// action
	_, e := NewApplicationContext(beans...)

	// assert
	require.NotNil(t, e)
}

func Test_checkFactory_returnTwoValue(t *testing.T) {
	// arrange
	type beanStruct struct{}
	beans := Beans(
		Bean(beanStruct{}).Factory(
			func() (*beanStruct, error) {
				return nil, nil
			},
		),
	)

	// action
	_, e := NewApplicationContext(beans...)

	// assert
	require.Nil(t, e)
}

func Test_checkFactory_returnTwoValueAndNotPointer(t *testing.T) {
	// arrange
	type beanStruct struct{}
	beans := Beans(
		Bean(beanStruct{}).Factory(
			func() (beanStruct, error) {
				return beanStruct{}, nil
			},
		),
	)

	// action
	_, e := NewApplicationContext(beans...)

	// assert
	require.NotNil(t, e)
}

func Test_checkFactory_returnTwoValueAndNoterror(t *testing.T) {
	// arrange
	type beanStruct struct{}
	beans := Beans(
		Bean(beanStruct{}).Factory(
			func() (*beanStruct, interface{}) {
				return nil, nil
			},
		),
	)

	// action
	_, e := NewApplicationContext(beans...)

	// assert
	require.NotNil(t, e)
}

func Test_checkFactory_returnThreeValue(t *testing.T) {
	// arrange
	type beanStruct struct{}
	beans := Beans(
		Bean(beanStruct{}).Factory(
			func() (*beanStruct, error, interface{}) {
				return nil, nil, nil
			},
		),
	)

	// action
	_, e := NewApplicationContext(beans...)

	// assert
	require.NotNil(t, e)
}

func Test_NewApplicationContext_empty(t *testing.T) {
	// arrange

	// action
	_, e := NewApplicationContext()

	// assert
	assert.Nil(t, e)
}
