package privacy

import (
	"context"
	"errors"

	"attractify.io/platform/analytics"
	"attractify.io/platform/app"
	"attractify.io/platform/db"
	"github.com/gofrs/uuid"
)

type Deletion struct {
	app            *app.App
	ctx            context.Context
	organizationID uuid.UUID
	userID         string
	profileID      uuid.UUID
}

func NewDeletionByUserID(app *app.App, organizationID uuid.UUID, userID string) *Deletion {
	return &Deletion{
		app:            app,
		ctx:            context.Background(),
		organizationID: organizationID,
		userID:         userID,
	}
}

func NewDeletionByProfileID(app *app.App, organizationID uuid.UUID, profileID uuid.UUID) *Deletion {
	return &Deletion{
		app:            app,
		ctx:            context.Background(),
		organizationID: organizationID,
		profileID:      profileID,
	}
}

func (d Deletion) Run() error {
	if len(d.userID) > 0 {
		identity, err := d.getIdentityByUserID()
		if err != nil {
			return err
		}
		d.profileID = identity.ProfileID
	}

	identities, err := d.getIdentities(d.profileID)
	if err != nil {
		return err
	}

	identityIDs := []uuid.UUID{}
	for _, i := range identities {
		identityIDs = append(identityIDs, i.ID)
	}

	if len(identityIDs) > 0 {
		if err := d.deleteEvents(identityIDs); err != nil {
			return err
		}

		if err := d.deleteReactions(identityIDs); err != nil {
			return err
		}
	}

	if err := d.deleteProfile(d.profileID); err != nil {
		return err
	}

	return nil
}

func (d Deletion) DeleteSingleIdentity() error {
	if len(d.userID) == 0 {
		return errors.New("Deletion does not contain userID")
	}
	identityID, err := d.getIdentityByUserID()
	if err != nil {
		return err
	}
	params := analytics.DeleteEventsByIdentityIDParams{
		OrganizationID: d.organizationID,
		IdentityID:     identityID.ID,
	}
	err = d.app.Analytics.DeleteEventsByIdentityID(params)
	if err != nil {
		return err
	}
	return d.app.DB.DeleteProfileIdentityByID(d.ctx, d.organizationID, identityID.ID)
}

func (d Deletion) getIdentityByUserID() (db.ProfileIdentity, error) {
	return d.app.DB.GetProfileIdentityForUserID(d.ctx, d.organizationID, d.userID)
}

func (d Deletion) getIdentities(profileID uuid.UUID) ([]db.ProfileIdentity, error) {
	return d.app.DB.GetProfileIdentitiesForProfile(d.ctx, d.organizationID, profileID)
}

func (d Deletion) deleteEvents(identityIDs []uuid.UUID) error {
	p := analytics.DeleteEventsByIdentityIDsParams{
		OrganizationID: d.organizationID,
		IdentityIDs:    identityIDs,
	}
	return d.app.Analytics.DeleteEventsByIdentityIDs(p)
}

func (d Deletion) deleteReactions(identityIDs []uuid.UUID) error {
	p := analytics.DeleteReactionByIdentityIDsParams{
		OrganizationID: d.organizationID,
		IdentityIDs:    identityIDs,
	}
	return d.app.Analytics.DeleteReactionByIdentityIDs(p)
}

func (d Deletion) deleteProfile(profileID uuid.UUID) error {
	return d.app.DB.DeleteProfile(d.ctx, d.organizationID, profileID)
}
