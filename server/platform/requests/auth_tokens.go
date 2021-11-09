package requests

type AuthToken struct {
	Channel string `json:"channel" binding:"required,min=1"`
}
