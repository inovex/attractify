package responses

import (
	"time"

	"github.com/gofrs/uuid"
)

type AuthToken struct {
	ID        uuid.UUID `json:"id"`
	Channel   string    `json:"channel"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"createdAt"`
}
