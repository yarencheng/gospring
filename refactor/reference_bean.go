package refactor

func Ref(id string) ReferenceBeanI {
	return &referenceBean{
		id: id,
	}
}

type referenceBean struct {
	id string
}

func (bean *referenceBean) ID(id string) ReferenceBeanI {
	return bean
}

func (bean *referenceBean) GetID() *string {
	return &bean.id
}
