package refactor

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func Test_NewApplicationContext_empty(t *testing.T) {
	// arrange

	// action
	_, e := NewApplicationContext()

	// assert
	assert.Nil(t, e)
}

func Test_NewApplicationContext_ValueBeanI(t *testing.T) {
	// arrange
	beans := Beans(111)

	// action
	_, e := NewApplicationContext(beans...)
	require.Nil(t, e)

	// assert
	assert.Nil(t, e)
}

func Test_NewApplicationContext_ReferenceBeanI(t *testing.T) {
	// arrange
	beans := Beans(111)

	// action
	_, e := NewApplicationContext(beans...)

	// assert
	assert.Nil(t, e)
}

func Test_NewApplicationContext_ReferenceBeanI_loop(t *testing.T) {
	// arrange
	type beanStract struct{}
	beans := Beans(
		Bean(beanStract{}).
			ID("id_1").
			Property("aaa", Bean(beanStract{}).
				ID("id_2").
				Property("aaa", Ref("id_1")),
			),
	)

	// action
	_, e := NewApplicationContext(beans...)

	// assert
	assert.NotNil(t, e)
}

func Test_NewApplicationContext_StructBeanI(t *testing.T) {
	// arrange
	type beanStract struct{}
	beans := Beans(
		Bean(beanStract{}),
	)

	// action
	_, e := NewApplicationContext(beans...)

	// assert
	assert.Nil(t, e)
}

func Test_NewApplicationContext_StructBeanI_idIsDuplicated(t *testing.T) {
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

func Test_NewApplicationContext_StructBeanI_haveProperty(t *testing.T) {
	// arrange
	type beanStract struct{}
	beans := Beans(
		Bean(beanStract{}).
			ID("id").
			Property("a", 123),
	)

	// action
	_, e := NewApplicationContext(beans...)

	// assert
	assert.Nil(t, e)
}
