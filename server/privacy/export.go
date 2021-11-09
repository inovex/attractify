package privacy

import (
	"bytes"
	"context"
	"encoding/json"

	"attractify.io/platform/analytics"
	"attractify.io/platform/app"
	"attractify.io/platform/db"
	"github.com/gofrs/uuid"
)

type Export struct {
	app            *app.App
	ctx            context.Context
	organizationID uuid.UUID
	userID         string
	email          string
}

type DataExport struct {
	Profile    db.Profile           `json:"profile"`
	Identities []db.ProfileIdentity `json:"identities"`
	Events     []analytics.Event    `json:"events"`
	Reactions  []analytics.Reaction `json:"reactions"`
}

func (d DataExport) Marshal() ([]byte, error) {
	return json.MarshalIndent(d, "", "  ")
}

func NewExport(app *app.App, organizationID uuid.UUID, userID, email string) *Export {
	return &Export{
		app:            app,
		ctx:            context.Background(),
		organizationID: organizationID,
		userID:         userID,
		email:          email,
	}
}

func (e Export) Run() error {
	identity, err := e.getIdentity()
	if err != nil {
		return err
	}

	identities, err := e.getIdentities(identity.ProfileID)
	if err != nil {
		return err
	}

	profile, err := e.getProfile(identity.ProfileID)
	if err != nil {
		return err
	}

	identityIDs := []uuid.UUID{}
	for _, i := range identities {
		identityIDs = append(identityIDs, i.ID)
	}

	events, err := e.getEvents(identityIDs)
	if err != nil {
		return err
	}

	reactions, err := e.getReactions(identityIDs)
	if err != nil {
		return err
	}

	data := DataExport{
		Profile:    profile,
		Identities: identities,
		Events:     events,
		Reactions:  reactions,
	}

	b, err := data.Marshal()
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(b)

	return e.app.Mailer.SendDataExport(e.email, buf)
}

func (e Export) getIdentity() (db.ProfileIdentity, error) {
	return e.app.DB.GetProfileIdentityForUserID(e.ctx, e.organizationID, e.userID)
}

func (e Export) getIdentities(profileID uuid.UUID) ([]db.ProfileIdentity, error) {
	return e.app.DB.GetProfileIdentitiesForProfile(e.ctx, e.organizationID, profileID)
}

func (e Export) getProfile(id uuid.UUID) (db.Profile, error) {
	return e.app.DB.GetProfile(e.ctx, e.organizationID, id)
}

func (e Export) getEvents(identities []uuid.UUID) ([]analytics.Event, error) {
	p := analytics.GetEventsForIdentitiesParams{
		OrganizationID: e.organizationID,
		IdentityIDs:    identities,
	}
	return e.app.Analytics.GetEventsForIdentities(p)
}

func (e Export) getReactions(identities []uuid.UUID) ([]analytics.Reaction, error) {
	p := analytics.GetReactionsForIdentitiesParams{
		OrganizationID: e.organizationID,
		IdentityIDs:    identities,
	}
	return e.app.Analytics.GetReactionsForIdentities(p)
}
