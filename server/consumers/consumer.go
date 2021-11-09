package consumers

import (
	"context"
	"log"

	"attractify.io/platform/app"
	"attractify.io/platform/stream"
)

type Consumer struct {
	ctx context.Context
	app *app.App
}

func New(ctx context.Context, app *app.App) *Consumer {
	return &Consumer{ctx: ctx, app: app}
}

func (c Consumer) ProcessMsg(m *stream.Msg) error {
	var err error
	switch m.Type {
	case stream.MsgTypeTrack:
		err = c.Track(&m.TrackMsg, m.Time)
	case stream.MsgTypeIdentify:
		err = c.Identify(&m.IdentifyMsg, m.Time)
	}

	if err != nil {
		return err
	}

	log.Println("Successfully processed", m.Type, m.Time)
	return nil
}
