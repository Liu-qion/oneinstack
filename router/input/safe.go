package input

type IptablesRuleParam struct {
	Q      string `json:"q"`
	ID     int64  `json:"id"`
	Target string `json:"target"`
	Page
}
