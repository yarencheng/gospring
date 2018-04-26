package beans

type PropertyMetaData struct {
	name      string
	reference string
	value     string
}

func NewPropertyMetaData(name string, reference string, value string) *PropertyMetaData {
	return &PropertyMetaData{
		name:      name,
		reference: reference,
		value:     value,
	}
}

func (meta *PropertyMetaData) GetName() string {
	return meta.name
}

func (meta *PropertyMetaData) GetReference() string {
	return meta.reference
}

func (meta *PropertyMetaData) GetValue() string {
	return meta.value
}
