package api

import (
	"net/http"

	"attractify.io/platform/api/requests"
	"attractify.io/platform/app"
	"attractify.io/platform/db"
	"attractify.io/platform/stream"
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

	msg := stream.TrackMsg{
		OrganizationID: auth.OrganizationID,
		UserID:         req.UserID,
		Channel:        auth.Channel,
		Event:          req.Event,
		Properties:     req.Properties,
		Context:        req.Context,
	}

	if err := tc.App.Stream.Track(c.Request.Context(), msg); err != nil {
		tc.App.Logger.Warn("api.track.track.sendToStream", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}
