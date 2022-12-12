package profiles

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"attractify.io/platform/app"
	"attractify.io/platform/db"
	"github.com/gofrs/uuid"
)

const maxIdentities = 500

type Profile struct {
	ctx        context.Context
	app        *app.App
	params     Params
	traits     db.Traits
	traitsJSON json.RawMessage
}

type Params struct {
	Time           time.Time
	OrganizationID uuid.UUID
	UserID         string
	PreviousUserID string
	Channel        string
	Type           string
	IsAnonymous    bool
	Traits         *json.RawMessage
}

func New(ctx context.Context, app *app.App, params Params) *Profile {
	return &Profile{
		ctx:        ctx,
		app:        app,
		params:     params,
		traits:     db.Traits{},
		traitsJSON: json.RawMessage("{}"),
	}
}

func (p *Profile) GetOrCreate() (*db.Profile, *db.ProfileIdentity, error) {
	// Get Profile identity by user ID.
	profileIdentities, err := p.getIdentities(p.params.UserID)
	if err != nil {
		return nil, nil, fmt.Errorf("could not get identities from DB (%s): %w", p.params.UserID, err)
	}

	// Profile identity does not exist.
	var profileIdentity db.ProfileIdentity
	if len(profileIdentities) == 0 {
		profile, err := p.createProfile()
		if err != nil {
			return nil, nil, fmt.Errorf("could not create profile (%v): %w", p, err)
		}

		profileIdentity, err = p.createProfileIdentity(profile.ID)
		if err != nil {
			if err := p.deleteProfile(profile.ID); err != nil {
				return nil, nil, fmt.Errorf("could not delete profile (%s): %w", profile.ID.String(), err)
			}

			// TODO: Find better way than recursion.
			return p.GetOrCreate()
		}

		return &profile, &profileIdentity, nil
	}

	// Profile exists, check if identity exists for channel and type.
	identityFound := false
	for _, pi := range profileIdentities {
		if pi.Channel == p.params.Channel {
			profileIdentity = pi
			identityFound = true
		}
	}

	// Identity does not exist in channel, create new identity.
	if !identityFound {
		// Ensure that maxIdentities limit is not exceeded.
		count, err := p.countProfileIdentities(profileIdentities[0].ProfileID)
		if err != nil {
			return nil, nil, fmt.Errorf("could not count identities (%s): %w", profileIdentities[0].ProfileID.String(), err)
		}
		if count >= maxIdentities {
			return nil, nil, fmt.Errorf("too many identities for profile (%s)", profileIdentities[0].ProfileID.String())
		}

		profileIdentity, err = p.createProfileIdentity(profileIdentities[0].ProfileID)
		if err != nil {
			return nil, nil, fmt.Errorf("could not create identity (%s): %w", profileIdentities[0].ProfileID.String(), err)
		}
	}

	// Return profile.
	profile, err := p.getProfile(profileIdentity.ProfileID)
	if err != nil {
		return nil, nil, fmt.Errorf("could not get profile from DB (%s): %w", profileIdentity.ProfileID.String(), err)
	}

	return &profile, &profileIdentity, nil
}

func (p *Profile) UpdateOrCreate() error {
	if p.params.Traits != nil {
		if err := p.validateTraits(); err != nil {
			return fmt.Errorf("could not validate traits (%v): %w", p.params.Traits, err)
		}
		if err := p.prepareTraits(); err != nil {
			return fmt.Errorf("could not prepare traits (%v): %w", p.params.Traits, err)
		}
	}

	profile, identity, err := p.GetOrCreate()
	if err != nil {
		return err
	}

	// Update custom traits.
	if p.params.Traits != nil {
		var traits db.Traits
		if err := json.Unmarshal(profile.CustomTraits, &traits); err != nil {
			return err
		}
		if !reflect.DeepEqual(traits, p.traits) {
			p.traitsJSON = *p.params.Traits
			profile.CustomTraits = *p.params.Traits
			profile.UpdatedAt = time.Now().UTC()

			if err := p.updateProfile(profile); err != nil {
				return fmt.Errorf("could not update profile (%v): %w", profile, err)
			}
		}
	}

	// If type or anonymity flag has changed, update identity.
	if identity.Type != p.params.Type || identity.IsAnonymous != p.params.IsAnonymous {
		identity.UpdatedAt = time.Now().UTC()
		identity.Type = p.params.Type
		identity.IsAnonymous = p.params.IsAnonymous
		if err := p.updateIdentityType(identity); err != nil {
			return fmt.Errorf("could not update identity type (%v): %w", identity.ID.String(), err)
		}
	}

	// If a previous profile ID exists, merge profiles.
	if len(p.params.PreviousUserID) > 0 {
		prevIdentity, err := p.getIdentity(p.params.PreviousUserID)
		if err != nil {
			return fmt.Errorf("could not get previous identity (%s): %w ", p.params.PreviousUserID, err)
		}
		prevProfile, err := p.getProfile(prevIdentity.ProfileID)
		if err != nil {
			return fmt.Errorf("could not get previous profile (%s): %w ", prevIdentity.ProfileID.String(), err)
		}

		if prevProfile.ID == profile.ID {
			return nil
		}

		return p.mergeProfiles(profile, &prevProfile, identity)
	}

	return nil
}

