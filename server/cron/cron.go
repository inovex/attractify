package cron

import (
	"context"

	"attractify.io/platform/app"
	"attractify.io/platform/audiences"
	"github.com/robfig/cron/v3"
)

type Cron struct {
	ctx  context.Context
	app  *app.App
	cron *cron.Cron
}

func New(ctx context.Context, app *app.App) *Cron {
	c := cron.New(cron.WithChain(
		cron.SkipIfStillRunning(cron.DiscardLogger),
	))
	return &Cron{ctx: ctx, app: app, cron: c}
}

func (c Cron) Init() {
	c.cron.AddFunc("@every 12h", func() { c.updateAudiences() })
	c.cron.Start()
}

func (c Cron) updateAudiences() error {
	a, err := c.app.DB.GetAllAudiences(c.ctx)
	if err != nil {
		return err
	}

	for _, audience := range a {
		a := audiences.New(c.ctx, c.app, &audience)
		if _, err := a.Refresh(); err != nil {
			return err
		}
	}

	return nil
}
