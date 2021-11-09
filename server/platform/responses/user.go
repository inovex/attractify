package responses

import (
	"time"

	"github.com/gofrs/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UserSession struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
	Token string `json:"token"`
}
