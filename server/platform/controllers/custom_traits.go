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
	"go.uber.org/zap"
)

type CustomTraitsController struct {
	Router *gin.RouterGroup
	App    *app.App
}

func InitCustomTraits(router *gin.RouterGroup, app *app.App) {
	c := CustomTraitsController{Router: router, App: app}
	c.Router.GET("/custom-traits", c.Show)
	c.Router.POST("/custom-traits", c.Upsert)
	c.Router.GET("/custom-traits/properties", c.ListProperties)
}

func (cc CustomTraitsController) Show(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	ct, err := cc.App.DB.GetCustomTraits(c.Request.Context(), user.OrganizationID)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	res := responses.CustomTraits{
		OrganizationID: ct.OrganizationID,
		Structure:      ct.Structure,
		CreatedAt:      ct.CreatedAt,
		UpdatedAt:      ct.UpdatedAt,
	}

	c.JSON(http.StatusOK, res)
}

func (cc CustomTraitsController) Upsert(c *gin.Context) {
	var req requests.CustomTraitsUpsert
	if err := c.ShouldBindJSON(&req); err != nil {
		cc.App.Logger.Warn("customTraits.upsert.parseRequest", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user := c.MustGet("user").(*db.User)

	def := events.NewDefinition()
	jsonSchema, properties, err := def.Prepare(req.Structure)
	if err != nil {
		cc.App.Logger.Warn("customTraits.upsert.prepare", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	args := db.UpsertCustomTraitsParams{
		OrganizationID: user.OrganizationID,
		Structure:      req.Structure,
		JSONSchema:     jsonSchema,
		Properties:     properties,
	}
	ct, err := cc.App.DB.UpsertCustomTraits(c.Request.Context(), args)
	if err != nil {
		cc.App.Logger.Error("customTraits.upsert.upsertCustomTraits", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	res := responses.Event{
		OrganizationID: ct.OrganizationID,
		Structure:      ct.Structure,
		CreatedAt:      ct.CreatedAt,
		UpdatedAt:      ct.UpdatedAt,
	}

	c.JSON(http.StatusOK, res)
}

func (cc CustomTraitsController) ListProperties(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	ctp, err := cc.App.DB.GetCustomTraitsProperties(c.Request.Context(), user.OrganizationID)
	if err != nil {
		cc.App.Logger.Error("customTraits.listProperties.getCustomTraitsProperties", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var properties []db.CustomTraitsProperty
	if err := json.Unmarshal(ctp, &properties); err != nil {
		cc.App.Logger.Error("customTraits.listProperties.prepareResult", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	res := []responses.CustomTraitsProperty{}
	for _, p := range properties {
		res = append(res, responses.CustomTraitsProperty{
			Key:  p.Key,
			Type: p.Type,
		})
	}

	c.JSON(http.StatusOK, res)
}
