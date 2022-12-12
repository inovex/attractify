package requests

type Analyze struct {
	ActionID string `json:"actionId" validate:"required,uuid"`
	Start    string `json:"start" validate:"omitempty,datetime=2006-01-02"`
	End      string `json:"end" validate:"omitempty,datetime=2006-01-02"`
	Interval string `json:"interval" validate:"omitempty,oneof=year month day week_day hour"`
}
