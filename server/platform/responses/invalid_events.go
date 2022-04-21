package responses

import (
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
)

type InvalidEvent struct {
	ID             uuid.UUID       `json:"id"`
	OrganizationID uuid.UUID       `json:"organizationId"`
	Channel        string          `json:"channel"`
	Name           string          `json:"name"`
	Properties     json.RawMessage `json:"properties"`
	Context        json.RawMessage `json:"context"`
	Type           string          `json:"type"`
	CreatedAt      time.Time       `json:"createdAt"`
}

type InvalidEventOverview struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type InvalidEventProperty struct {
	Key  string `json:"key"`
	Type string `json:"type"`
}
