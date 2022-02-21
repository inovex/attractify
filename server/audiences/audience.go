package audiences

import (
	"context"
	"encoding/json"

	"attractify.io/platform/analytics"
	"attractify.io/platform/app"
	"attractify.io/platform/db"
	"github.com/gofrs/uuid"
)

const (
	batchSize   = 1000
	previewSize = 1000
)

type Audience struct {
	ctx        context.Context
	app        *app.App
	audience   *db.Audience
	conditions conditions
}

func New(ctx context.Context, app *app.App, audience *db.Audience) *Audience {
	return &Audience{
		ctx:        ctx,
		app:        app,
		audience:   audience,
		conditions: conditions{},
	}
}

func (a Audience) Preview() ([]analytics.AudienceProfile, error) {
	if err := a.parseConditions(); err != nil {
		return nil, err
	}
	return a.runQuery()
}

func (a Audience) Refresh() (int, error) {
	if err := a.parseConditions(); err != nil {
		return 0, err
	}

	// Build and run query.
	profiles, err := a.runQuery()
	if err != nil {
		return 0, err
	}

	// Write retrieved profiles to mapping table.
	count := 0
	setID, _ := uuid.NewV4()
	var argList []db.CreateAudienceProfilesParams
	for i, p := range profiles {
		argList = append(argList, db.CreateAudienceProfilesParams{
			OrganizationID: a.audience.OrganizationID,
			AudienceID:     a.audience.ID,
			ProfileID:      p.ID,
			SetID:          setID,
		})

		// Persist to DB.
		if i%batchSize == 0 {
			if err := a.app.DB.CreateAudienceProfiles(a.ctx, argList); err != nil {
				return 0, err
			}
			argList = []db.CreateAudienceProfilesParams{}
		}
		count++
	}
	if err := a.app.DB.CreateAudienceProfiles(a.ctx, argList); err != nil {
		return 0, err
	}

	// Set audience to sucessfully processed.
	params := db.UpdateAudienceProfilesParams{
		CurrentSetID:   setID,
		ProfileCount:   count,
		OrganizationID: a.audience.OrganizationID,
		ID:             a.audience.ID,
	}
	if err := a.app.DB.UpdateAudienceProfiles(a.ctx, params); err != nil {
		return 0, err
	}

	// Remove old set of audience profiles.
	if a.audience.CurrentSetID.Valid {
		if err := a.app.DB.DeleteAudienceProfilesBySetID(
			a.ctx,
			a.audience.OrganizationID,
			a.audience.CurrentSetID.UUID,
		); err != nil {
			return 0, err
		}
	}

	return count, nil
}

func (a *Audience) parseConditions() error {
	if err := json.Unmarshal(a.audience.Events, &a.conditions.events); err != nil {
		return err
	}
	if err := json.Unmarshal(a.audience.Traits, &a.conditions.traits); err != nil {
		return err
	}
	a.conditions.includeAnonymous = a.audience.IncludeAnonymous

	return nil
}

func (a *Audience) runQuery() ([]analytics.AudienceProfile, error) {
	query := a.conditions.BuildQuery(a.audience.OrganizationID)
	println(query)
	return a.app.Analytics.RunAudience(query)
}
