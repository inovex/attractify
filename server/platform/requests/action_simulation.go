package requests

import (
	"encoding/json"

	"github.com/gofrs/uuid"
)

// `json:"name" binding:"required,min=1"`

type ActionSimulation struct {
	User ActionSimulationUser `json:"user" binding:"required"`
}

type ActionSimulationUser struct {
	UserID   uuid.UUID `json:"userId"`
	Channel  string    `json:"channel"`
	ActionID string    `json:"actionId"`
	//Time           int64           `json:"time"`
	ComputedTraits json.RawMessage `json:"computedTraits"`
	CustomTraits   json.RawMessage `json:"customTraits"`
	Context        json.RawMessage `json:"context"`
}
