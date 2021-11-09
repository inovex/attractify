package platform

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"attractify.io/platform/analytics"
	"attractify.io/platform/app"
	"attractify.io/platform/db"
	"attractify.io/platform/platform/requests"
	"attractify.io/platform/platform/responses"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
)

type EventLogController struct {
	Router *gin.RouterGroup
	App    *app.App
}

func InitEventLog(router *gin.RouterGroup, app *app.App) {
	c := EventLogController{Router: router, App: app}
	c.Router.GET("/event-log", c.List)
	c.Router.DELETE("/event-log/:id", c.Delete)
}

func (ec EventLogController) List(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	var req requests.EventsList
	if err := c.ShouldBindQuery(&req); err != nil {
		ec.App.Logger.Warn("events.list.parseRequest", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var identityIDs []uuid.UUID
	if len(req.UserID) > 0 {
		identity, err := ec.App.DB.GetProfileIdentityForUserID(
			c.Request.Context(), user.OrganizationID, req.UserID,
		)
		if err == nil {
			identityIDs = append(identityIDs, identity.ID)
		} else {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
	}

	// If no start date is given, set start date to the lowest possible.
	if len(req.Start) == 0 {
		req.Start = "2000-01-01"
	}
	start, _ := time.Parse("2006-01-02", req.Start)

	// If no end date is given use a ridiculous high end date (FIX ME)
	// and add one day to make it the end of the day.
	if len(req.End) == 0 {
		req.End = "2099-01-01"
	}
	end, _ := time.Parse("2006-01-02", req.End)
	end = end.AddDate(0, 0, 1)

	var (
		events []analytics.Event
		err    error
	)

	eventList := []uuid.UUID{}
	if len(req.Events) > 0 {
		for _, e := range strings.Split(req.Events, ",") {
			eventList = append(eventList, uuid.FromStringOrNil(e))
		}
	}

	offset := req.ItemsPerPage * (req.Page - 1)
	args := analytics.GetEventsParams{
		OrganizationID: user.OrganizationID,
		EventIDs:       eventList,
		IdentityIDs:    identityIDs,
		Start:          start,
		End:            end,
		Limit:          req.ItemsPerPage,
		Offset:         offset,
	}
	events, err = ec.App.Analytics.GetEvents(args)

	if err != nil {
		ec.App.Logger.Error("events.list.getEvents", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	eventRes := []responses.EventLog{}
	for _, e := range events {
		eventRes = append(eventRes, responses.EventLog{
			ID:         e.ID,
			EventID:    e.EventID,
			Channel:    e.Channel,
			Context:    json.RawMessage(e.Context),
			Properties: json.RawMessage(e.Properties),
			CreatedAt:  e.CreatedAt,
		})
	}

	count := req.ItemsPerPage * (req.Page + 1)
	if len(eventRes) < req.ItemsPerPage {
		count = req.ItemsPerPage * req.Page
	}

	res := responses.EventLogList{
		Events: eventRes,
		// FIX ME: Find better way to calculate number of results.
		Count: count,
	}

	c.JSON(http.StatusOK, &res)
}

func (ec EventLogController) Delete(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	args := analytics.DeleteEventParams{
		OrganizationID: user.OrganizationID,
		ID:             uuid.FromStringOrNil(c.Param("id")),
	}
	if err := ec.App.Analytics.DeleteEvent(args); err != nil {
		ec.App.Logger.Error("events.delete.deleteAnalyticsEvent", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}
