package requests

type Reactions struct {
	ActionID     string `form:"actionId" binding:"omitempty"`
	UserID       string `form:"userId" binding:"omitempty"`
	Events       string `form:"events" binding:"omitempty"`
	Page         int    `form:"page" binding:"omitempty,min=1"`
	ItemsPerPage int    `form:"itemsPerPage" binding:"omitempty,min=1,max=50"`
	Start        string `form:"start" binding:"omitempty,datetime=2006-01-02"`
	End          string `form:"end" binding:"omitempty,datetime=2006-01-02"`
}
