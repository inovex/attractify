package platform

import (
	"net/http"

	"attractify.io/platform/analytics"
	"attractify.io/platform/app"
	"attractify.io/platform/db"
	"attractify.io/platform/platform/requests"
	"attractify.io/platform/platform/responses"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
)

type AnalyzeController struct {
	Router *gin.RouterGroup
	App    *app.App
}

func InitAnalyze(router *gin.RouterGroup, app *app.App) {
	c := AnalyzeController{Router: router, App: app}
	c.Router.GET("/analyze/events", c.Events)
	c.Router.GET("/analyze/rates", c.Rates)
	c.Router.GET("/analyze/reach", c.Reach)
}

func (ac AnalyzeController) Events(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	var req requests.Analyze
	if err := c.ShouldBindQuery(&req); err != nil {
		ac.App.Logger.Warn("analyze.events.parseRequest", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	args := analytics.GetReactionEventCountForIntervalParams{
		OrganizationID: user.OrganizationID,
		ActionID:       uuid.FromStringOrNil(req.ActionID),
		Timezone:       "Europe/Berlin",
		Start:          req.Start,
		End:            req.End,
		Interval:       req.Interval,
	}
	events, err := ac.App.Analytics.GetEventOverview(args)
	if err != nil {
		ac.App.Logger.Error("analyze.events.getEventOverview", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var res []responses.AnalyzeEvents
	for _, e := range events {
		res = append(res, responses.AnalyzeEvents{
			Event:   e.Event,
			Total:   e.Total,
			Year:    e.Year,
			Month:   e.Month,
			Day:     e.Day,
			WeekDay: e.WeekDay,
			Hour:    e.Hour,
		})
	}

	c.JSON(http.StatusOK, res)
}

func (ac AnalyzeController) Rates(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	var req requests.Analyze
	if err := c.ShouldBindQuery(&req); err != nil {
		ac.App.Logger.Warn("analyze.rates.parseRequest", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	args := analytics.GetReactionEventCountParams{
		OrganizationID: user.OrganizationID,
		ActionID:       uuid.FromStringOrNil(req.ActionID),
		Start:          req.Start,
		End:            req.End,
	}
	agg, err := ac.App.Analytics.GetReactionEventCount(args)
	if err != nil {
		ac.App.Logger.Error("analyze.rates.GetActionEventCount", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	res := responses.AnalyzeRates{}
	for _, a := range agg {
		switch a.Event {
		case analytics.ReactionEventDelivered:
			res.Delivered = a.Total
		case analytics.ReactionEventShown:
			res.Shown = a.Total
		case analytics.ReactionEventHidden:
			res.Hidden = a.Total
		case analytics.ReactionEventDeclined:
			res.Declined = a.Total
		case analytics.ReactionEventAccepted:
			res.Accepted = a.Total
		}
	}

	c.JSON(http.StatusOK, res)
}

func (ac AnalyzeController) Reach(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	var req requests.Analyze
	if err := c.ShouldBindQuery(&req); err != nil {
		ac.App.Logger.Warn("analyze.reach.parseRequest", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	args := analytics.GetReactionChannelDeliveriesParams{
		OrganizationID: user.OrganizationID,
		ActionID:       uuid.FromStringOrNil(req.ActionID),
		Start:          req.Start,
		End:            req.End,
	}
	agg, err := ac.App.Analytics.GetReactionChannelDeliveries(args)
	if err != nil {
		ac.App.Logger.Error("analyze.reach.GetActionEventDeliveries", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	res := []responses.AnalyzeReach{}
	for _, a := range agg {
		res = append(res, responses.AnalyzeReach{
			Channel: a.Channel,
			Total:   a.Total,
		})
	}

	c.JSON(http.StatusOK, res)
}
