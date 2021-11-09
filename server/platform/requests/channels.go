package requests

type Channel struct {
	Name string `json:"name" binding:"required,min=1"`
	Key  string `json:"key" binding:"required,min=1"`
}
