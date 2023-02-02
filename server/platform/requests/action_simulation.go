package requests

import (
	"encoding/json"
)

type ActionSimulation struct {
	User ActionSimulationUser `json:"user" binding:"required"`
}

type ActionSimulationUser struct {
	UserID         string          `json:"userId"`
	Channel        string          `json:"channel"`
	ActionID       string          `json:"actionId"`
	ComputedTraits json.RawMessage `json:"computedTraits"`
	CustomTraits   json.RawMessage `json:"customTraits"`
	Context        json.RawMessage `json:"context"`
}
