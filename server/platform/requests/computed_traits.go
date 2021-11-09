package requests

import (
	"encoding/json"

	"github.com/gofrs/uuid"
)

type ComputedTraitCreate struct {
	Name       string          `json:"name" binding:"required,min=1"`
	Key        string          `json:"key" binding:"required,min=1"`
	Type       string          `json:"type" binding:"required,oneof=count_events aggregation most_frequent first_event last_event unique_list unique_list_count"`
	EventID    uuid.UUID       `json:"eventId" binding:"required,min=1"`
	Conditions json.RawMessage `json:"conditions" binding:"required"`
	Properties json.RawMessage `json:"properties" binding:"required"`
}
