package responses

import (
	"time"

	"github.com/gofrs/uuid"
)

type Privacy struct {
	ID        uuid.UUID `json:"id"`
	UserID    string    `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
}
