package stream

import (
	"context"
	"encoding/json"

	"github.com/gofrs/uuid"
)

type IdentifyMsg struct {
	OrganizationID uuid.UUID        `json:"organizationID"`
	UserID         string           `json:"userID"`
	PreviousUserID string           `json:"previousUserID"`
	Channel        string           `json:"channel"`
	Type           string           `json:"event"`
	IsAnonymous    bool             `json:"isAnonymous"`
	Traits         *json.RawMessage `json:"traits"`
}

func (t IdentifyMsg) Marshal() []byte {
	b, _ := json.Marshal(t)
	return b
}

func (s Stream) Identify(ctx context.Context, msg IdentifyMsg) error {
	k := []byte(msg.UserID)
	m := Msg{Type: MsgTypeIdentify, RawMsg: msg.Marshal()}
	return s.Write(ctx, k, m.Marshal())
}
