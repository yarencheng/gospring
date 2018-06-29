package application_context

type ApplicationContextI interface {
	GetByID(id string) (interface{}, error)
}
