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

//=========================

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
	_, ok := bean.(*beanStract)
	assert.True(t, ok)
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

func Test_applicationContext_GetBean_fromFactory_withPointerParameter(t *testing.T) {
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

func Test_applicationContext_GetBean_fromFactory_withNonPointerParameter(t *testing.T) {
	// arrange
	beans := Beans(
		Bean(string("")).ID("id_1").Factory(func(in string) *string {
			s := "Hi " + in
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

func Test_applicationContext_GetBean_beanHaveIntProperty(t *testing.T) {
	// arrange
	type myBean struct {
		I int
	}
	beans := Beans(
		Bean(myBean{}).ID("id_1").Property("I", 123),
	)

	// action
	ctx, e1 := NewApplicationContext(beans...)
	require.Nil(t, e1)
	bean, e2 := ctx.GetBean("id_1")
	require.Nil(t, e2)

	// assert
	my, ok := bean.(*myBean)
	require.True(t, ok)
	assert.Equal(t, 123, my.I)
}

func Test_applicationContext_GetBean_beanHaveStringProperty(t *testing.T) {
	// arrange
	type myBean struct {
		S string
	}
	beans := Beans(
		Bean(myBean{}).ID("id_1").Property("S", "abc"),
	)

	// action
	ctx, e1 := NewApplicationContext(beans...)
	require.Nil(t, e1)
	bean, e2 := ctx.GetBean("id_1")
	require.Nil(t, e2)

	// assert
	my, ok := bean.(*myBean)
	require.True(t, ok)
	assert.Equal(t, "abc", my.S)
}

func Test_applicationContext_GetBean_beanHaveStructProperty(t *testing.T) {
	// arrange
	type myBean struct {
		S string
	}
	type yourBean struct {
		B myBean
	}
	beans := Beans(
		Bean(yourBean{}).ID("id_1").
			Property("B", Bean(myBean{}).Property("S", "abc")),
	)

	// action
	ctx, e1 := NewApplicationContext(beans...)
	require.Nil(t, e1)
	bean, e2 := ctx.GetBean("id_1")
	require.Nil(t, e2)

	// assert
	your, ok := bean.(*yourBean)
	require.True(t, ok)
	assert.Equal(t, "abc", your.B.S)
}

func Test_applicationContext_GetBean_beanHaveIntSliceProperty(t *testing.T) {
	// arrange
	type myBean struct {
		Is []int
	}
	beans := Beans(
		Bean(myBean{}).ID("id_1").Property("Is", 123, 456, 789),
	)

	// action
	ctx, e1 := NewApplicationContext(beans...)
	require.Nil(t, e1)
	bean, e2 := ctx.GetBean("id_1")
	require.Nil(t, e2)

	// assert
	my, ok := bean.(*myBean)
	require.True(t, ok)
	assert.Equal(t, []int{123, 456, 789}, my.Is)
}

func Test_applicationContext_GetBean_beanHaveStringSliceProperty(t *testing.T) {
	// arrange
	type myBean struct {
		Ss []string
	}
	beans := Beans(
		Bean(myBean{}).ID("id_1").Property("Ss", "aaa", "bbb", "ccc"),
	)

	// action
	ctx, e1 := NewApplicationContext(beans...)
	require.Nil(t, e1)
	bean, e2 := ctx.GetBean("id_1")
	require.Nil(t, e2)

	// assert
	my, ok := bean.(*myBean)
	require.True(t, ok)
	assert.Equal(t, []string{"aaa", "bbb", "ccc"}, my.Ss)
}

func Test_applicationContext_GetBean_beanHaveStructSliceProperty(t *testing.T) {
	// arrange
	type myBean struct {
		S string
	}
	type yourBean struct {
		B []myBean
	}
	beans := Beans(
		Bean(yourBean{}).ID("id_1").
			Property("B",
				Bean(myBean{}).Property("S", "aaa"),
				Bean(myBean{}).Property("S", "bbb"),
				Bean(myBean{}).Property("S", "ccc"),
			),
	)

	// action
	ctx, e1 := NewApplicationContext(beans...)
	require.Nil(t, e1)
	bean, e2 := ctx.GetBean("id_1")
	require.Nil(t, e2)

	// assert
	your, ok := bean.(*yourBean)
	require.True(t, ok)
	assert.Equal(t, []myBean{
		myBean{S: "aaa"},
		myBean{S: "bbb"},
		myBean{S: "ccc"},
	}, your.B)
}

var Test_applicationContext_GetBean_withDefaultInitFunc_isRun = false

type Test_applicationContext_GetBean_withDefaultInitFunc_struct struct {
	I int
}

func (s *Test_applicationContext_GetBean_withDefaultInitFunc_struct) Init() {
	Test_applicationContext_GetBean_withDefaultInitFunc_isRun = true
}

func Test_applicationContext_GetBean_withDefaultInitFunc(t *testing.T) {
	// arrange
	beans := Beans(
		Bean(Test_applicationContext_GetBean_withDefaultInitFunc_struct{}).
			ID("id_1"),
	)

	// action
	ctx, e1 := NewApplicationContext(beans...)
	require.Nil(t, e1)
	bean, e2 := ctx.GetBean("id_1")
	require.Nil(t, e2)

	// assert
	_, ok := bean.(*Test_applicationContext_GetBean_withDefaultInitFunc_struct)
	require.True(t, ok)
	assert.True(t, Test_applicationContext_GetBean_withDefaultInitFunc_isRun)
}

var Test_applicationContext_GetBean_withCostumeInitFunc_isRun = false

type Test_applicationContext_GetBean_withCostumeInitFunc_struct struct {
	I int
}

func (s *Test_applicationContext_GetBean_withCostumeInitFunc_struct) Aaa() {
	Test_applicationContext_GetBean_withCostumeInitFunc_isRun = true
}

func Test_applicationContext_GetBean_withCostumeInitFunc(t *testing.T) {
	// arrange
	beans := Beans(
		Bean(Test_applicationContext_GetBean_withCostumeInitFunc_struct{}).
			ID("id_1").Init("Aaa"),
	)

	// action
	ctx, e1 := NewApplicationContext(beans...)
	require.Nil(t, e1)
	bean, e2 := ctx.GetBean("id_1")
	require.Nil(t, e2)

	// assert
	_, ok := bean.(*Test_applicationContext_GetBean_withCostumeInitFunc_struct)
	require.True(t, ok)
	assert.True(t, Test_applicationContext_GetBean_withCostumeInitFunc_isRun)
}

var Test_applicationContext_Finalize_withDefaultFunction_isRun = false

type Test_applicationContext_Finalize_withDefaultFunction_struct struct {
	I int
}

func (s *Test_applicationContext_Finalize_withDefaultFunction_struct) Finalize() {
	Test_applicationContext_Finalize_withDefaultFunction_isRun = true
}

func Test_applicationContext_Finalize_withDefaultFunction(t *testing.T) {
	// arrange
	beans := Beans(
		Bean(Test_applicationContext_Finalize_withDefaultFunction_struct{}).
			ID("id_1"),
	)
	ctx, e1 := NewApplicationContext(beans...)
	require.Nil(t, e1)
	_, e2 := ctx.GetBean("id_1")
	require.Nil(t, e2)

	// action
	ctx.Finalize()

	// assert
	assert.True(t, Test_applicationContext_Finalize_withDefaultFunction_isRun)
}

var Test_applicationContext_Finalize_withCostumeFunction_isRun = false

type Test_applicationContext_Finalize_withCostumeFunction_struct struct {
	I int
}

func (s *Test_applicationContext_Finalize_withCostumeFunction_struct) Aaa() {
	Test_applicationContext_Finalize_withCostumeFunction_isRun = true
}

func Test_applicationContext_Finalize_withCostumeFunction(t *testing.T) {
	// arrange
	beans := Beans(
		Bean(Test_applicationContext_Finalize_withCostumeFunction_struct{}).
			ID("id_1").Finalize("Aaa"),
	)
	ctx, e1 := NewApplicationContext(beans...)
	require.Nil(t, e1)
	_, e2 := ctx.GetBean("id_1")
	require.Nil(t, e2)

	// action
	ctx.Finalize()

	// assert
	assert.True(t, Test_applicationContext_Finalize_withCostumeFunction_isRun)
}
