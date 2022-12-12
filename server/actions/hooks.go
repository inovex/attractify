package actions

import (
	"encoding/json"

	"attractify.io/platform/app"
	"attractify.io/platform/db"
	"github.com/gofrs/uuid"
)

type Hook struct {
	App               *app.App
	Config            json.RawMessage
	OrganizationID    uuid.UUID
	Action            *db.Action
	ProfileIdentityID uuid.UUID
	UserID            string
	Event             string
	Channel           string
	Context           *json.RawMessage
	Properties        *json.RawMessage
}

func (a *Action) parseHooks() error {
	return json.Unmarshal(a.action.Hooks, &a.hooks)
}

func (a Action) RunHooks(userID, event, channel string, context, properties *json.RawMessage) (json.RawMessage, bool, error) {
	if err := a.parseHooks(); err != nil {
		return nil, nil, err
	}

	isHookSuccessful := true
	var result json.RawMessage
	for _, hook := range a.hooks {
		if len(hook.Channels) > 0 && !a.inChannel(hook.Channels, channel) {
			continue
		}

		if hook.Event != event {
			continue
		}

		h := Hook{
			App:               a.app,
			Config:            hook.Properties,
			OrganizationID:    a.organizationID,
			Action:            a.action,
			ProfileIdentityID: a.profileIdentity.ID,
			UserID:            userID,
			Event:             event,
			Channel:           channel,
			Context:           context,
			Properties:        properties,
		}

		switch hook.Type {
		case "execute_webhook":
			res, err := h.ExecuteWebhook()
			if err != nil {
				return nil, nil, err
			}
			if res != nil {
				if res.StatusCode >= 300 {
					isHookSuccessful = false
				}
				result, _ = json.Marshal(res)
			}
		case "track_event":
			if err := h.TrackEvent(); err != nil {
				return nil, nil, err
			}
		}
	}

	return result, isHookSuccessful, nil
}

func (a Action) inChannel(channels []string, channel string) bool {
	for _, c := range channels {
		if c == channel {
			return true
		}
	}
	return false
}
