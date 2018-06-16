package application_context

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yarencheng/gospring/v1"
)

func Test_New_v1Bean(t *testing.T) {
	// arrange
	config := v1.Bean{
		ID:   "aa",
		Type: reflect.TypeOf(""),
	}

	// action
	_, err := New(config)

	// assert
	assert.NoError(t, err)
}
