package responses

import (
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
)

type ComputedTrait struct {
	ID          uuid.UUID       `json:"id"`
	Name        string          `json:"name"`
	Key         string          `json:"key"`
	Type        string          `json:"type"`
	EventID     uuid.UUID       `json:"eventId"`
	Conditions  json.RawMessage `json:"conditions"`
	Properties  json.RawMessage `json:"properties"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
	RefreshedAt time.Time       `json:"refreshedAt"`
}

type ComputedTraitRefresh struct {
	Count         int `json:"count"`
	ExecutionTime int `json:"executionTime"`
}
