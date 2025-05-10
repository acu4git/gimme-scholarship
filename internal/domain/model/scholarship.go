package model

type Scholarship struct {
	ID             int64
	Name           string
	Address        string
	Targets        []string
	TargetDetail   string
	AmountDetail   string
	TypeDetail     string
	CapacityDetail string
	DeadlineDetail string
	ContactPoint   string
	Remark         string
	PostingDate    string
}
