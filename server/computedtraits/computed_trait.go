package computedtraits

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"attractify.io/platform/app"
	"attractify.io/platform/db"
)

type ComputedTrait struct {
	ctx        context.Context
	app        *app.App
	ct         *db.ComputedTrait
	conditions []condition
	properties *properties
	profile    *db.Profile
}

func New(ctx context.Context, app *app.App, ct *db.ComputedTrait) *ComputedTrait {
	return &ComputedTrait{
		ctx: ctx,
		app: app,
		ct:  ct,
	}
}

func (c *ComputedTrait) parseConditions() error {
	return json.Unmarshal(c.ct.Conditions, &c.conditions)
}

func (c *ComputedTrait) parseProperties() error {
	return json.Unmarshal(c.ct.Properties, &c.properties)
}

func (c *ComputedTrait) profileIdentities() string {
	return fmt.Sprintf("AND identity_id GLOBAL IN (SELECT id FROM identities WHERE organization_id = '%s' AND profile_id = '%s')", c.profile.OrganizationID, c.profile.ID)
}

func (c *ComputedTrait) mergeTraits(newK string, newV interface{}) error {
	var traits db.Traits
	if err := json.Unmarshal(c.profile.ComputedTraits, &traits); err != nil {
		return err
	}

	traits[newK] = newV

	var err error
	c.profile.ComputedTraits, err = json.Marshal(traits)
	return err
}

func (c *ComputedTrait) Refresh(profile *db.Profile) error {
	if err := c.parseConditions(); err != nil {
		return err
	}

	if err := c.parseProperties(); err != nil {
		return err
	}

	c.profile = profile

	var (
		res interface{}
		err error
	)
	switch c.ct.Type {
	case "count_events":
		res, err = c.eventCount()
	case "aggregation":
		res, err = c.aggregation()
	case "most_frequent":
		res, err = c.mostFrequent()
	case "first_event":
		res, err = c.edgeEvent(true)
	case "last_event":
		res, err = c.edgeEvent(false)
	case "unique_list":
		res, err = c.uniqueList()
	case "unique_list_count":
		res, err = c.uniqueListCount()
	}
	if err != nil {
		return err
	}

	if err := c.mergeTraits(c.ct.Key, res); err != nil {
		return err
	}

	now := time.Now().UTC()

	dbArgs := db.UpdateProfileComputedTraitsParams{
		ComputedTraits: c.profile.ComputedTraits,
		UpdatedAt:      now,
		ID:             c.profile.ID,
	}

	return c.app.DB.UpdateProfileComputedTraits(c.ctx, dbArgs)
}
