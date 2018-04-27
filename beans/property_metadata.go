package beans

type PropertyMetaData interface {
	GetName() string
	GetReference() string
	GetBean() BeanMetaData
	IsReference() bool
}

type PropertyMetaData_old struct {
	name      string
	reference string
	value     string
}

func NewPropertyMetaData_old(name string, reference string, value string) *PropertyMetaData_old {
	return &PropertyMetaData_old{
		name:      name,
		reference: reference,
		value:     value,
	}
}

func (meta *PropertyMetaData_old) GetName() string {
	return meta.name
}

func (meta *PropertyMetaData_old) GetReference() string {
	return meta.reference
}

func (meta *PropertyMetaData_old) GetValue() string {
	return meta.value
}
