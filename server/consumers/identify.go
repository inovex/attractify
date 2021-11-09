package consumers

import (
	"time"

	"attractify.io/platform/privacy"
	"attractify.io/platform/profiles"
	"attractify.io/platform/stream"
	"go.uber.org/zap"
)

func (c Consumer) Identify(m *stream.IdentifyMsg, t time.Time) error {
	l := privacy.NewLocked(c.app, m.OrganizationID, m.UserID)
	if l.IsLocked() {
		return errUserIDisLocked
	}

	params := profiles.Params{
		Time:           t,
		OrganizationID: m.OrganizationID,
		UserID:         m.UserID,
		PreviousUserID: m.PreviousUserID,
		Channel:        m.Channel,
		Type:           m.Type,
		IsAnonymous:    m.IsAnonymous,
		Traits:         m.Traits,
	}
	p := profiles.New(c.ctx, c.app, params)
	if err := p.UpdateOrCreate(); err != nil {
		c.app.Logger.Warn("api.identify.identify.profileHandler", zap.Error(err))
		return nil
	}

	return nil
}
