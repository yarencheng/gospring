package refactor

type ReferenceBeanI interface {
	ID(id string) ReferenceBeanI
	GetID() *string
}
