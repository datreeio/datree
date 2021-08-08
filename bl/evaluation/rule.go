package evaluation

type Rule struct {
	ID                 int
	Name               string
	FailSuggestion     string
	OccurrencesDetails []OccurrenceDetails
}

func (rp *Rule) GetCount() int {
	return len(rp.OccurrencesDetails)
}

type OccurrenceDetails struct {
	MetadataName string
	Kind         string
}
