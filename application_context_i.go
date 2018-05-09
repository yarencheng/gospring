package refactor

type ApplicationContextI interface {
	GetBean(id string) (interface{}, error)
	Finalize() error
}
