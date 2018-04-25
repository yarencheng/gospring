package beans

type PropertyMetaData struct {
	name      string
	reference string
}

func NewPropertyMetaData(name string, reference string) *PropertyMetaData {
	return &PropertyMetaData{
		name:      name,
		reference: reference,
	}
}

func (meta *PropertyMetaData) GetName() string {
	return meta.name
}

func (meta *PropertyMetaData) GetReference() string {
	return meta.reference
}
