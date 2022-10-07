package requests

import (
	"github.com/jmoiron/sqlx/types"
)

type property struct {
	Channels   []string `json:"channels" binding:"dive,min=1"`
	Type       string   `json:"type" binding:"required,oneof=text custom_trait computed_trait"`
	Name       string   `json:"name" binding:"required,min=1"`
	Value      string   `json:"value" binding:"omitempty,min=1"`
	SourceKey  string   `json:"sourceKey" binding:"omitempty,min=1"`
	SourceType string   `json:"sourceType" binding:"omitempty,min=1"`
}

type dateTime struct {
	Date *string `json:"date" binding:"omitempty,datetime=2006-01-02"`
	Time *string `json:"time" binding:"omitempty,datetime=15:04"`
}

type traitConditions struct {
	Source   string      `json:"source" binding:"required,oneof=custom computed"`
	Key      string      `json:"key" binding:"required,min=1"`
	Type     string      `json:"type" binding:"required,min=1"`
	Operator string      `json:"operator" binding:"required,min=1"`
	Value    interface{} `json:"value" binding:"omitempty,min=1"`
}

type contextConditions struct {
	Channel  string      `json:"channel" binding:"required,min=1"`
	Key      string      `json:"key" binding:"required,min=1"`
	Type     string      `json:"type" binding:"required,min=1"`
	Operator string      `json:"operator" binding:"required,min=1"`
	Value    interface{} `json:"value" binding:"omitempty,min=1"`
}

type targeting struct {
	Audiences         []string            `json:"audiences" binding:"dive,uuid"`
	Channels          []string            `json:"channels" binding:"dive,min=1"`
	TraitConditions   []traitConditions   `json:"traitConditions" binding:"dive"`
	ContextConditions []contextConditions `json:"contextConditions" binding:"dive"`
	Start             dateTime            `json:"start" binding:"omitempty"`
	End               dateTime            `json:"end" binding:"omitempty"`
}

type capping struct {
	Channels []string `json:"channels" binding:"dive,min=1"`
	Event    string   `json:"event" binding:"required,oneof=shown hidden declined accepted"`
	Group    string   `json:"group" binding:"required,oneof=all user"`
	Count    int      `json:"count" binding:"required,min=0"`
	Within   int      `json:"within" binding:"omitempty,min=0"`
}

type webhook struct {
	Events       []string `json:"events" binding:"dive,oneof=close decline accept"`
	URL          string   `json:"url" binding:"omitempty,min=10"`
	Properties   string   `json:"properties" binding:"omitempty,min=1"`
	ReturnResult bool     `json:"returnResult" binding:"omitempty"`
}

type hook struct {
	Channels   []string               `json:"channels" binding:"required,dive,min=1"`
	Event      string                 `json:"event" binding:"required,oneof=shown hidden declined accepted"`
	Type       string                 `json:"type" binding:"required,oneof=execute_webhook track_event"`
	Properties map[string]interface{} `json:"properties" binding:"omitempty"`
}

type testUser struct {
	Channels      []string `json:"channels" binding:"dive,min=1"`
	UserID        string   `json:"userId" binding:"required,min=1"`
	Description   string   `json:"description" binding:"omitempty"`
	SkipTargeting bool     `json:"skipTargeting" binding:"omitempty"`
}

type ActionCreate struct {
	Type           string     `json:"type" binding:"required,min=1"`
	Version        int        `json:"version" binding:"required,min=1"`
	Name           string     `json:"name" binding:"required,min=1"`
	State          string     `json:"state" binding:"required,oneof=inactive staging active"`
	Tags           []string   `json:"tags" binding:"dive,min=1"`
	Properties     []property `json:"properties" binding:"dive"`
	TypeProperties []property `json:"typeProperties" binding:"dive"`
	Targeting      targeting  `json:"targeting"`
	Capping        []capping  `json:"capping"`
	Hooks          []hook     `json:"hooks"`
	TestUsers      []testUser `json:"testUsers"`
}

type ActionState struct {
	State string `json:"state" binding:"required,oneof=inactive staging active"`
}

type ActionWebhookTest struct {
	Event      string         `json:"event" binding:"required,oneof=show close dismiss accept"`
	UserID     string         `json:"userId" binding:"required,min=1"`
	Channel    string         `json:"channel" binding:"required,min=1"`
	Properties types.JSONText `json:"properties" binding:"omitempty"`
}
