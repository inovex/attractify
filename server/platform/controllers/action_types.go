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
	c.Router.GET("/action-types", c.List)
	c.Router.GET("/action-types/:id", c.Show)
	c.Router.GET("/action-types/exists/:name", c.Exists)
	c.Router.POST("/action-types", c.Create)
	c.Router.DELETE("/action-types/:name", c.Archive)
	c.Router.PUT("/action-types/:id", c.Update)
}

func (ac ActionTypesController) List(c *gin.Context) {
	user := c.MustGet("user").(*db.User)
	actionTypes, err := ac.App.DB.GetActionTypes(c.Request.Context(), user.OrganizationID)
	if err != nil {
		ac.App.Logger.Warn("actionTypes.list.getActiontypes", zap.Error(err))
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	res := []responses.ActionType{}
	for _, a := range actionTypes {
		res = append(res, responses.ActionType{
			ID:             a.ID,
			OrganizationID: a.OrganizationID,
			Name:           a.Name,
			Version:        a.Version,
			Properties:     a.Properties,
			IsArchived:     a.IsArchived,
			IsInUse:        a.IsInUse,
			CreatedAt:      a.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, res)
}

func (ac ActionTypesController) Show(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	id := uuid.FromStringOrNil(c.Param("id"))
	actionType, err := ac.App.DB.GetActionTypeByUUID(c.Request.Context(), user.OrganizationID, id)
	if err != nil {
		ac.App.Logger.Warn("actionTypes.show.getActiontype", zap.Error(err))
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	res := responses.ActionType{
		ID:             actionType.ID,
		OrganizationID: actionType.OrganizationID,
		Name:           actionType.Name,
		Version:        actionType.Version,
		Properties:     actionType.Properties,
		IsInUse:        actionType.IsInUse,
		CreatedAt:      actionType.CreatedAt,
	}

	c.JSON(http.StatusOK, res)
}

func (ac ActionTypesController) Create(c *gin.Context) {
	var req requests.ActionTypeCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		ac.App.Logger.Warn("actionTypes.create.parseRequest", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user := c.MustGet("user").(*db.User)

	properties, _ := json.Marshal(req.Properties)

	archivedTypes, err := ac.App.DB.GetActionTypesByName(c, user.OrganizationID, req.Name)
	if err != nil {
		ac.App.Logger.Error("actionTypes.create.listArchivedVersions", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	for _, archivedType := range archivedTypes {
		ac.App.DB.UnArchiveActionType(c, user.OrganizationID, archivedType.ID)
	}

	version, err := ac.App.DB.GetNewActionTypeVersion(c, user.OrganizationID, req.Name)
	if err != nil {
		ac.App.Logger.Error("actionTypes.create.getNewVersion", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	args := db.CreateActionTypeParams{
		OrganizationID: user.OrganizationID,
		Name:           req.Name,
		Version:        version,
		Properties:     properties,
	}

	actionType, err := ac.App.DB.CreateActionType(c.Request.Context(), args)
	if err != nil {
		ac.App.Logger.Error("actionTypes.create.createActiontype", zap.Error(err))
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

func (ac ActionTypesController) Update(c *gin.Context) {
	var req requests.ActionTypeUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		ac.App.Logger.Warn("actionTypes.update.parseRequest", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user := c.MustGet("user").(*db.User)
	id := uuid.FromStringOrNil(c.Param("id"))

	properties, _ := json.Marshal(req.Properties)

	args := db.UpdateActionTypeParams{
		ID:             id,
		OrganizationID: user.OrganizationID,
		Properties:     properties,
	}

	err := ac.App.DB.UpdateActionType(c.Request.Context(), args)
	if err != nil {
		ac.App.Logger.Error("actionTypes.update", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

func (ac ActionTypesController) Archive(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	name := c.Param("name")
	if err := ac.App.DB.ArchiveActionType(c.Request.Context(), user.OrganizationID, name); err != nil {
		ac.App.Logger.Error("actionTypes.archive.archiveType", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

func (ac ActionTypesController) Exists(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	name := c.Param("name")
	typesByName, err := ac.App.DB.GetActionTypesByName(c, user.OrganizationID, name)
	if err != nil {
		ac.App.Logger.Error("actionTypes.checkExistance.listVersions", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if len(typesByName) > 0 {
		c.JSON(http.StatusOK, true)
		return
	}

	c.JSON(http.StatusOK, false)
}
