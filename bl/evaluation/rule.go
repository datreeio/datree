package evaluation

type Rule struct {
	Identifier         string
	Name               string
	MessageOnFailure   string
	OccurrencesDetails []OccurrenceDetails
}

func (rp *Rule) GetOccurrencesCount() int {
	count := 0
	for _, occurrence := range rp.OccurrencesDetails {
		count += occurrence.Occurrences
	}
	return count
}

type OccurrenceDetails struct {
	MetadataName string `yaml:"metadataName" json:"metadataName" xml:"metadataName"`
	Kind         string `yaml:"kind" json:"kind" xml:"kind"`
	Occurrences  int    `yaml:"occurrences" json:"occurrences" xml:"occurrences"`
}
