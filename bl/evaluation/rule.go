package evaluation

type Rule struct {
	ID             int
	Name           string
	FailSuggestion string
	Count          int
}

func (rp *Rule) IncrementCount() {
	rp.Count++
}
