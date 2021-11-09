package privacy

import (
	"context"

	"attractify.io/platform/app"
	"github.com/gofrs/uuid"
)

type Locked struct {
	app            *app.App
	ctx            context.Context
	organizationID uuid.UUID
	userID         string
}

func NewLocked(app *app.App, organizationID uuid.UUID, userID string) *Locked {
	return &Locked{
		app:            app,
		ctx:            context.Background(),
		organizationID: organizationID,
		userID:         userID,
	}
}

func (l Locked) IsLocked() bool {
	_, err := l.app.DB.GetLockedProfileIdentityForUserID(l.ctx, l.organizationID, l.userID)
	return err == nil
}
