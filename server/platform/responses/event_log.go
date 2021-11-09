package responses

import (
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
)

type EventLog struct {
	ID         uuid.UUID       `json:"id"`
	EventID    uuid.UUID       `json:"eventId"`
	Channel    string          `json:"channel"`
	Context    json.RawMessage `json:"context"`
	Properties json.RawMessage `json:"properties"`
	CreatedAt  time.Time       `json:"createdAt"`
}

type EventLogList struct {
	Events []EventLog `json:"events"`
	Count  int        `json:"count"`
}
