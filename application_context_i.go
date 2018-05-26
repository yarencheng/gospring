package gospring

// ApplicationContextI is an interface to management beans.
type ApplicationContextI interface {

	// Accuire a bean from its ID.
	GetBean(id string) (interface{}, error)

	// A destory function of this instance
	Finalize() error
}
