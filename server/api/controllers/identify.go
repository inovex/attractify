package api

import (
	"net/http"
	"time"

	"attractify.io/platform/api/requests"
	"attractify.io/platform/app"
	"attractify.io/platform/db"
	"attractify.io/platform/privacy"
	"attractify.io/platform/profiles"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type IdentifyController struct {
	Router *gin.RouterGroup
	App    *app.App
}

func InitIdentify(router *gin.RouterGroup, app *app.App) {
	i := IdentifyController{Router: router, App: app}
	i.Router.POST("/identify", i.Identify)
}

func (ic IdentifyController) Identify(c *gin.Context) {
	var req requests.Identify
	if err := c.ShouldBindJSON(&req); err != nil {
		ic.App.Logger.Warn("api.identify.identify.parseRequest", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	auth := c.MustGet("auth").(*db.OrganizationAuth)

	// Is user locked
	l := privacy.NewLocked(ic.App, auth.OrganizationID, req.UserID)
	if l.IsLocked() {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ut := req.Type
	if req.IsAnonymous {
		ut = "anonymous_id"
	}
	params := profiles.Params{
		Time:           time.Now().UTC(),
		OrganizationID: auth.OrganizationID,
		UserID:         req.UserID,
		PreviousUserID: req.PreviousUserID,
		Channel:        auth.Channel,
		Type:           ut,
		IsAnonymous:    req.IsAnonymous,
		Traits:         req.Traits,
	}
	p := profiles.New(c.Request.Context(), ic.App, params)
	if err := p.UpdateOrCreate(); err != nil {
		ic.App.Logger.Warn("api.identify.identify.profileHandler", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}
