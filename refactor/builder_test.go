package refactor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Beans_notAnyBean(t *testing.T) {
	// arrange

	// action
	beans := Beans(
		123,
	)

	// assert
	assert.Implements(t, new(ValueBeanI), beans[0])
}

func Test_Beans_structBean(t *testing.T) {
	// arrange

	// action
	beans := Beans(
		new(structBean),
	)

	// assert
	assert.Implements(t, new(StructBeanI), beans[0])
}

func Test_Beans_referenceBean(t *testing.T) {
	// arrange

	// action
	beans := Beans(
		new(referenceBean),
	)

	// assert
	assert.Implements(t, new(ReferenceBeanI), beans[0])
}

func Test_Beans_valueBean(t *testing.T) {
	// arrange

	// action
	beans := Beans(
		new(valueBean),
	)

	// assert
	assert.Implements(t, new(ValueBeanI), beans[0])
}
