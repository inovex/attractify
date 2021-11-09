package platform

import (
	"encoding/json"
	"net/http"

	"attractify.io/platform/app"
	"attractify.io/platform/db"
	"attractify.io/platform/events"
	"attractify.io/platform/platform/requests"
	"attractify.io/platform/platform/responses"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
)

type EventsController struct {
	Router *gin.RouterGroup
	App    *app.App
}

func InitEvents(router *gin.RouterGroup, app *app.App) {
	c := EventsController{Router: router, App: app}
	c.Router.GET("/events", c.List)
	c.Router.GET("/events/:id", c.Show)
	c.Router.POST("/events", c.Create)
	c.Router.DELETE("/events/:id", c.Delete)
	c.Router.PUT("/events/:id", c.Update)
	c.Router.GET("/events/:id/properties", c.ListEventProperties)
}

func (ec EventsController) List(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	ed, err := ec.App.DB.GetEvents(c.Request.Context(), user.OrganizationID)
	if err != nil {
		ec.App.Logger.Error("events.list.getEvents", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	res := []responses.Event{}
	for _, t := range ed {
		res = append(res, responses.Event{
			ID:             t.ID,
			OrganizationID: t.OrganizationID,
			Name:           t.Name,
			Description:    t.Description,
			Structure:      t.Structure,
			Properties:     t.Properties,
			CreatedAt:      t.CreatedAt,
			UpdatedAt:      t.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, res)
}

func (ec EventsController) Show(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	id := uuid.FromStringOrNil(c.Param("id"))
	ed, err := ec.App.DB.GetEvent(c.Request.Context(), user.OrganizationID, id)
	if err != nil {
		ec.App.Logger.Warn("events.show.getEvent", zap.Error(err))
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	// FIXME: Extract to standalone methode.
	res := responses.Event{
		ID:             ed.ID,
		OrganizationID: ed.OrganizationID,
		Name:           ed.Name,
		Description:    ed.Description,
		Structure:      ed.Structure,
		CreatedAt:      ed.CreatedAt,
		UpdatedAt:      ed.UpdatedAt,
	}

	c.JSON(http.StatusOK, res)
}

func (ec EventsController) Create(c *gin.Context) {
	var req requests.EventCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		ec.App.Logger.Warn("events.create.parseRequest", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user := c.MustGet("user").(*db.User)

	def := events.NewDefinition()
	jsonSchema, properties, err := def.Prepare(req.Structure)
	if err != nil {
		ec.App.Logger.Warn("events.create.prepare", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	args := db.CreateEventParams{
		OrganizationID: user.OrganizationID,
		Name:           req.Name,
		Description:    req.Description,
		Structure:      req.Structure,
		JSONSchema:     jsonSchema,
		Properties:     properties,
	}
	ed, err := ec.App.DB.CreateEvent(c.Request.Context(), args)
	if err != nil {
		ec.App.Logger.Error("events.create.createEvent", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// FIXME: Extract to standalone methode.
	res := responses.Event{
		ID:             ed.ID,
		OrganizationID: ed.OrganizationID,
		Name:           ed.Name,
		Description:    ed.Description,
		Structure:      ed.Structure,
		CreatedAt:      ed.CreatedAt,
		UpdatedAt:      ed.UpdatedAt,
	}

	c.JSON(http.StatusOK, res)
}

func (ec EventsController) Update(c *gin.Context) {
	var req requests.EventCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		ec.App.Logger.Warn("events.update.parseRequest", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user := c.MustGet("user").(*db.User)

	def := events.NewDefinition()
	jsonSchema, properties, err := def.Prepare(req.Structure)
	if err != nil {
		ec.App.Logger.Warn("events.create.prepare", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	args := db.UpdateEventParams{
		OrganizationID: user.OrganizationID,
		Name:           req.Name,
		Description:    req.Description,
		Structure:      req.Structure,
		JSONSchema:     jsonSchema,
		Properties:     properties,
		ID:             uuid.FromStringOrNil(c.Param("id")),
	}

	if err := ec.App.DB.UpdateEvent(c.Request.Context(), args); err != nil {
		ec.App.Logger.Error("events.update.updateEvent", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

func (ec EventsController) Delete(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	id := uuid.FromStringOrNil(c.Param("id"))
	if err := ec.App.DB.DeleteEvent(c.Request.Context(), user.OrganizationID, id); err != nil {
		ec.App.Logger.Warn("events.delete.deleteEvent", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

func (ec EventsController) ListEventNames(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	events, err := ec.App.DB.GetEventEvents(c.Request.Context(), user.OrganizationID)
	if err != nil {
		ec.App.Logger.Error("events.searchEvents.getEventEvents", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	res := []responses.EventOverview{}
	for _, e := range events {
		// FIXME: Extract to standalone methode.
		res = append(res, responses.EventOverview{
			ID:   e.ID,
			Name: e.Name,
		})
	}

	c.JSON(http.StatusOK, res)
}

func (ec EventsController) ListEventProperties(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	id := uuid.FromStringOrNil(c.Param("id"))
	ed, err := ec.App.DB.GetEventProperties(c.Request.Context(), user.OrganizationID, id)
	if err != nil {
		ec.App.Logger.Error("events.searchProperties.getEventProperties", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var properties []db.EventProperty
	if err := json.Unmarshal(ed, &properties); err != nil {
		ec.App.Logger.Error("events.searchProperties.prepareResult", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	res := []responses.EventProperty{}
	for _, p := range properties {
		res = append(res, responses.EventProperty{
			Key:  p.Key,
			Type: p.Type,
		})
	}

	c.JSON(http.StatusOK, res)
}