func (p *Profile) mergeProfiles(leader, follower *db.Profile, identity *db.ProfileIdentity) error {
	var err error
	leader.CustomTraits, err = p.mergeTraits(leader.CustomTraits, follower.CustomTraits)
	if err != nil {
		return err
	}

	leader.ComputedTraits, err = p.mergeTraits(leader.ComputedTraits, follower.ComputedTraits)
	if err != nil {
		return err
	}

	if err := p.updateProfile(leader); err != nil {
		return fmt.Errorf("could not update profile (%s) while merging: %w", leader.ID.String(), err)
	}

	// Ensure that maxIdentities limit is not exceeded.
	countLeader, err := p.countProfileIdentities(leader.ID)
	if err != nil {
		return fmt.Errorf("could not count identities (%s): %w", leader.ID.String(), err)
	}
	countFollower, err := p.countProfileIdentities(follower.ID)
	if err != nil {
		return fmt.Errorf("could not count identities (%s): %w", follower.ID.String(), err)
	}
	if countLeader+countFollower >= maxIdentities {
		return fmt.Errorf("too many identities for profiles (%s - %s)", follower.ID.String(), leader.ID.String())
	}

	if err := p.relinkIdentitiesInDB(leader.ID, follower.ID); err != nil {
		return fmt.Errorf(
			"could not relink identities (%s - %s) while merging: %w",
			leader.ID.String(), follower.ID.String(), err,
		)
	}

	if err := p.deleteProfile(follower.ID); err != nil {
		return fmt.Errorf("could not delete profile (%s) while merging: %w", follower.ID.String(), err)
	}

	return nil
}

func (p Profile) getIdentity(userID string) (db.ProfileIdentity, error) {
	return p.app.DB.GetProfileIdentityForUserID(p.ctx, p.params.OrganizationID, userID)
}

func (p Profile) getIdentities(userID string) ([]db.ProfileIdentity, error) {
	return p.app.DB.GetProfileIdentitiesForUserID(p.ctx, p.params.OrganizationID, userID)
}

func (p *Profile) createProfile() (db.Profile, error) {
	args := db.CreateProfileParams{
		OrganizationID: p.params.OrganizationID,
		CustomTraits:   p.traitsJSON,
		CreatedAt:      p.params.Time,
	}
	return p.app.DB.CreateProfile(p.ctx, args)
}

func (p *Profile) createProfileIdentity(profileID uuid.UUID) (db.ProfileIdentity, error) {
	args := db.CreateProfileIdentityParams{
		OrganizationID: p.params.OrganizationID,
		ProfileID:      profileID,
		Channel:        p.params.Channel,
		Type:           p.params.Type,
		UserID:         p.params.UserID,
		IsAnonymous:    p.params.IsAnonymous,
		CreatedAt:      p.params.Time,
	}
	return p.app.DB.CreateProfileIdentity(p.ctx, args)
}

func (p *Profile) countProfileIdentities(profileID uuid.UUID) (int, error) {
	return p.app.DB.CountProfileIdentities(p.ctx, p.params.OrganizationID, profileID)
}

func (p *Profile) deleteProfile(id uuid.UUID) error {
	return p.app.DB.DeleteProfile(p.ctx, p.params.OrganizationID, id)
}

func (p *Profile) getProfile(id uuid.UUID) (db.Profile, error) {
	return p.app.DB.GetProfile(p.ctx, p.params.OrganizationID, id)
}

func (p Profile) updateProfile(profile *db.Profile) error {
	args := db.UpdateProfileParams{
		CustomTraits:   profile.CustomTraits,
		ComputedTraits: profile.ComputedTraits,
		OrganizationID: profile.OrganizationID,
		ID:             profile.ID,
		UpdatedAt:      profile.UpdatedAt,
	}
	return p.app.DB.UpdateProfile(p.ctx, args)
}

func (p *Profile) relinkIdentitiesInDB(leaderID, followerID uuid.UUID) error {
	return p.app.DB.UpdateProfileIdentitiesWithProfileID(
		p.ctx, p.params.OrganizationID, followerID, leaderID,
	)
}

func (p *Profile) updateIdentityType(identity *db.ProfileIdentity) error {
	args := db.UpdateProfileIdentityParams{
		OrganizationID: p.params.OrganizationID,
		ID:             identity.ID,
		UpdatedAt:      identity.UpdatedAt,
		Type:           identity.Type,
		IsAnonymous:    identity.IsAnonymous,
	}
	return p.app.DB.UpdateProfileIdentity(p.ctx, args)
}
