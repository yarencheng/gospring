package gospring

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"

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
	assert.NotNilf(t, e, "%v", e)
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

var Test_Finalize_checkOrder_flag = 1

type Test_Finalize_checkOrder_struct1 struct {
	i int
	V interface{}
}
type Test_Finalize_checkOrder_struct2 struct {
	i int
	V interface{}
}
type Test_Finalize_checkOrder_struct3 struct {
	i int
	V interface{}
}

func (s *Test_Finalize_checkOrder_struct1) Finalize() {
	s.i = Test_Finalize_checkOrder_flag
	Test_Finalize_checkOrder_flag++
}

func (s *Test_Finalize_checkOrder_struct2) Finalize() {
	s.i = Test_Finalize_checkOrder_flag
	Test_Finalize_checkOrder_flag++
}

func (s *Test_Finalize_checkOrder_struct3) Finalize() {
	s.i = Test_Finalize_checkOrder_flag
	Test_Finalize_checkOrder_flag++
}

func Test_Finalize_checkOrder(t *testing.T) {
	// arrange
	beans := Beans(
		Bean(Test_Finalize_checkOrder_struct1{}).
			ID("id_1").
			Property("V", Ref("id_2")),
		Bean(Test_Finalize_checkOrder_struct2{}).
			ID("id_2"),
		Bean(Test_Finalize_checkOrder_struct3{}).
			ID("id_3"),
	)

	ctx, e := NewApplicationContext(beans...)
	require.Nil(t, e)
	bean1, e := ctx.GetBean("id_1")
	require.Nil(t, e)
	bean3, e := ctx.GetBean("id_3")
	require.Nil(t, e)
	bean2, e := ctx.GetBean("id_2")
	require.Nil(t, e)

	// action
	ef := ctx.Finalize()
	require.Nil(t, ef)

	// assert
	assert.Equal(t, 2, bean1.(*Test_Finalize_checkOrder_struct1).i)
	assert.Equal(t, 3, bean2.(*Test_Finalize_checkOrder_struct2).i)
	assert.Equal(t, 1, bean3.(*Test_Finalize_checkOrder_struct3).i)
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

func Test_checkScope_singleton(t *testing.T) {
	// arrange
	type beanStruct struct{}
	beans := Beans(
		Bean(beanStruct{}).Singleton(),
	)

	// action
	_, e := NewApplicationContext(beans...)

	// assert
	require.Nil(t, e)
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

func Test_addBean(t *testing.T) {
	// arrange
	type beanStruct struct{}
	beans := Beans(
		Bean(beanStruct{}),
	)

	// action
	_, e := NewApplicationContext(beans...)

	// assert
	assert.Nil(t, e)
}

func Test_addBean_recursiveError(t *testing.T) {
	// arrange
	type beanStruct struct {
		B interface{}
	}
	beans := Beans(
		Bean(beanStruct{}).
			ID("id").
			Property(
				"B",
				Bean(beanStruct{}).
					ID("id"),
			),
	)

	// action
	_, e := NewApplicationContext(beans...)

	// assert
	assert.NotNil(t, e)
}

func Test_checkDependencyLoop_noPrentID(t *testing.T) {
	// arrange
	type beanStruct struct {
		B interface{}
	}
	beans := Beans(
		Bean(beanStruct{}).
			Property(
				"B",
				Bean(beanStruct{}),
			),
	)

	// action
	_, e := NewApplicationContext(beans...)

	// assert
	assert.Nil(t, e)
}

func Test_checkDependencyLoop_referenceLoop_1(t *testing.T) {
	// arrange
	type beanStruct1 struct {
		B2 interface{}
	}
	type beanStruct2 struct {
		B1 interface{}
	}
	beans := Beans(
		Bean(beanStruct1{}).ID("id1").Property("B2", Ref("id2")),
		Bean(beanStruct2{}).ID("id2").Property("B1", Ref("id1")),
	)

	// action
	_, e := NewApplicationContext(beans...)

	// assert
	assert.NotNil(t, e)
}

func Test_checkDependencyLoop_referenceLoop_2(t *testing.T) {
	// arrange
	type beanStruct struct {
		B interface{}
	}
	beans := Beans(
		Bean(beanStruct{}).ID("1").Property("B", Ref("2")),
		Bean(beanStruct{}).ID("2").Property("B", Ref("3")),
		Bean(beanStruct{}).ID("3").Property("B", Ref("1")),
	)

	// action
	_, e := NewApplicationContext(beans...)

	// assert
	assert.NotNil(t, e)
}

func Test_checkDependencyLoop_referenceLoop_3(t *testing.T) {
	// arrange
	type beanStruct struct {
		B interface{}
	}
	beans := Beans(
		Bean(beanStruct{}).ID("1").Property("B",
			Ref("1")),
	)

	// action
	_, e := NewApplicationContext(beans...)

	// assert
	assert.NotNil(t, e)
}

func Test_checkDependencyLoop_referenceLoop_4(t *testing.T) {
	// arrange
	type beanStruct struct {
		B interface{}
	}
	beans := Beans(
		Bean(beanStruct{}).ID("1").Property("B",
			Bean(beanStruct{}).Property("B",
				Ref("1")),
		),
	)

	// action
	_, e := NewApplicationContext(beans...)

	// assert
	assert.NotNil(t, e)
}

func Test_checkDependencyLoop_referenceLoop_6(t *testing.T) {
	// arrange
	type beanStruct struct {
		B interface{}
	}
	beans := Beans(
		Bean(beanStruct{}).ID("1").Property("B",
			Bean(beanStruct{}).Property("B",
				Bean(beanStruct{}).Property("B",
					Ref("1")),
			),
		),
	)

	// action
	_, e := NewApplicationContext(beans...)

	// assert
	assert.NotNil(t, e)
}

func Test_GetBean_int(t *testing.T) {
	// arrange
	beans := Beans(
		Bean(int(123)).ID("1"),
	)
	ctx, e := NewApplicationContext(beans...)
	require.Nil(t, e)

	// action
	bean, e := ctx.GetBean("1")
	require.Nil(t, e)

	// assert
	actual, ok := bean.(*int)
	require.True(t, ok)
	assert.Equal(t, int(123), *actual)
}

func Test_GetBean_string(t *testing.T) {
	// arrange
	beans := Beans(
		Bean(string("abcd")).ID("1"),
	)
	ctx, e := NewApplicationContext(beans...)
	require.Nil(t, e)

	// action
	bean, e := ctx.GetBean("1")
	require.Nil(t, e)

	// assert
	actual, ok := bean.(*string)
	require.True(t, ok)
	assert.Equal(t, "abcd", *actual)
}

func Test_GetBean_Singleton(t *testing.T) {
	// arrange
	type beanStruct struct {
		I int
	}
	beans := Beans(
		Bean(beanStruct{}).ID("1").Singleton(),
	)
	ctx, e := NewApplicationContext(beans...)
	require.Nil(t, e)

	// action
	bean1, e := ctx.GetBean("1")
	require.Nil(t, e)
	bean2, e := ctx.GetBean("1")
	require.Nil(t, e)

	// assert
	assert.Equal(t,
		unsafe.Pointer(bean1.(*beanStruct)),
		unsafe.Pointer(bean2.(*beanStruct)),
	)
}

func Test_GetBean_Property(t *testing.T) {
	// arrange
	type beanStruct struct {
		I int
	}
	beans := Beans(
		Bean(beanStruct{}).ID("1").Prototype(),
	)
	ctx, e := NewApplicationContext(beans...)
	require.Nil(t, e)

	// action
	bean1, e := ctx.GetBean("1")
	require.Nil(t, e)
	bean2, e := ctx.GetBean("1")
	require.Nil(t, e)

	// assert
	assert.NotEqual(t,
		unsafe.Pointer(bean1.(*beanStruct)),
		unsafe.Pointer(bean2.(*beanStruct)),
	)
}

func Test_GetBean_factoryFail(t *testing.T) {
	// arrange
	type beanStruct struct {
		I int
	}
	beans := Beans(
		Bean(beanStruct{}).ID("1").Factory(
			func() (*beanStruct, error) {
				return nil, fmt.Errorf("")
			},
		),
	)
	ctx, e := NewApplicationContext(beans...)
	require.Nil(t, e)

	// action
	bean, e := ctx.GetBean("1")

	// assert
	assert.Nil(t, bean)
	assert.NotNil(t, e)
}

func Test_GetBean_injectFail(t *testing.T) {
	// arrange
	type beanStruct struct {
		I *int
	}
	beans := Beans(
		Bean(beanStruct{}).ID("1"),
		Bean(beanStruct{}).ID("2").Property("I", Ref("1")),
	)
	ctx, e := NewApplicationContext(beans...)
	require.Nil(t, e)

	// action
	bean, e := ctx.GetBean("2")

	// assert
	assert.Nil(t, bean)
	assert.NotNil(t, e)
}

func Test_GetBean_injectSliceFail(t *testing.T) {
	// arrange
	type beanStruct struct {
		I []*int
	}
	beans := Beans(
		Bean(beanStruct{}).ID("1"),
		Bean(beanStruct{}).ID("2").Property("I", Ref("1")),
	)
	ctx, e := NewApplicationContext(beans...)
	require.Nil(t, e)

	// action
	bean, e := ctx.GetBean("2")

	// assert
	assert.Nil(t, bean)
	assert.NotNil(t, e)
}

func Test_createBeanByFactory_getArgvFailed(t *testing.T) {
	// arrange
	type beanStruct struct {
		I interface{}
	}
	beans := Beans(
		Bean(beanStruct{}).ID("1").Factory(
			func(*beanStruct) *beanStruct { return &beanStruct{} },
			Bean(beanStruct{}).Factory( // argv
				func() (*beanStruct, error) {
					return nil, fmt.Errorf("") // return error from factory
				},
			),
		),
	)
	ctx, e := NewApplicationContext(beans...)
	require.Nil(t, e)

	// action
	bean, e := ctx.GetBean("1")

	// assert
	assert.Nil(t, bean)
	assert.NotNil(t, e)
}

func Test_createBeanByFactory_injectArgvAsPointer(t *testing.T) {
	// arrange
	type beanStruct struct {
		I interface{}
	}
	beans := Beans(
		Bean(beanStruct{}).ID("1").Factory(
			func(*beanStruct) *beanStruct { return &beanStruct{} },
			Bean(beanStruct{}),
		),
	)
	ctx, e := NewApplicationContext(beans...)
	require.Nil(t, e)

	// action
	bean, e := ctx.GetBean("1")

	// assert
	assert.NotNil(t, bean)
	assert.Nil(t, e)
}

func Test_createBeanByFactory_injectArgvAsElem(t *testing.T) {
	// arrange
	type beanStruct struct {
		I interface{}
	}
	beans := Beans(
		Bean(beanStruct{}).ID("1").Factory(
			func(beanStruct) *beanStruct { return &beanStruct{} },
			Bean(beanStruct{}),
		),
	)
	ctx, e := NewApplicationContext(beans...)
	require.Nil(t, e)

	// action
	bean, e := ctx.GetBean("1")

	// assert
	assert.NotNil(t, bean)
	assert.Nil(t, e)
}

func Test_createBeanByFactory_injectSingletonArgvAsElem(t *testing.T) {
	// arrange
	type beanStruct struct {
		I interface{}
	}
	beans := Beans(
		Bean(beanStruct{}).ID("1").Factory(
			func(beanStruct) *beanStruct { return &beanStruct{} },
			Bean(beanStruct{}).Singleton(),
		),
	)
	ctx, e := NewApplicationContext(beans...)
	require.Nil(t, e)

	// action
	bean, e := ctx.GetBean("1")

	// assert
	assert.Nil(t, bean)
	assert.NotNil(t, e)
}

func Test_createBeanByFactory_convertArgvFailed(t *testing.T) {
	// arrange
	type beanStruct struct {
		I interface{}
	}
	beans := Beans(
		Bean(beanStruct{}).ID("1").Factory(
			func(int) *beanStruct { return &beanStruct{} },
			Bean(beanStruct{}),
		),
	)
	ctx, e := NewApplicationContext(beans...)
	require.Nil(t, e)

	// action
	bean, e := ctx.GetBean("1")

	// assert
	assert.Nil(t, bean)
	assert.NotNil(t, e)
}

func Test_createBeanByFactory_returnError(t *testing.T) {
	// arrange
	type beanStruct struct {
		I interface{}
	}
	beans := Beans(
		Bean(beanStruct{}).ID("1").Factory(
			func() error { return fmt.Errorf("") },
		),
	)
	ctx, e := NewApplicationContext(beans...)
	require.Nil(t, e)

	// action
	bean, e := ctx.GetBean("1")

	// assert
	assert.Nil(t, bean)
	assert.NotNil(t, e)
}

func Test_createBeanByFactory_returnNilError(t *testing.T) {
	// arrange
	type beanStruct struct {
		I interface{}
	}
	beans := Beans(
		Bean(beanStruct{}).ID("1").Factory(
			func() (*beanStruct, error) { return &beanStruct{}, nil },
		),
	)
	ctx, e := NewApplicationContext(beans...)
	require.Nil(t, e)

	// action
	bean, e := ctx.GetBean("1")

	// assert
	assert.NotNil(t, bean)
	assert.Nil(t, e)
}

func Test_createBeanByFactory_returnTwoValuesAndTheSecondIsError(t *testing.T) {
	// arrange
	type beanStruct struct {
		I interface{}
	}
	beans := Beans(
		Bean(beanStruct{}).ID("1").Factory(
			func() (*beanStruct, error) { return nil, fmt.Errorf("") },
		),
	)
	ctx, e := NewApplicationContext(beans...)
	require.Nil(t, e)

	// action
	bean, e := ctx.GetBean("1")

	// assert
	assert.Nil(t, bean)
	assert.NotNil(t, e)
}

func Test_inject_getBeanFailed(t *testing.T) {
	// arrange
	type beanStruct struct {
		I interface{}
	}
	beans := Beans(
		Bean(beanStruct{}).ID("1").Property("I",
			Bean(beanStruct{}).Factory(func() error { return fmt.Errorf("") }),
		),
	)
	ctx, e := NewApplicationContext(beans...)
	require.Nil(t, e)

	// action
	bean, e := ctx.GetBean("1")

	// assert
	assert.Nil(t, bean)
	assert.NotNil(t, e)
}

func Test_inject_convertFailed(t *testing.T) {
	// arrange
	type beanStruct struct {
		I *beanStruct
	}
	beans := Beans(
		Bean(beanStruct{}).ID("1").Property("I",
			Bean(123).Factory(func() *int {
				i := 123
				return &i
			}),
		),
	)
	ctx, e := NewApplicationContext(beans...)
	require.Nil(t, e)

	// action
	bean, e := ctx.GetBean("1")

	// assert
	assert.Nil(t, bean)
	assert.NotNil(t, e)
}

func Test_inject_singletonToElem(t *testing.T) {
	// arrange
	type beanStruct1 struct {
	}
	type beanStruct2 struct {
		I beanStruct1
	}
	beans := Beans(
		Bean(beanStruct2{}).ID("1").Property("I",
			Bean(beanStruct1{}).
				Singleton(),
		),
	)
	ctx, e := NewApplicationContext(beans...)
	require.Nil(t, e)

	// action
	bean, e := ctx.GetBean("1")

	// assert
	assert.Nil(t, bean)
	assert.NotNil(t, e)
}

func Test_inject(t *testing.T) {
	// arrange
	type beanStruct1 struct {
	}
	type beanStruct2 struct {
		I *beanStruct1
	}
	beans := Beans(
		Bean(beanStruct2{}).ID("1").Property("I",
			Bean(beanStruct1{}),
		),
	)
	ctx, e := NewApplicationContext(beans...)
	require.Nil(t, e)

	// action
	bean, e := ctx.GetBean("1")

	// assert
	assert.NotNil(t, bean)
	assert.Nil(t, e)
}

func Test_injectSlice(t *testing.T) {
	// arrange
	type beanStruct1 struct {
	}
	type beanStruct2 struct {
		I []*beanStruct1
	}
	beans := Beans(
		Bean(beanStruct2{}).ID("1").Property("I",
			Bean(beanStruct1{}),
		),
	)
	ctx, e := NewApplicationContext(beans...)
	require.Nil(t, e)

	// action
	bean, e := ctx.GetBean("1")

	// assert
	assert.NotNil(t, bean)
	assert.Nil(t, e)
}

func Test_injectSlice_getBeanFailed(t *testing.T) {
	// arrange
	type beanStruct struct {
		I []interface{}
	}
	beans := Beans(
		Bean(beanStruct{}).ID("1").Property("I",
			Bean(beanStruct{}).Factory(func() error { return fmt.Errorf("") }),
		),
	)
	ctx, e := NewApplicationContext(beans...)
	require.Nil(t, e)

	// action
	bean, e := ctx.GetBean("1")

	// assert
	assert.Nil(t, bean)
	assert.NotNil(t, e)
}

func Test_injectSlice_singletonToElem(t *testing.T) {
	// arrange
	type beanStruct1 struct {
	}
	type beanStruct2 struct {
		I []beanStruct1
	}
	beans := Beans(
		Bean(beanStruct2{}).ID("1").Property("I",
			Bean(beanStruct1{}).
				Singleton(),
		),
	)
	ctx, e := NewApplicationContext(beans...)
	require.Nil(t, e)

	// action
	bean, e := ctx.GetBean("1")

	// assert
	assert.Nil(t, bean)
	assert.NotNil(t, e)
}

func Test_injectSlice_prototypeToElem(t *testing.T) {
	// arrange
	type beanStruct1 struct {
	}
	type beanStruct2 struct {
		I []beanStruct1
	}
	beans := Beans(
		Bean(beanStruct2{}).ID("1").Property("I",
			Bean(beanStruct1{}).
				Prototype(),
		),
	)
	ctx, e := NewApplicationContext(beans...)
	require.Nil(t, e)

	// action
	bean, e := ctx.GetBean("1")

	// assert
	assert.NotNil(t, bean)
	assert.Nil(t, e)
}

var Test_callInitFunc_isCalled = false

type Test_callInitFunc_struct struct{}

func (t *Test_callInitFunc_struct) Init() {
	Test_callInitFunc_isCalled = true
}

func Test_callInitFunc(t *testing.T) {
	// arrange
	beans := Beans(
		Bean(Test_callInitFunc_struct{}).ID("1"),
	)
	ctx, e := NewApplicationContext(beans...)
	require.Nil(t, e)

	// action
	ctx.GetBean("1")

	// assert
	assert.True(t, Test_callInitFunc_isCalled)
}

type Test_callInitFunc_returnNilError_struct struct{}

func (t *Test_callInitFunc_returnNilError_struct) Init() error {
	return nil
}

func Test_callInitFunc_returnNilError(t *testing.T) {
	// arrange
	beans := Beans(
		Bean(Test_callInitFunc_returnNilError_struct{}).ID("1"),
	)
	ctx, e := NewApplicationContext(beans...)
	require.Nil(t, e)

	// action
	bean, e := ctx.GetBean("1")

	// assert
	assert.NotNil(t, bean)
	assert.Nil(t, e)
}

func Test_callInitFunc_cantFindMethod(t *testing.T) {
	// arrange
	type beanStruct struct{}
	beans := Beans(
		Bean(beanStruct{}).ID("1").Init("aaa"),
	)
	ctx, e := NewApplicationContext(beans...)
	require.Nil(t, e)

	// action
	bean, e := ctx.GetBean("1")

	// assert
	assert.Nil(t, bean)
	assert.NotNil(t, e)
}

type Test_callInitFunc_returnNonError_struct struct{}

func (t *Test_callInitFunc_returnNonError_struct) Init() int {
	return 123
}

func Test_callInitFunc_returnNonError(t *testing.T) {
	// arrange
	beans := Beans(
		Bean(Test_callInitFunc_returnNonError_struct{}).ID("1"),
	)
	ctx, e := NewApplicationContext(beans...)
	require.Nil(t, e)

	// action
	bean, e := ctx.GetBean("1")

	// assert
	assert.Nil(t, bean)
	assert.NotNil(t, e)
}

type Test_callInitFunc_returnMoreThanTwoValues_struct struct{}

func (t *Test_callInitFunc_returnMoreThanTwoValues_struct) Init() (int, int) {
	return 123, 123
}

func Test_callInitFunc_returnMoreThanTwoValues(t *testing.T) {
	// arrange
	beans := Beans(
		Bean(Test_callInitFunc_returnMoreThanTwoValues_struct{}).ID("1"),
	)
	ctx, e := NewApplicationContext(beans...)
	require.Nil(t, e)

	// action
	bean, e := ctx.GetBean("1")

	// assert
	assert.Nil(t, bean)
	assert.NotNil(t, e)
}

func Test_callFinalizeFunc_cantFindMethod(t *testing.T) {
	// arrange
	type beanStruct struct{}
	beans := Beans(
		Bean(beanStruct{}).ID("1").Finalize("aaa"),
	)
	ctx, e := NewApplicationContext(beans...)
	require.Nil(t, e)
	bean, e := ctx.GetBean("1")
	require.NotNil(t, bean)
	require.Nil(t, e)

	// action
	e = ctx.Finalize()

	// assert
	assert.NotNil(t, e)
}

func Test_callFinalizeFunc_notDefine(t *testing.T) {
	// arrange
	type beanStruct struct{}
	beans := Beans(
		Bean(beanStruct{}).ID("1"),
	)
	ctx, e := NewApplicationContext(beans...)
	require.Nil(t, e)
	bean, e := ctx.GetBean("1")
	require.NotNil(t, bean)
	require.Nil(t, e)

	// action
	e = ctx.Finalize()

	// assert
	assert.Nil(t, e)
}

type Test_callFinalizeFunc_returnNilError_struct struct{}

func (t *Test_callFinalizeFunc_returnNilError_struct) Finalize() error {
	return nil
}

func Test_callFinalizeFunc_returnNilError(t *testing.T) {
	// arrange
	beans := Beans(
		Bean(Test_callFinalizeFunc_returnNilError_struct{}).ID("1"),
	)
	ctx, e := NewApplicationContext(beans...)
	require.Nil(t, e)
	bean, e := ctx.GetBean("1")
	require.NotNil(t, bean)
	require.Nil(t, e)

	// action
	e = ctx.Finalize()

	// assert
	assert.Nil(t, e)
}

type Test_callFinalizeFunc_returnNonError_struct struct{}

func (t *Test_callFinalizeFunc_returnNonError_struct) Finalize() int {
	return 123
}

func Test_callFinalizeFunc_returnNonError(t *testing.T) {
	// arrange
	beans := Beans(
		Bean(Test_callFinalizeFunc_returnNonError_struct{}).ID("1"),
	)
	ctx, e := NewApplicationContext(beans...)
	require.Nil(t, e)
	bean, e := ctx.GetBean("1")
	require.NotNil(t, bean)
	require.Nil(t, e)

	// action
	e = ctx.Finalize()

	// assert
	assert.NotNil(t, e)
}

type Test_callFinalizeFunc_returnMoreThanTwoValues_struct struct{}

func (t *Test_callFinalizeFunc_returnMoreThanTwoValues_struct) Finalize() (int, int) {
	return 123, 123
}

func Test_callFinalizeFunc_returnMoreThanTwoValues(t *testing.T) {
	// arrange
	beans := Beans(
		Bean(Test_callFinalizeFunc_returnMoreThanTwoValues_struct{}).ID("1"),
	)
	ctx, e := NewApplicationContext(beans...)
	require.Nil(t, e)
	bean, e := ctx.GetBean("1")
	require.NotNil(t, bean)
	require.Nil(t, e)

	// action
	e = ctx.Finalize()

	// assert
	assert.NotNil(t, e)
}

func Test_GetBean_channel(t *testing.T) {
	// arrange
	beans := Beans(
		Chan(int(123), 1).ID("1"),
	)
	ctx, e := NewApplicationContext(beans...)
	require.Nil(t, e)

	// action
	bean, e := ctx.GetBean("1")

	// assert
	assert.NotNil(t, bean)
	assert.Nil(t, e)

	// assert - check channel is usable
	c := bean.(chan int)
	c <- 456
	assert.Equal(t, 456, <-c)
}

func Test_GetBean_channelWithFactory(t *testing.T) {
	// arrange
	isCalled := false
	beans := Beans(
		Chan(int(123), 1).
			ID("1").
			Factory(func() chan int {
				isCalled = true
				return make(chan int, 1)
			}),
	)
	ctx, e := NewApplicationContext(beans...)
	require.Nil(t, e)

	// action
	bean, e := ctx.GetBean("1")

	// assert
	assert.True(t, isCalled)
	assert.NotNil(t, bean)
	assert.Nil(t, e)

	// assert - check channel is usable
	c := bean.(chan int)
	c <- 456
	assert.Equal(t, 456, <-c)
}

func Test_GetBean_withChannel(t *testing.T) {
	// arrange
	type myStruct struct {
		C chan int
	}
	beans := Beans(
		Bean(myStruct{}).
			ID("1").
			Property("C", Ref("2")),
		Chan(int(123), 1).
			ID("2"),
	)
	ctx, e := NewApplicationContext(beans...)
	require.Nil(t, e)

	// action
	bean, e := ctx.GetBean("1")
	require.Nil(t, e)
	my, ok := bean.(*myStruct)
	require.True(t, ok)

	// assert - check channel is usable
	my.C <- 456
	assert.Equal(t, 456, <-my.C)
}
