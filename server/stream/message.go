package stream

import (
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	MsgTypeTrack    = "track"
	MsgTypeIdentify = "identify"
)

type MsgType string

type Msg struct {
	Type        MsgType         `json:"type"`
	Time        time.Time       `json:"-"`
	RawMsg      json.RawMessage `json:"payload"`
	TrackMsg    TrackMsg        `json:"-"`
	IdentifyMsg IdentifyMsg     `json:"-"`
}

func (m Msg) Marshal() []byte {
	b, _ := json.Marshal(m)
	return b
}

func NewMsg(m kafka.Message) (*Msg, error) {
	msg := Msg{Time: m.Time}
	if err := json.Unmarshal(m.Value, &msg); err != nil {
		return nil, err
	}
	switch msg.Type {
	case MsgTypeTrack:
		if err := json.Unmarshal(msg.RawMsg, &msg.TrackMsg); err != nil {
			return nil, err
		}
	case MsgTypeIdentify:
		if err := json.Unmarshal(msg.RawMsg, &msg.IdentifyMsg); err != nil {
			return nil, err
		}
	}
	return &msg, nil
}
