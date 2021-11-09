package requests

type OrganizationCreate struct {
	Email            string `json:"email" validate:"required,email"`
	Password         string `json:"password" validate:"required,min=8"`
	OrganizationName string `json:"organization_name" validate:"required,min=1"`
	Name             string `json:"name" validate:"required,min=1"`
	Timezone         string `json:"timezone" validate:"required,min=1"`
}

type OrganizationUpdate struct {
	Email    string `json:"email" validate:"omitempty,email"`
	Name     string `json:"name" validate:"omitempty,min=1"`
	Timezone string `json:"timezone" validate:"omitempty,min=1"`
}

type OrganizationKey struct {
	Password string `json:"password" validate:"required,min=8"`
}
