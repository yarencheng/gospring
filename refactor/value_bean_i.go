package refactor

type ValueBeanI interface {
	GetID() *string
	GetValue() interface{}
}
