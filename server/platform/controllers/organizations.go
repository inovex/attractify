package platform

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"attractify.io/platform/app"
	"attractify.io/platform/auth"
	"attractify.io/platform/db"
	"attractify.io/platform/platform/requests"
	"attractify.io/platform/platform/responses"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type OrganizationController struct {
	Router *gin.RouterGroup
	App    *app.App
}

func InitOrganization(router *gin.RouterGroup, app *app.App) {
	c := OrganizationController{Router: router, App: app}
	c.Router.GET("/organization", c.Show)
	c.Router.POST("/organization", c.Create)
	c.Router.PUT("/organization", c.Update)
	c.Router.POST("/organization/token", c.Token)
	c.Router.POST("/organization/key", c.Key)
}

func (oc OrganizationController) Show(c *gin.Context) {
	user := c.MustGet("user").(*db.User)
	if user.Role != db.RoleAdmin {
		oc.App.Logger.Warn("organizations.show.userNotAdmin", zap.String("userID", user.ID.String()))
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	org, err := oc.App.DB.GetOrganization(c.Request.Context(), user.OrganizationID)
	if err != nil {
		oc.App.Logger.Error("organizations.show.getOrganization", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	res := responses.Organization{
		Name:     org.Name,
		Timezone: org.Timezone,
	}

	c.JSON(http.StatusOK, res)
}

func (oc OrganizationController) Create(c *gin.Context) {
	var req requests.OrganizationCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		oc.App.Logger.Error("organizations.create.genKey", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	org, err := oc.App.DB.CreateOrganization(
		c.Request.Context(),
		req.OrganizationName,
		req.Timezone,
		key,
	)
	if err != nil {
		oc.App.Logger.Error("organizations.create.createOrganization", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	pw := auth.NewPassword(req.Password)
	ua := db.CreateUserParams{
		OrganizationID: org.ID,
		Email:          req.Email,
		Password:       pw.Password,
		Salt:           pw.Salt,
		Name:           req.Name,
		Role:           db.RoleAdmin,
	}
	user, err := oc.App.DB.CreateUser(c.Request.Context(), ua)
	if err != nil {
		oc.App.Logger.Error("organizations.create.createUser", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	res := responses.UserSession{
		Email: user.Email,
		Name:  user.Name,
		Role:  string(user.Role),
		Token: auth.JWT("platform", user.ID, oc.App.Config.AuthKey),
	}

	c.JSON(http.StatusOK, res)
}

func (oc OrganizationController) Update(c *gin.Context) {
	user := c.MustGet("user").(*db.User)
	if user.Role != db.RoleAdmin {
		oc.App.Logger.Warn("organizations.update.userNotAdmin", zap.String("userID", user.ID.String()))
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	var req requests.OrganizationUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		oc.App.Logger.Warn("organizations.update.parseRequest", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := oc.App.DB.UpdateOrganization(c.Request.Context(), user.OrganizationID, req.Name, req.Timezone); err != nil {
		oc.App.Logger.Error("organizations.update.updateOrganization", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

func (oc OrganizationController) Token(c *gin.Context) {
	user := c.MustGet("user").(*db.User)
	if user.Role != db.RoleAdmin {
		oc.App.Logger.Warn("organizations.token.userNotAdmin", zap.String("userID", user.ID.String()))
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	token := auth.JWT("platform", user.ID, oc.App.Config.AuthKey)
	res := responses.OrganizationToken{Token: token}

	c.JSON(http.StatusOK, res)
}

func (oc OrganizationController) Key(c *gin.Context) {
	user := c.MustGet("user").(*db.User)
	if user.Role != db.RoleAdmin {
		oc.App.Logger.Warn("organizations.key.userNotAdmin", zap.String("userID", user.ID.String()))
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	var req requests.OrganizationKey
	if err := c.ShouldBindJSON(&req); err != nil {
		oc.App.Logger.Warn("organizations.key.parseRequest", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	pw := auth.Password{Password: user.Password, Salt: user.Salt}
	if !pw.Compare(req.Password) {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	key, err := oc.App.DB.GetKeyForOrganization(c.Request.Context(), user.OrganizationID)
	if err != nil {
		oc.App.Logger.Error("organizations.token.getOrganization", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	res := responses.OrganizationKey{
		Key: hex.EncodeToString(key),
	}

	c.JSON(http.StatusOK, res)
}
