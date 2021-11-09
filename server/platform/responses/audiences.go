package responses

import (
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
)

type Audience struct {
	ID               uuid.UUID       `json:"id"`
	Name             string          `json:"name"`
	Description      string          `json:"description"`
	IncludeAnonymous bool            `json:"includeAnonymous"`
	Events           json.RawMessage `json:"events"`
	Traits           json.RawMessage `json:"traits"`
	ProfileCount     int             `json:"profileCount"`
	CreatedAt        time.Time       `json:"createdAt"`
	UpdatedAt        time.Time       `json:"updatedAt"`
	RefreshedAt      time.Time       `json:"refreshedAt"`
}

type AudienceProfile struct {
	ID             uuid.UUID       `json:"id"`
	UserID         string          `json:"userId"`
	CustomTraits   json.RawMessage `json:"customTraits"`
	ComputedTraits json.RawMessage `json:"computedTraits"`
	CreatedAt      time.Time       `json:"createdAt"`
}

type AudiencePreview struct {
	ExecutionTime int64             `json:"executionTime"`
	Profiles      []AudienceProfile `json:"profiles"`
}

type AudienceRefresh struct {
	Count         int `json:"count"`
	ExecutionTime int `json:"executionTime"`
}
