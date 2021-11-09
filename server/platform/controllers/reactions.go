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

type ReactionsController struct {
	Router *gin.RouterGroup
	App    *app.App
}

func InitReactions(router *gin.RouterGroup, app *app.App) {
	c := ReactionsController{Router: router, App: app}
	c.Router.GET("/reactions", c.List)
	c.Router.DELETE("/reactions/:id", c.Delete)
}

func (rc ReactionsController) List(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	var req requests.Reactions
	if err := c.ShouldBindQuery(&req); err != nil {
		rc.App.Logger.Warn("actions.list.parseRequest", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// If no start date is given, set start date to the lowest possible.
	if len(req.Start) == 0 {
		req.Start = "2000-01-01"
	}
	start, _ := time.Parse("2006-01-02", req.Start)

	// If no end date is given use a ridiculously high end date (FIX ME)
	// and add one day to make it the end of the day.
	if len(req.End) == 0 {
		req.End = "2099-01-01"
	}
	end, _ := time.Parse("2006-01-02", req.End)
	end = end.AddDate(0, 0, 1)

	var (
		actions []analytics.Reaction
		err     error
	)

	eventList := []string{}
	if len(req.Events) > 0 {
		eventList = strings.Split(req.Events, ",")
	}

	offset := req.ItemsPerPage * (req.Page - 1)
	args := analytics.GetReactionsParams{
		OrganizationID: user.OrganizationID,
		ActionID:       uuid.FromStringOrNil(req.ActionID),
		UserID:         strings.NewReplacer(`\`, `\\`, `'`, `\'`).Replace(req.UserID),
		Events:         eventList,
		Start:          start,
		End:            end,
		Limit:          req.ItemsPerPage,
		Offset:         offset,
	}
	actions, err = rc.App.Analytics.GetReactions(args)

	if err != nil {
		rc.App.Logger.Error("actions.list.getActions", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	actionRes := []responses.Reaction{}
	for _, a := range actions {
		context := json.RawMessage("{}")
		if len(a.Context) > 0 {
			context = json.RawMessage(a.Context)
		}

		properties := json.RawMessage("{}")
		if len(a.Properties) > 0 {
			properties = json.RawMessage(a.Properties)
		}

		result := json.RawMessage("{}")
		if len(a.Result) > 0 {
			result = json.RawMessage(a.Result)
		}
		actionRes = append(actionRes, responses.Reaction{
			ID:         a.ID,
			Event:      a.Event,
			Channel:    a.Channel,
			Context:    context,
			Properties: properties,
			Result:     result,
			CreatedAt:  a.CreatedAt,
		})
	}

	count := req.ItemsPerPage * (req.Page + 1)
	if len(actionRes) < req.ItemsPerPage {
		count = req.ItemsPerPage * req.Page
	}

	res := responses.ReactionList{
		Reactions: actionRes,
		// FIX ME: Find better way to calculate number of results.
		Count: count,
	}

	c.JSON(http.StatusOK, &res)
}

func (rc ReactionsController) Delete(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	args := analytics.DeleteReactionParams{
		OrganizationID: user.OrganizationID,
		ID:             uuid.FromStringOrNil(c.Param("id")),
	}
	if err := rc.App.Analytics.DeleteReaction(args); err != nil {
		rc.App.Logger.Error("actions.delete.deleteAnalyticsAction", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}
