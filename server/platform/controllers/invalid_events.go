package platform

import (
	"net/http"

	"attractify.io/platform/app"
	"attractify.io/platform/db"
	"attractify.io/platform/platform/requests"
	"attractify.io/platform/platform/responses"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
)

type InvalidEventsController struct {
	Router *gin.RouterGroup
	App    *app.App
}

func InitInvalidEvents(router *gin.RouterGroup, app *app.App) {
	c := InvalidEventsController{Router: router, App: app}
	c.Router.GET("/invalid-events", c.List)
	c.Router.DELETE("/invalid-events/:id", c.Delete)
	c.Router.PUT("/invalid-events/update", c.Update)
}

func (ec InvalidEventsController) List(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	ed, err := ec.App.DB.GetInvalidEvents(c.Request.Context(), user.OrganizationID)
	if err != nil {
		ec.App.Logger.Error("Invalidevents.list.getInvalidEvents", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	res := []responses.InvalidEvent{}
	for _, t := range ed {
		res = append(res, responses.InvalidEvent{
			ID:             t.ID,
			EventID:        t.EventID,
			OrganizationID: t.OrganizationID,
			Channel:        t.Channel,
			Name:           t.Name,
			Properties:     t.Properties,
			Context:        t.Context,
			Type:           t.Type,
			CreatedAt:      t.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, res)
}

func (ec InvalidEventsController) Delete(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	id := uuid.FromStringOrNil(c.Param("id"))
	if err := ec.App.DB.DeleteInvalidEvent(c.Request.Context(), user.OrganizationID, id); err != nil {
		ec.App.Logger.Warn("invalidEvents.delete.deleteInvalidEvent", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

func (ec InvalidEventsController) Update(c *gin.Context) {
	var req requests.InvalidEventUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		ec.App.Logger.Warn("invalidEvents.update.parseRequest", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	user := c.MustGet("user").(*db.User)

	id := uuid.FromStringOrNil(req.EventId)
	newName := req.NewName
	if err := ec.App.DB.UpdateInvalidEvent(c.Request.Context(), newName, user.OrganizationID, id); err != nil {
		ec.App.Logger.Warn("invalidEvents.update.updateInvalidEvent", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}
