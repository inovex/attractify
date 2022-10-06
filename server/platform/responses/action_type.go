package responses

import (
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
)

type ActionType struct {
	ID             uuid.UUID       `json:"id"`
	OrganizationID uuid.UUID       `json:"organizationId"`
	Name           string          `json:"name"`
	Version        int             `json:"version"`
	Properties     json.RawMessage `json:"properties"`
	Archived       bool            `json:"archived"`
	CreatedAt      time.Time       `json:"createdAt"`
}
