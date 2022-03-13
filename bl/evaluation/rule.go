package evaluation

type Rule struct {
	Identifier         string
	Name               string
	MessageOnFailure   string
	OccurrencesDetails []OccurrenceDetails
}

func (rp *Rule) GetCount() int {
	return len(rp.OccurrencesDetails)
}

type OccurrenceDetails struct {
	MetadataName string `yaml:"metadataName" json:"metadataName" xml:"metadataName"`
	Kind         string `yaml:"kind" json:"kind" xml:"kind"`
}
