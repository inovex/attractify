package requests

import (
	"encoding/json"
)

type Track struct {
	UserID     string           `form:"userId" binding:"required,min=1"`
	Event      string           `form:"event" binding:"required,min=1"`
	Context    *json.RawMessage `json:"context" binding:"omitempty"`
	Properties *json.RawMessage `json:"properties" binding:"omitempty"`
}
