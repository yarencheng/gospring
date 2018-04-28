package config

import (
	"reflect"
	"testing"
	"unsafe"

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

func (s *config_suite) Test_Config_GetBean_beanIsSingleton() {

	// arrange

	type beanStruct struct{ i int }
	id := "iiiidddd"

	config := Config(
		Bean(id, reflect.TypeOf(beanStruct{})).Singleton(), // singleton
	)

	// action

	var e error
	var ctx *applicationContext
	var i1, i2 interface{}
	var s1, s2 *beanStruct
	var ok bool

	ctx, e = ApplicationContext(config)
	assert.NoError(s.T(), e)

	i1, e = ctx.GetBean(id)
	assert.NoError(s.T(), e)

	i2, e = ctx.GetBean(id)
	assert.NoError(s.T(), e)

	s1, ok = i1.(*beanStruct)
	assert.True(s.T(), ok)

	s2, ok = i2.(*beanStruct)
	assert.True(s.T(), ok)

	// assert
	assert.Equal(s.T(), unsafe.Pointer(s1), unsafe.Pointer(s2))
}

func (s *config_suite) Test_Config_GetBean_beanIsPrototype() {

	// arrange

	type beanStruct struct{ i int }
	id := "iiiidddd"

	config := Config(
		Bean(id, reflect.TypeOf(beanStruct{})).Prototype(), // prototype
	)

	// action

	var e error
	var ctx *applicationContext
	var i1, i2 interface{}
	var s1, s2 *beanStruct
	var ok bool

	ctx, e = ApplicationContext(config)
	assert.NoError(s.T(), e)

	i1, e = ctx.GetBean(id)
	assert.NoError(s.T(), e)

	i2, e = ctx.GetBean(id)
	assert.NoError(s.T(), e)

	s1, ok = i1.(*beanStruct)
	assert.True(s.T(), ok)

	s2, ok = i2.(*beanStruct)
	assert.True(s.T(), ok)

	// assert
	assert.NotEqual(s.T(), unsafe.Pointer(s1), unsafe.Pointer(s2))
}

func (s *config_suite) Test_Config_GetBean_withIntPropertyInside() {

	// arrange

	type beanStruct struct{ I int } // with a int property
	id := "iiiidddd"

	config := Config(
		Bean(id, reflect.TypeOf(beanStruct{})).Prototype().With(
			Value("I", "123"),
		),
	)

	// action

	var e error
	var ctx *applicationContext
	var i interface{}
	var bean *beanStruct
	var ok bool

	ctx, e = ApplicationContext(config)
	assert.NoError(s.T(), e)

	i, e = ctx.GetBean(id)
	assert.NoError(s.T(), e)

	bean, ok = i.(*beanStruct)
	assert.True(s.T(), ok)

	// assert
	assert.Equal(s.T(), 123, bean.I)
}

func (s *config_suite) Test_Config_GetBean_withStringPropertyInside() {

	// arrange

	type beanStruct struct{ S string } // with string property
	id := "iiiidddd"

	config := Config(
		Bean(id, reflect.TypeOf(beanStruct{})).Prototype().With(
			Value("S", "a string property"),
		),
	)

	// action

	var e error
	var ctx *applicationContext
	var i interface{}
	var bean *beanStruct
	var ok bool

	ctx, e = ApplicationContext(config)
	assert.NoError(s.T(), e)

	i, e = ctx.GetBean(id)
	assert.NoError(s.T(), e)

	bean, ok = i.(*beanStruct)
	assert.True(s.T(), ok)

	// assert
	assert.Equal(s.T(), "a string property", bean.S)
}
