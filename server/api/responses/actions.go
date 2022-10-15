package responses

import (
	"encoding/json"

	"github.com/gofrs/uuid"
)

type Action struct {
	ID         uuid.UUID              `json:"id"`
	Type       string                 `json:"type"`
	Version    int                    `json:"version"`
	Tags       json.RawMessage        `json:"tags"`
	Properties map[string]interface{} `json:"properties"`
}
