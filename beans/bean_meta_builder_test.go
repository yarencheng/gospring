package beans

import (
	"reflect"
	"testing"
)

func Test_bean_GetProperties(t *testing.T) {

	expect := make([]PropertyMetaData, 1)
	expect[0] = &propertyBean{}

	type A struct{ i int }

	expect[0] = &propertyBean{}
	// expect[0] = &A{}

	reflect.TypeOf("")

	//vean := bean{]}
}
