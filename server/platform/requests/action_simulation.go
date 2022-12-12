package requests

import (
	"encoding/json"

	"github.com/gofrs/uuid"
)

type ActionSimulation struct {
	User ActionSimulationUser `json:"user" binding:"required"`
}

type ActionSimulationUser struct {
	UserID   uuid.UUID `json:"userId"`
	Channel  string    `json:"channel"`
	ActionID string    `json:"actionId"`
	ComputedTraits json.RawMessage `json:"computedTraits"`
	CustomTraits   json.RawMessage `json:"customTraits"`
	Context        json.RawMessage `json:"context"`
}
