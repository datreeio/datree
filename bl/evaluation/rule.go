package evaluation

import "github.com/datreeio/datree/pkg/cliClient"

type Rule struct {
	ID             int
	Name           string
	FailSuggestion string
	Matches        []*cliClient.Match
}

func (rp *Rule) GetCount() int {
	return len(rp.Matches)
}
