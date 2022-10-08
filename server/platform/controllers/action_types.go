package platform

import (
	"encoding/json"
	"net/http"

	"attractify.io/platform/app"
	"attractify.io/platform/db"
	"attractify.io/platform/platform/requests"
	"attractify.io/platform/platform/responses"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
)

type ActionTypesController struct {
	Router *gin.RouterGroup
	App    *app.App
}

func InitActionTypes(router *gin.RouterGroup, app *app.App) {
	c := ActionTypesController{Router: router, App: app}
	c.Router.GET("/actiontypes", c.List)
	c.Router.GET("/actiontypes/:id", c.Show)
	c.Router.GET("/actiontypes/:id/used", c.InUse)
	c.Router.POST("/actiontypes", c.Create)
	c.Router.DELETE("/actiontypes/:id", c.Archive)
	c.Router.PUT("/actiontypes/:id", c.Create)
}

func (ac ActionTypesController) List(c *gin.Context) {
	user := c.MustGet("user").(*db.User)
	actiontypes, err := ac.App.DB.GetActionTypes(c.Request.Context(), user.OrganizationID)
	if err != nil {
		ac.App.Logger.Warn("actiontypes.list.getActiontypes", zap.Error(err))
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	res := []responses.ActionType{}
	for _, a := range actiontypes {
		res = append(res, responses.ActionType{
			ID:             a.ID,
			OrganizationID: a.OrganizationID,
			Name:           a.Name,
			Version:        a.Version,
			Properties:     a.Properties,
			Archived:       a.Archived,
			CreatedAt:      a.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, res)
}

func (ac ActionTypesController) Show(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	id := uuid.FromStringOrNil(c.Param("id"))
	action, err := ac.App.DB.GetActionTypeByUUID(c.Request.Context(), user.OrganizationID, id)
	if err != nil {
		ac.App.Logger.Warn("actiontypes.show.getActiontype", zap.Error(err))
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	res := responses.Action{
		ID:             action.ID,
		OrganizationID: action.OrganizationID,
		Name:           action.Name,
		Version:        action.Version,
		Properties:     action.Properties,
		CreatedAt:      action.CreatedAt,
	}

	c.JSON(http.StatusOK, res)
}

func (ac ActionTypesController) Create(c *gin.Context) {
	var req requests.ActionTypeCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		ac.App.Logger.Warn("actiontypes.create.parseRequest", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user := c.MustGet("user").(*db.User)

	properties, _ := json.Marshal(req.Properties)
	name := req.Name
	version := req.Version

	archivedTypes, err := ac.App.DB.GetActionTypesByName(c, user.OrganizationID, name)
	if err != nil {
		ac.App.Logger.Error("actiontypes.create.listArchivedVersions", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	for _, archivedType := range archivedTypes {
		ac.App.DB.UnArchiveActionType(c, user.OrganizationID, archivedType.ID)
		version = archivedType.Version + 1
	}

	args := db.CreateActionTypeParams{
		OrganizationID: user.OrganizationID,
		Name:           req.Name,
		Version:        version,
		Properties:     properties,
	}

	actionType, err := ac.App.DB.CreateActionType(c.Request.Context(), args)
	if err != nil {
		ac.App.Logger.Error("actiontypes.create.createActiontype", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	res := responses.ActionType{
		ID:             actionType.ID,
		OrganizationID: actionType.OrganizationID,
		Name:           actionType.Name,
		Version:        actionType.Version,
		Properties:     actionType.Properties,
		CreatedAt:      actionType.CreatedAt,
	}

	c.JSON(http.StatusOK, res)
}

func (ac ActionTypesController) Archive(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	name := c.Param("id")
	if err := ac.App.DB.ArchiveActionType(c.Request.Context(), user.OrganizationID, name); err != nil {
		ac.App.Logger.Error("actiontypes.archive.archiveType", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

type Usage struct {
	InUse bool `json:"inUse"`
}

func (ac ActionTypesController) InUse(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	id := uuid.FromStringOrNil(c.Param("id"))

	actionType, err := ac.App.DB.GetActionTypeByUUID(c, user.OrganizationID, id)
	if err != nil {
		ac.App.Logger.Error("actiontypes.inuse.actionTypeNotFound", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	actions, err := ac.App.DB.GetActions(c, user.OrganizationID)

	if err != nil {
		ac.App.Logger.Error("actiontypes.inuse.getActions", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	res := Usage{
		InUse: true,
	}

	for _, action := range actions {
		if action.Type == actionType.Name && action.Version == actionType.Version {
			c.JSON(http.StatusOK, res)
			return
		}
	}
	res.InUse = false

	c.JSON(http.StatusOK, res)
}
