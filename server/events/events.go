package events

import (
	"context"
	"encoding/json"
	"time"

	"attractify.io/platform/analytics"
	"attractify.io/platform/app"
	"attractify.io/platform/db"
	"attractify.io/platform/profiles"
	"github.com/gofrs/uuid"
)

type Params struct {
	Time           time.Time
	OrganizationID uuid.UUID
	UserID         string
	Channel        string
	Event          string
	Context        *json.RawMessage
	Properties     *json.RawMessage
}

type Event struct {
	ctx        context.Context
	app        *app.App
	params     Params
	profile    *db.Profile
	identity   *db.ProfileIdentity
	event      *db.Event
	context    analytics.Properties
	properties analytics.Properties
}

func New(ctx context.Context, app *app.App, params Params) *Event {
	return &Event{
		ctx:        ctx,
		app:        app,
		params:     params,
		context:    analytics.Properties{},
		properties: analytics.Properties{},
	}
}

func (e Event) Profile() *db.Profile {
	return e.profile
}

func (e *Event) Track() error {
	if err := e.getEvent(); err != nil {
		return err
	}

	// Validate and prepare context properties.
	if e.params.Context != nil {
		if err := e.validateContext(); err != nil {
			invalidParams := db.CreateInvalidEventParams{
				EventID:        e.event.ID,
				OrganizationID: e.params.OrganizationID,
				Channel:        e.params.Channel,
				Properties:     *e.params.Properties,
				Context:        *e.params.Context,
				Type:           "context",
				CreatedAt:      e.params.Time,
			}

			e.app.DB.CreateInvalidEvent(e.ctx, invalidParams)
			return err
		}
	}

	// Validate and prepare event properties.
	if e.params.Properties != nil {
		if err := e.validateProperties(); err != nil {
			invalidParams := db.CreateInvalidEventParams{
				EventID:        e.event.ID,
				OrganizationID: e.params.OrganizationID,
				Channel:        e.params.Channel,
				Properties:     *e.params.Properties,
				Context:        *e.params.Context,
				Type:           "properties",
				CreatedAt:      e.params.Time,
			}

			e.app.DB.CreateInvalidEvent(e.ctx, invalidParams)
			return err
		}
	}

	// Resolve profile
	params := profiles.Params{
		Time:           e.params.Time,
		OrganizationID: e.params.OrganizationID,
		UserID:         e.params.UserID,
		Channel:        e.params.Channel,
		Type:           "anonymous_id",
		IsAnonymous:    true,
	}
	p := profiles.New(e.ctx, e.app, params)

	var err error
	e.profile, e.identity, err = p.GetOrCreate()
	if err != nil {
		return err
	}

	// Create event.
	if err := e.createEvent(); err != nil {
		return err
	}

	return nil
}

func (e *Event) createEvent() error {
	context := "{}"
	if e.params.Context != nil {
		context = string(*e.params.Context)
	}

	properties := "{}"
	if e.params.Properties != nil {
		properties = string(*e.params.Properties)
	}

	args := analytics.CreateEventParams{
		ID:             uuid.Must(uuid.NewV4()),
		OrganizationID: e.params.OrganizationID,
		IdentityID:     e.identity.ID,
		EventID:        e.event.ID,
		Channel:        e.params.Channel,
		Context:        context,
		Properties:     properties,
		CreatedAt:      e.params.Time.UTC(),
	}
	return e.app.Analytics.CreateEvent(args)
}

func (e *Event) EventID() uuid.UUID {
	return e.event.ID
}

func (e *Event) getEvent() error {
	evt, err := e.app.DB.GetEventByName(e.ctx, e.params.OrganizationID, e.params.Event)
	e.event = &evt
	return err
}
