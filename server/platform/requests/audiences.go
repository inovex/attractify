package requests

import (
	"encoding/json"
)

type AudienceCreate struct {
	Name             string          `json:"name" binding:"required,min=1"`
	Description      string          `json:"description" binding:"required,min=1"`
	IncludeAnonymous bool            `json:"includeAnonymous" binding:"omitempty"`
	Events           json.RawMessage `json:"events" binding:"omitempty"`
	Traits           json.RawMessage `json:"traits" binding:"omitempty"`
}
