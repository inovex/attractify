package requests

import (
	"encoding/json"
)

type Identify struct {
	UserID         string           `json:"userId" binding:"required,min=1"`
	PreviousUserID string           `json:"previousUserId" binding:"omitempty,min=1"`
	Type           string           `json:"type" binding:"required,min=1"`
	IsAnonymous    bool             `json:"isAnonymous" binding:"omitempty"`
	Traits         *json.RawMessage `json:"traits" binding:"omitempty"`
}
