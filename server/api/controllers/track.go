package api

import (
	"net/http"
	"time"

	"attractify.io/platform/api/requests"
	"attractify.io/platform/app"
	"attractify.io/platform/computedtraits"
	"attractify.io/platform/db"
	"attractify.io/platform/events"
	"attractify.io/platform/privacy"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type TrackController struct {
	Router *gin.RouterGroup
	App    *app.App
}

func InitTrack(router *gin.RouterGroup, app *app.App) {
	t := TrackController{Router: router, App: app}
	t.Router.POST("/track", t.Track)
}

func (tc TrackController) Track(c *gin.Context) {
	var req requests.Track
	if err := c.ShouldBindJSON(&req); err != nil {
		tc.App.Logger.Warn("api.track.track.parseRequest", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	auth := c.MustGet("auth").(*db.OrganizationAuth)

	l := privacy.NewLocked(tc.App, auth.OrganizationID, req.UserID)
	if l.IsLocked() {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}

	params := events.Params{
		Time:           time.Now().UTC(),
		OrganizationID: auth.OrganizationID,
		UserID:         req.UserID,
		Channel:        auth.Channel,
		Event:          req.Event,
		Properties:     req.Properties,
		Context:        req.Context,
	}
	e := events.New(c.Request.Context(), tc.App, params)
	if err := e.Track(); err != nil {
		tc.App.Logger.Warn("api.track.track.processEvent", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
	}

	computedTraits, err := tc.App.DB.GetComputedTraitsForEvent(
		c.Request.Context(), auth.OrganizationID, e.EventID())
	if err != nil {
		tc.App.Logger.Warn("api.track.track.getComputedTraitsForEvent", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
	}

	for _, ct := range computedTraits {
		ctr := computedtraits.New(c.Request.Context(), tc.App, &ct)
		if err := ctr.Refresh(e.Profile()); err != nil {
			tc.App.Logger.Warn("api.track.track.processComputedTraits", zap.Error(err))
			c.AbortWithStatus(http.StatusBadRequest)
		}
	}

	c.AbortWithStatus(http.StatusNoContent)
}
