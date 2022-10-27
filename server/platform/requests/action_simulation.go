package requests

import (
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
)

// `json:"name" binding:"required,min=1"`

type ActionSimulation struct {
	User ActionSimulationUser `json:"user" binding:"required"`
}

type ActionSimulationUser struct {
	UserID         uuid.UUID       `json:"id"`
	Channel        string          `json:"channel"`
	Time           time.Time       `json:"time"`
	ComputedTraits json.RawMessage `json:"computedTraits" binding:"required"`
	CustomTraits   json.RawMessage `json:"customTraits" binding:"required"`
	Context        json.RawMessage `json:"context" binding:"required"`
}
