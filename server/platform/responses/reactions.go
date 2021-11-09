package responses

import (
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
)

type Reaction struct {
	ID         uuid.UUID       `json:"id"`
	Event      string          `json:"event"`
	Channel    string          `json:"channel"`
	Context    json.RawMessage `json:"context"`
	Properties json.RawMessage `json:"properties"`
	Result     json.RawMessage `json:"result"`
	CreatedAt  time.Time       `json:"createdAt"`
}

type ReactionList struct {
	Reactions []Reaction `json:"reactions"`
	Count     int        `json:"count"`
}
