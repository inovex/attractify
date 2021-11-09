package responses

import (
	"time"

	"github.com/gofrs/uuid"
)

type Channel struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Key       string    `json:"key"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
