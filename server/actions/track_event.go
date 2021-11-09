package actions

import (
	"encoding/json"
	"time"

	"attractify.io/platform/analytics"
	"github.com/gofrs/uuid"
)

type trackEvent struct {
	EventID uuid.UUID `json:"eventId"`
}

func (h Hook) TrackEvent() error {
	te := trackEvent{}
	if err := json.Unmarshal(h.Config, &te); err != nil {
		return err
	}

	context := "{}"
	if h.Context != nil {
		context = string(*h.Context)
	}

	properties := "{}"
	if h.Properties != nil {
		properties = string(*h.Properties)
	}

	args := analytics.CreateEventParams{
		ID:             uuid.Must(uuid.NewV4()),
		OrganizationID: h.OrganizationID,
		IdentityID:     h.ProfileIdentityID,
		EventID:        te.EventID,
		Channel:        h.Channel,
		Context:        context,
		Properties:     properties,
		CreatedAt:      time.Now().UTC(),
	}

	return h.App.Analytics.CreateEvent(args)
}
