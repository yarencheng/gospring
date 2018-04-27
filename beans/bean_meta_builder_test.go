package beans

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// MockPropertyMetaData is a mock of PropertyMetaData interface
type MockPropertyMetaData struct {
	mock.Mock
}

func (m *MockPropertyMetaData) GetName() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockPropertyMetaData) GetReference() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockPropertyMetaData) GetBean() BeanMetaData {
	args := m.Called()
	return args.Get(0).(BeanMetaData)
}

func (m *MockPropertyMetaData) IsReference() bool {
	args := m.Called()
	return args.Bool(0)
}

type bean_TestSuite struct {
	suite.Suite
}

func Test_bean(t *testing.T) {
	suite.Run(t, new(bean_TestSuite))
}

func (suite *bean_TestSuite) SetupTest() {

}

func (suite *bean_TestSuite) Test_GetProperties() {
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	expect := make([]PropertyMetaData, 1)
	expect[0] = &MockPropertyMetaData{}

	bean := bean{
		properties: expect,
	}

	actual := bean.GetProperties()

	assert.ElementsMatch(suite.T(), expect, actual)
}

type beanMetasBuilder_TestSuite struct {
	suite.Suite
}

func Test_beanMetasBuilder(t *testing.T) {
	suite.Run(t, new(beanMetasBuilder_TestSuite))
}

func (suite *beanMetasBuilder_TestSuite) SetupTest() {

}

func (suite *beanMetasBuilder_TestSuite) Test_Build_emptyProperty() {

	beans, e := Beans().Build()

	assert.NoError(suite.T(), e)

	assert.Len(suite.T(), beans, 0)

}

func (suite *beanMetasBuilder_TestSuite) Test_Build_oneProperty() {

	beans, e := Beans(&beanMetaBuilder{}).Build()

	assert.NoError(suite.T(), e)

	assert.Len(suite.T(), beans, 1)

}
