package responses

import (
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
)

type Context struct {
	ID             uuid.UUID       `json:"id"`
	OrganizationID uuid.UUID       `json:"organizationId"`
	Channel        string          `json:"channel"`
	Structure      json.RawMessage `json:"structure"`
	Properties     json.RawMessage `json:"properties"`
	CreatedAt      time.Time       `json:"createdAt"`
	UpdatedAt      time.Time       `json:"updatedAt"`
}

type ContextProperty struct {
	Key  string `json:"key"`
	Type string `json:"type"`
}
