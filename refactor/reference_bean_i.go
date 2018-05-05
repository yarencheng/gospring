package refactor

type ReferenceBeanI interface {
	GetReference() BeanI
	SetReference(BeanI)
}
