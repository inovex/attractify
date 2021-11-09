package stream

import (
	"context"
	"encoding/json"

	"github.com/gofrs/uuid"
)

type TrackMsg struct {
	OrganizationID uuid.UUID        `json:"organizationID"`
	UserID         string           `json:"userID"`
	Channel        string           `json:"channel"`
	Event          string           `json:"event"`
	Properties     *json.RawMessage `json:"properties"`
	Context        *json.RawMessage `json:"context"`
}

func (t TrackMsg) Marshal() []byte {
	b, _ := json.Marshal(t)
	return b
}

func (s Stream) Track(ctx context.Context, msg TrackMsg) error {
	k := []byte(msg.UserID)
	m := Msg{Type: MsgTypeTrack, RawMsg: msg.Marshal()}
	return s.Write(ctx, k, m.Marshal())
}
