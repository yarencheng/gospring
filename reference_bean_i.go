package gospring

type ReferenceBeanI interface {
	GetReference() BeanI
	SetReference(bean BeanI)
}
