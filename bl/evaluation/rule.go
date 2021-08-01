package evaluation

import "github.com/datreeio/datree/pkg/cliClient"

type Rule struct {
	ID             int
	Name           string
	FailSuggestion string
	Count          int
	Matches		   []*cliClient.Match
}

func (rp *Rule) IncrementCount() {
	rp.Count++
}
