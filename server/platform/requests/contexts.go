package requests

import (
	"encoding/json"
)

type Context struct {
	Channel   string          `json:"channel" binding:"required,min=1"`
	Structure json.RawMessage `json:"structure" binding:"omitempty"`
}
