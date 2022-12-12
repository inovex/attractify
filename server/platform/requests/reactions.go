package requests

type Reactions struct {
	ActionID     string `json:"actionId" binding:"omitempty"`
	UserID       string `json:"userId" binding:"omitempty"`
	Events       string `json:"events" binding:"omitempty"`
	Page         int    `json:"page" binding:"omitempty,min=1"`
	ItemsPerPage int    `json:"itemsPerPage" binding:"omitempty,min=1,max=50"`
	Start        string `json:"start" binding:"omitempty,datetime=2006-01-02"`
	End          string `json:"end" binding:"omitempty,datetime=2006-01-02"`
}
