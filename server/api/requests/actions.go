package requests

import (
	"encoding/json"
)

type Actions struct {
	UserID   string          `form:"userId" json:"userId" binding:"required,min=1"`
	TypeName string          `form:"type" json:"type" binding:"omitempty,min=1"`
	Tags     []string        `form:"tags" json:"tags" binding:"omitempty,dive,min=1"`
	Context  json.RawMessage `json:"context" binding:"omitempty,min=1"`
}

type ActionsAct struct {
	ActionID   string           `json:"actionId" binding:"required,min=1"`
	UserID     string           `json:"userId" binding:"required,min=1"`
	Event      string           `json:"event" binding:"required,oneof=show hide decline accept"`
	Context    *json.RawMessage `json:"context" binding:"omitempty"`
	Properties *json.RawMessage `json:"properties" binding:"omitempty"`
}
