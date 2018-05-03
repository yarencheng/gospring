package refactor

type ApplicationContextI interface {
	GetBean() (interface{}, error)
	Finalize() error
}
