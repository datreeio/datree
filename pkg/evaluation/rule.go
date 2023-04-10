package evaluation

import "github.com/datreeio/datree/pkg/cliClient"

type Rule struct {
	Identifier         string
	Name               string
	MessageOnFailure   string
	DocumentationUrl   string
	OccurrencesDetails []OccurrenceDetails
}

func (rp *Rule) GetFailedOccurrencesCount() int {
	count := 0
	for _, occurrence := range rp.OccurrencesDetails {
		if !occurrence.IsSkipped {
			count += occurrence.Occurrences
		}
	}
	return count
}

type OccurrenceDetails struct {
	MetadataName      string                      `yaml:"metadataName" json:"metadataName" xml:"metadataName"`
	Kind              string                      `yaml:"kind" json:"kind" xml:"kind"`
	SkipMessage       string                      `yaml:"skipMessage" json:"skipMessage" xml:"skipMessage"`
	Occurrences       int                         `yaml:"occurrences" json:"occurrences" xml:"occurrences"`
	IsSkipped         bool                        `yaml:"isSkipped" json:"isSkipped" xml:"isSkipped"`
	FailureLocations  []cliClient.FailureLocation `yaml:"failureLocations" json:"failureLocations" xml:"failureLocations"`
	ValidationResults []string                    `yaml:"validationResults" json:"validationResults" xml:"validationResults"`
}
