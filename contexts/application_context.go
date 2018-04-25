package contexts

type ApplicatoinContext interface {
	GetBean(name string) (interface{}, error)
	Start() error
	Stop() error
}
