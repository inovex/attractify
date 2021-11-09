package responses

import (
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
)

type CampaignWebhookTest struct {
	StatusCode int    `json:"statusCode"`
	Body       string `json:"body"`
	Signature  string `json:"signature"`
	Error      string `json:"error"`
}

type Action struct {
	ID             uuid.UUID       `json:"id"`
	OrganizationID uuid.UUID       `json:"organizationId"`
	Name           string          `json:"name"`
	Type           string          `json:"type"`
	State          string          `json:"state"`
	Tags           json.RawMessage `json:"tags"`
	Properties     json.RawMessage `json:"properties"`
	Targeting      json.RawMessage `json:"targeting"`
	Capping        json.RawMessage `json:"capping"`
	Hooks          json.RawMessage `json:"hooks"`
	TestUsers      json.RawMessage `json:"testUsers"`
	CreatedAt      time.Time       `json:"createdAt"`
	UpdatedAt      time.Time       `json:"updatedAt"`
}
