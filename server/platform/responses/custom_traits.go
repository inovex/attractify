package responses

import (
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
)

type CustomTraits struct {
	OrganizationID uuid.UUID       `json:"organizationId"`
	Structure      json.RawMessage `json:"structure"`
	Properties     json.RawMessage `json:"properties"`
	CreatedAt      time.Time       `json:"createdAt"`
	UpdatedAt      time.Time       `json:"updatedAt"`
}

type CustomTraitsProperty struct {
	Key  string `json:"key"`
	Type string `json:"type"`
}
