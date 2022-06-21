package platform

import (
	"net/http"

	"attractify.io/platform/app"
	"attractify.io/platform/db"
	"attractify.io/platform/platform/responses"
	"github.com/gin-gonic/gin"
)

type DashboardController struct {
	Router *gin.RouterGroup
	App    *app.App
}

func InitDashboard(router *gin.RouterGroup, app *app.App) {
	c := DashboardController{Router: router, App: app}
	c.Router.GET("/dashboard", c.show)
}

func (dc DashboardController) show(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	rRes, err := dc.App.Analytics.GetReactionsForInterval(user.OrganizationID)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	reactions := []responses.DashboardBucket{}
	for _, p := range rRes {
		reactions = append(reactions, responses.DashboardBucket{Bucket: p.Bucket, Count: p.Count})
	}

	pRes, err := dc.App.DB.GetNewProfilesLast24h(c.Request.Context(), user.OrganizationID)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	profiles := []responses.DashboardBucket{}
	for _, p := range pRes {
		profiles = append(profiles, responses.DashboardBucket{Bucket: p.Bucket, Count: p.Count})
	}

	events, err := dc.App.Analytics.GetLastDaysEventCount(user.OrganizationID)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	actions, err := dc.App.DB.GetActiveActionsCount(c.Request.Context(), user.OrganizationID)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	res := responses.Dashboard{
		Reactions: reactions,
		Profiles:  profiles,
		Events:    events,
		Actions:   actions,
	}
	c.JSON(http.StatusOK, res)
}
