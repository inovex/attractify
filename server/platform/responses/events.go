package responses

import (
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
)

type Event struct {
	ID             uuid.UUID       `json:"id"`
	OrganizationID uuid.UUID       `json:"organizationId"`
	Name           string          `json:"name"`
	Description    string          `json:"description"`
	Structure      json.RawMessage `json:"structure"`
	Properties     json.RawMessage `json:"properties"`
	CreatedAt      time.Time       `json:"createdAt"`
	UpdatedAt      time.Time       `json:"updatedAt"`
}

type EventOverview struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type EventProperty struct {
	Key  string `json:"key"`
	Type string `json:"type"`
}
