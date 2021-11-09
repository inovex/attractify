package consumers

import (
	"errors"
	"time"

	"attractify.io/platform/computedtraits"
	"attractify.io/platform/db"
	"attractify.io/platform/events"
	"attractify.io/platform/privacy"
	"attractify.io/platform/stream"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
)

func (c Consumer) Track(m *stream.TrackMsg, t time.Time) error {
	l := privacy.NewLocked(c.app, m.OrganizationID, m.UserID)
	if l.IsLocked() {
		return errUserIDisLocked
	}

	params := events.Params{
		Time:           t,
		OrganizationID: m.OrganizationID,
		UserID:         m.UserID,
		Channel:        m.Channel,
		Event:          m.Event,
		Properties:     m.Properties,
		Context:        m.Context,
	}
	e := events.New(c.ctx, c.app, params)
	if err := e.Track(); err != nil {
		c.app.Logger.Warn("api.consumers.track.track", zap.Error(err))
		return nil
	}

	if err := c.processComputedTraits(m.OrganizationID, e.EventID(), e.Profile()); err != nil {
		c.app.Logger.Warn("api.consumers.track.processComputedTraits", zap.Error(err))
		return nil
	}

	return nil
}

func (c Consumer) processComputedTraits(organizationID, eventID uuid.UUID, profile *db.Profile) error {
	computedTraits, err := c.app.DB.GetComputedTraitsForEvent(c.ctx, organizationID, eventID)
	if err != nil {
		return err
	}

	for _, ct := range computedTraits {
		ctr := computedtraits.New(c.ctx, c.app, &ct)
		if err := ctr.Refresh(profile); err != nil {
			c.app.Logger.Warn("api.consumers.track.refresh", zap.Error(err))
		}
	}

	return nil
}

var (
	errUserIDisLocked = errors.New("user ID is locked")
)
