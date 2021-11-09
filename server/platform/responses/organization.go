package responses

type Organization struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Timezone string `json:"timezone"`
}

type OrganizationToken struct {
	Token string `json:"token"`
}

type OrganizationKey struct {
	Key string `json:"key"`
}
