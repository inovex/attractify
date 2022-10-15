package requests

import (
	"encoding/json"

	"github.com/gofrs/uuid"
)

type Actions struct {
	UserID  string          `form:"userId" binding:"required,min=1"`
	Type    uuid.UUID       `form:"type" binding:"omitempty,min=1"`
	Tags    []string        `form:"tags" binding:"omitempty,dive,min=1"`
	Context json.RawMessage `form:"context" binding:"omitempty,min=1"`
}

type ActionsAct struct {
	ActionID   string           `form:"actionId" binding:"required,min=1"`
	UserID     string           `form:"userId" binding:"required,min=1"`
	Event      string           `form:"event" binding:"required,oneof=show hide decline accept"`
	Context    *json.RawMessage `json:"context" binding:"omitempty"`
	Properties *json.RawMessage `json:"properties" binding:"omitempty"`
}
