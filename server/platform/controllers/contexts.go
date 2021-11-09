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

type ContextsController struct {
	Router *gin.RouterGroup
	App    *app.App
}

func InitContexts(router *gin.RouterGroup, app *app.App) {
	c := ContextsController{Router: router, App: app}
	c.Router.GET("/contexts", c.List)
	c.Router.GET("/contexts/:id", c.Show)
	c.Router.POST("/contexts", c.Create)
	c.Router.DELETE("/contexts/:id", c.Delete)
	c.Router.PUT("/contexts/:id", c.Update)
	c.Router.GET("/contexts/:id/properties", c.ListProperties)
}
func (cc ContextsController) List(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	cd, err := cc.App.DB.GetContexts(c.Request.Context(), user.OrganizationID)
	if err != nil {
		cc.App.Logger.Error("contexts.list.getContexts", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	res := []responses.Context{}
	for _, c := range cd {
		res = append(res, responses.Context{
			ID:             c.ID,
			OrganizationID: c.OrganizationID,
			Channel:        c.Channel,
			Structure:      c.Structure,
			Properties:     c.Properties,
			CreatedAt:      c.CreatedAt,
			UpdatedAt:      c.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, res)
}

func (cc ContextsController) Show(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	id := uuid.FromStringOrNil(c.Param("id"))
	cd, err := cc.App.DB.GetContextByID(c.Request.Context(), user.OrganizationID, id)
	if err != nil {
		cc.App.Logger.Warn("contexts.show.getContextByID", zap.Error(err))
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	// FIXME: Extract to standalone methode.
	res := responses.Context{
		ID:             cd.ID,
		OrganizationID: cd.OrganizationID,
		Channel:        cd.Channel,
		Structure:      cd.Structure,
		Properties:     cd.Properties,
		CreatedAt:      cd.CreatedAt,
		UpdatedAt:      cd.UpdatedAt,
	}

	c.JSON(http.StatusOK, res)
}

func (cc ContextsController) Create(c *gin.Context) {
	var req requests.Context
	if err := c.ShouldBindJSON(&req); err != nil {
		cc.App.Logger.Warn("contexts.create.parseRequest", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user := c.MustGet("user").(*db.User)

	def := events.NewDefinition()
	jsonSchema, properties, err := def.Prepare(req.Structure)
	if err != nil {
		cc.App.Logger.Warn("contexts.create.prepare", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	args := db.CreateContextParams{
		OrganizationID: user.OrganizationID,
		Channel:        req.Channel,
		Structure:      req.Structure,
		JSONSchema:     jsonSchema,
		Properties:     properties,
	}
	cd, err := cc.App.DB.CreateContext(c.Request.Context(), args)
	if err != nil {
		cc.App.Logger.Error("contexts.create.createContext", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// FIXME: Extract to standalone methode.
	res := responses.Context{
		ID:             cd.ID,
		OrganizationID: cd.OrganizationID,
		Channel:        cd.Channel,
		Structure:      cd.Structure,
		Properties:     cd.Properties,
		CreatedAt:      cd.CreatedAt,
		UpdatedAt:      cd.UpdatedAt,
	}

	c.JSON(http.StatusOK, res)
}

func (cc ContextsController) Update(c *gin.Context) {
	var req requests.Context
	if err := c.ShouldBindJSON(&req); err != nil {
		cc.App.Logger.Warn("contexts.update.parseRequest", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user := c.MustGet("user").(*db.User)

	def := events.NewDefinition()
	jsonSchema, properties, err := def.Prepare(req.Structure)
	if err != nil {
		cc.App.Logger.Warn("contexts.create.prepare", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	id := uuid.FromStringOrNil(c.Param("id"))
	args := db.UpdateContextParams{
		OrganizationID: user.OrganizationID,
		ID:             id,
		Structure:      req.Structure,
		JSONSchema:     jsonSchema,
		Properties:     properties,
	}

	if err := cc.App.DB.UpdateContext(c.Request.Context(), args); err != nil {
		cc.App.Logger.Error("contexts.update.updateContext", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

func (cc ContextsController) Delete(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	if err := cc.App.DB.DeleteContext(c.Request.Context(), user.OrganizationID, uuid.FromStringOrNil(c.Param("id"))); err != nil {
		cc.App.Logger.Warn("contexts.delete.deleteContext", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

func (cc ContextsController) ListProperties(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	ctp, err := cc.App.DB.GetContextProperties(c.Request.Context(), user.OrganizationID, c.Param("id"))
	if err != nil {
		cc.App.Logger.Error("contexts.listProperties.getContextProperties", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var properties []db.ContextProperty
	if err := json.Unmarshal(ctp, &properties); err != nil {
		cc.App.Logger.Error("contexts.listProperties.prepareResult", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	res := []responses.ContextProperty{}
	for _, p := range properties {
		res = append(res, responses.ContextProperty{
			Key:  p.Key,
			Type: p.Type,
		})
	}

	c.JSON(http.StatusOK, res)
}
