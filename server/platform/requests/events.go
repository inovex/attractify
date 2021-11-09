package requests

import (
	"encoding/json"
)

type EventCreate struct {
	Name        string          `json:"name" binding:"omitempty"`
	Description string          `json:"description" binding:"omitempty"`
	Structure   json.RawMessage `json:"structure" binding:"omitempty"`
}
