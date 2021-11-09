package requests

type Privacy struct {
	UserID string `json:"userId" validate:"required,min=1"`
	Email  string `json:"email" validate:"omitempty,email"`
}
