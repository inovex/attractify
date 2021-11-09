package requests

type UserCreate struct {
	Email string `json:"email" validate:"required,email"`
	Role  string `json:"role" validate:"required,oneof=admin marketeer"`
}

type UserActivate struct {
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name" validate:"required,min=1"`
}

type UserSession struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type UserUpdate struct {
	Email       string `json:"email" validate:"omitempty,email"`
	OldPassword string `json:"oldPassword" validate:"omitempty,min=8"`
	NewPassword string `json:"newPassword" validate:"omitempty,min=8,nefield=OldPassword,required_with=OldPassword"`
	Name        string `json:"name" validate:"omitempty,min=1"`
}

type UserUpdateRole struct {
	Role string `json:"role" validate:"required,oneof=admin marketeer"`
}

type UserResetPassword struct {
	Email string `json:"email" validate:"required,email"`
}

type UserUpdatePassword struct {
	Password string `json:"password" validate:"required,min=8"`
}
