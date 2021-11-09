package requests

type Analyze struct {
	ActionID string `form:"actionId" validate:"required,uuid"`
	Start    string `form:"start" validate:"omitempty,datetime=2006-01-02"`
	End      string `form:"end" validate:"omitempty,datetime=2006-01-02"`
	Interval string `form:"interval" validate:"omitempty,oneof=year month day week_day hour"`
}
