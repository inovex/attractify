package requests

import (
	"encoding/json"
)

type Identify struct {
	UserID         string           `form:"userId" binding:"required,min=1"`
	PreviousUserID string           `form:"previousUserId" binding:"omitempty,min=1"`
	Type           string           `form:"type" binding:"required,min=1"`
	IsAnonymous    bool             `form:"isAnonymous" binding:"omitempty"`
	Traits         *json.RawMessage `form:"traits" binding:"omitempty"`
}
