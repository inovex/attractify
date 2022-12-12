package requests

type typeProperty struct {
	Channels   []string `json:"channels" binding:"dive,min=1"`
	Type       string   `json:"type" binding:"required,oneof=text custom_trait computed_trait"`
	Name       string   `json:"name" binding:"required,min=1"`
	Value      string   `json:"value" binding:"omitempty,min=1"`
	SourceKey  string   `json:"sourceKey" binding:"omitempty,min=1"`
	SourceType string   `json:"sourceType" binding:"omitempty,min=1"`
}

type ActionTypeCreate struct {
	Name       string     `json:"name" binding:"required,min=1"`
	Version    int        `json:"version" binding:"required,min=1"`
	Properties []property `json:"properties" binding:"dive"`
}

type ActionTypeUpdate struct {
	Properties []property `json:"properties" binding:"dive"`
}
