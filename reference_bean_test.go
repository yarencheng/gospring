package gospring

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type BeanMock struct {
	mock.Mock
}

func (m *BeanMock) GetID() *string {
	args := m.Called()
	s := args.Get(0)
	if s == nil {
		return nil
	} else {
		return s.(*string)
	}
}

func (m *BeanMock) GetValue() interface{} {
	return m.Called().Get(0)
}

func (m *BeanMock) GetScope() Scope {
	return m.Called().Get(0).(Scope)
}

func (m *BeanMock) GetFactory() (interface{}, []BeanI) {
	args := m.Called()
	return args.Get(0), args.Get(1).([]BeanI)
}

func (m *BeanMock) GetFinalize() *string {
	return m.Called().Get(0).(*string)
}

func (m *BeanMock) GetInit() *string {
	return m.Called().Get(0).(*string)
}

func (m *BeanMock) GetProperty(name string) []BeanI {
	return m.Called(name).Get(0).([]BeanI)
}

func (m *BeanMock) GetProperties() map[string][]BeanI {
	return m.Called().Get(0).(map[string][]BeanI)
}

func (m *BeanMock) GetType() reflect.Type {
	args := m.Called()
	s := args.Get(0)
	if s == nil {
		return nil
	} else {
		return s.(reflect.Type)
	}
}

func Test_GetID_withReference(t *testing.T) {
	// arrange
	expected := "a_id"
	mock := new(BeanMock)
	mock.On("GetID").Return(&expected)
	bean := referenceBean{
		reference: mock,
	}

	// action
	actual := bean.GetID()

	// assert
	assert.Equal(t, expected, *actual)
}

func Test_GetID_withOutReference(t *testing.T) {
	// arrange
	expected := "a_id"
	bean := referenceBean{
		id: expected,
	}

	// action
	actual := bean.GetID()

	// assert
	assert.Equal(t, expected, *actual)
}

func Test_GetScope_withReference(t *testing.T) {
	// arrange
	expected := Prototype
	mock := new(BeanMock)
	mock.On("GetScope").Return(expected)
	bean := referenceBean{
		reference: mock,
	}

	// action
	actual := bean.GetScope()

	// assert
	assert.Equal(t, expected, actual)
}

func Test_GetScope_withOutReference(t *testing.T) {
	// arrange
	bean := referenceBean{}

	// action
	actual := bean.GetScope()

	// assert
	assert.Equal(t, Default, actual)
}

func Test_GetFactory_withReference(t *testing.T) {
	// arrange
	isCalled := false
	mock := new(BeanMock)
	mock.On("GetFactory").Return(
		func() { isCalled = true },
		[]BeanI{nil},
	)
	bean := referenceBean{
		reference: mock,
	}

	// action
	actualFn, actualArgv := bean.GetFactory()
	reflect.ValueOf(actualFn).Call([]reflect.Value{})

	// assert
	assert.True(t, isCalled)
	assert.Equal(t, []BeanI{nil}, actualArgv)
}

func Test_GetFactory_withOutReference(t *testing.T) {
	// arrange
	bean := referenceBean{}

	// action
	actualFn, actualArgv := bean.GetFactory()

	// assert
	assert.Nil(t, actualFn)
	assert.Nil(t, actualArgv)
}

func Test_GetFinalize_withReference(t *testing.T) {
	// arrange
	expected := "a_method"
	mock := new(BeanMock)
	mock.On("GetFinalize").Return(&expected)
	bean := referenceBean{
		reference: mock,
	}

	// action
	actual := bean.GetFinalize()

	// assert
	assert.Equal(t, expected, *actual)
}

func Test_GetFinalize_withOutReference(t *testing.T) {
	// arrange
	bean := referenceBean{}

	// action
	actual := bean.GetFinalize()

	// assert
	assert.Nil(t, actual)
}

func Test_GetInit_withReference(t *testing.T) {
	// arrange
	expected := "a_method"
	mock := new(BeanMock)
	mock.On("GetInit").Return(&expected)
	bean := referenceBean{
		reference: mock,
	}

	// action
	actual := bean.GetInit()

	// assert
	assert.Equal(t, expected, *actual)
}

func Test_GetInit_withOutReference(t *testing.T) {
	// arrange
	bean := referenceBean{}

	// action
	actual := bean.GetInit()

	// assert
	assert.Nil(t, actual)
}

func Test_GetReference(t *testing.T) {
	// arrange
	expected := new(BeanMock)
	bean := referenceBean{
		reference: expected,
	}

	// action
	actual := bean.GetReference()

	// assert
	assert.Equal(t, expected, actual)
}

func Test_GetProperty_withReference(t *testing.T) {
	// arrange
	expected := []BeanI{nil}
	mock := new(BeanMock)
	mock.On("GetProperty", "name").Return(expected)
	bean := referenceBean{
		reference: mock,
	}

	// action
	actual := bean.GetProperty("name")

	// assert
	assert.Equal(t, expected, actual)
}

func Test_GetProperty_withOutReference(t *testing.T) {
	// arrange
	bean := referenceBean{}

	// action
	actual := bean.GetProperty("")

	// assert
	assert.Empty(t, actual)
}

func Test_GetProperties_withReference(t *testing.T) {
	// arrange
	expected := map[string][]BeanI{
		"name": []BeanI{nil},
	}
	mock := new(BeanMock)
	mock.On("GetProperties").Return(expected)
	bean := referenceBean{
		reference: mock,
	}

	// action
	actual := bean.GetProperties()

	// assert
	assert.Equal(t, expected, actual)
}

func Test_GetProperties_withOutReference(t *testing.T) {
	// arrange
	bean := referenceBean{}

	// action
	actual := bean.GetProperties()

	// assert
	assert.Empty(t, actual)
}

func Test_GetType_withReference(t *testing.T) {
	// arrange
	expected := reflect.TypeOf("")
	mock := new(BeanMock)
	mock.On("GetType").Return(expected)
	bean := referenceBean{
		reference: mock,
	}

	// action
	actual := bean.GetType()

	// assert
	assert.Equal(t, expected, actual)
}

func Test_GetType_withOutReference(t *testing.T) {
	// arrange
	bean := referenceBean{}

	// action
	actual := bean.GetType()

	// assert
	assert.Nil(t, actual)
}

func Test_SetReference(t *testing.T) {
	// arrange
	expected := new(BeanMock)
	bean := referenceBean{}

	// action
	bean.SetReference(expected)

	// assert
	assert.Equal(t, expected, bean.reference)
}
