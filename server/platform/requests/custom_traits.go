package requests

import (
	"encoding/json"
)

type CustomTraitsUpsert struct {
	Structure json.RawMessage `json:"structure" binding:"omitempty"`
}
