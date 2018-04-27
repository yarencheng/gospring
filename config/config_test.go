package config

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/suite"
)

type config_suite struct { // suite for level 1
	suite.Suite
}

func Test_sss(t *testing.T) {
	suite.Run(t, new(config_suite))
}

func (s *config_suite) SetupTest() {

}

func (s *config_suite) Test_Config_empty() {
	v := Config().Validate()
	assert.True(s.T(), v)
}

func (s *config_suite) Test_Config_GetBean() {

	type beanStruct struct{ i int }
	id := "iiiidddd"

	config := Config(
		Bean(id, reflect.TypeOf(beanStruct{})),
	)

	ctx, e1 := ApplicationContext(config)
	assert.NoError(s.T(), e1)

	bean, e2 := ctx.GetBean(id)
	assert.NoError(s.T(), e2)

	_, ok := bean.(*beanStruct)
	assert.True(s.T(), ok)
}
