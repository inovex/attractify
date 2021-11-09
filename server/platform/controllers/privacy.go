package platform

import (
	"net/http"

	"attractify.io/platform/app"
	"attractify.io/platform/db"
	"attractify.io/platform/platform/requests"
	"attractify.io/platform/platform/responses"
	"attractify.io/platform/privacy"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
)

type PrivacyController struct {
	Router *gin.RouterGroup
	App    *app.App
}

func InitPrivacy(router *gin.RouterGroup, app *app.App) {
	c := PrivacyController{Router: router, App: app}
	c.Router.POST("/privacy/export", c.Export)
	c.Router.POST("/privacy/deletion", c.Deletion)
	c.Router.GET("/privacy/locked-identities", c.GetLockedIdentity)
	c.Router.POST("/privacy/locked-identities", c.CreateLockedIdentity)
	c.Router.DELETE("/privacy/locked-identities/:id", c.DeleteLockedIdentity)
}

func (pc PrivacyController) Export(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	var req requests.Privacy
	if err := c.ShouldBindJSON(&req); err != nil {
		pc.App.Logger.Warn("privacy.export.parseRequest", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	e := privacy.NewExport(pc.App, user.OrganizationID, req.UserID, req.Email)
	if err := e.Run(); err != nil {
		pc.App.Logger.Warn("privacy.export.run", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

func (pc PrivacyController) Deletion(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	var req requests.Privacy
	if err := c.ShouldBindJSON(&req); err != nil {
		pc.App.Logger.Warn("privacy.delete.parseRequest", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	e := privacy.NewDeletionByUserID(pc.App, user.OrganizationID, req.UserID)
	if err := e.Run(); err != nil {
		pc.App.Logger.Warn("privacy.delete.run", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

func (pc PrivacyController) GetLockedIdentity(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	li, err := pc.App.DB.GetLockedProfileIdentities(c.Request.Context(), user.OrganizationID)
	if err != nil {
		pc.App.Logger.Warn("privacy.getLockedIdentity.getLockedProfileIdentity", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	res := []responses.Privacy{}
	for _, l := range li {
		res = append(res, responses.Privacy{
			ID:        l.ID,
			UserID:    l.UserID,
			CreatedAt: l.CreatedAt,
		})
	}
	c.JSON(http.StatusOK, res)
}

func (pc PrivacyController) CreateLockedIdentity(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	var req requests.Privacy
	if err := c.ShouldBindJSON(&req); err != nil {
		pc.App.Logger.Warn("privacy.delete.parseRequest", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	li, err := pc.App.DB.CreateLockedProfileIdentity(c.Request.Context(), user.OrganizationID, req.UserID)
	if err != nil {
		pc.App.Logger.Warn("privacy.createLockedIdentity.createLockedProfileIdentity", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	res := responses.Privacy{
		ID:        li.ID,
		UserID:    li.UserID,
		CreatedAt: li.CreatedAt,
	}
	c.JSON(http.StatusOK, res)
}

func (pc PrivacyController) DeleteLockedIdentity(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	id := uuid.FromStringOrNil(c.Param("id"))
	if err := pc.App.DB.DeleteLockedProfileIdentitiesByID(c.Request.Context(), user.OrganizationID, id); err != nil {
		pc.App.Logger.Warn("privacy.deleteLockedIdentity.deleteLockedProfileIdentity", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}
