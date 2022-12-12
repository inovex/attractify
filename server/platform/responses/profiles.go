package responses

import (
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
)

type Profile struct {
	ID             uuid.UUID       `json:"id"`
	CustomTraits   json.RawMessage `json:"customTraits"`
	ComputedTraits json.RawMessage `json:"computedTraits"`
	CreatedAt      time.Time       `json:"createdAt"`
	UpdatedAt      time.Time       `json:"updatedAt"`
}

type IdentityWithTraits struct {
	ID             uuid.UUID       `json:"id"`
	Channel        string          `json:"channel"`
	Type           string          `json:"type"`
	UserID         string          `json:"userId"`
	CustomTraits   json.RawMessage `json:"customTraits"`
	ComputedTraits json.RawMessage `json:"computedTraits"`
	IsAnonymous    bool            `json:"isAnonymous"`
	CreatedAt      time.Time       `json:"createdAt"`
	UpdatedAt      time.Time       `json:"updatedAt"`
}

type Identity struct {
	ID          uuid.UUID `json:"id"`
	Channel     string    `json:"channel"`
	Type        string    `json:"type"`
	UserID      string    `json:"userId"`
	IsAnonymous bool      `json:"isAnonymous"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type TrackingEvent struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"organizationId"`
	ProfileID      uuid.UUID `json:"profileId"`
	Name           string    `json:"name"`
	Channel        string    `json:"channel"`
	Context        string    `json:"context"`
	Properties     string    `json:"properties"`
	CreatedAt      time.Time `json:"createdAt"`
}
