package platform

import (
	"net/http"

	"attractify.io/platform/app"
	"attractify.io/platform/auth"
	"attractify.io/platform/db"
	"attractify.io/platform/platform/requests"
	"attractify.io/platform/platform/responses"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
)

type UsersController struct {
	Router *gin.RouterGroup
	App    *app.App
}

func InitUsers(router *gin.RouterGroup, app *app.App) {
	c := UsersController{Router: router, App: app}
	c.Router.GET("/users", c.List)
	c.Router.GET("/users/:id", c.Show)
	c.Router.POST("/users", c.Create)
	c.Router.POST("/users/:id/resend", c.ResendInvitation)
	c.Router.PUT("/users/:id", c.Update)
	c.Router.DELETE("/users/:id", c.Delete)
}

func (uc UsersController) List(c *gin.Context) {
	user := c.MustGet("user").(*db.User)
	if user.Role != db.RoleAdmin {
		uc.App.Logger.Warn("users.list.userNotAdmin", zap.String("userID", user.ID.String()))
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	users, err := uc.App.DB.GetUsersForOrganization(c.Request.Context(), user.OrganizationID)
	if err != nil {
		uc.App.Logger.Error("users.list.getUsersForOrganization", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	var res []responses.User
	for _, u := range users {
		res = append(res, responses.User{
			ID:        u.ID,
			Email:     u.Email,
			Name:      u.Name,
			Role:      string(u.Role),
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, res)
}

func (uc UsersController) Show(c *gin.Context) {
	user := c.MustGet("user").(*db.User)
	if user.Role != db.RoleAdmin {
		uc.App.Logger.Warn("users.show.userNotAdmin", zap.String("userID", user.ID.String()))
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	id := uuid.FromStringOrNil(c.Param("id"))
	u, err := uc.App.DB.GetUserForOrganization(c.Request.Context(), user.OrganizationID, id)
	if err != nil {
		uc.App.Logger.Warn("users.show.getUserForOrganization", zap.Error(err))
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	res := responses.User{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		Role:      string(u.Role),
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}

	c.JSON(http.StatusOK, res)
}

func (uc UsersController) Create(c *gin.Context) {
	user := c.MustGet("user").(*db.User)
	if user.Role != db.RoleAdmin {
		uc.App.Logger.Warn("users.create.userNotAdmin", zap.String("userID", user.ID.String()))
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	var req requests.UserCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		uc.App.Logger.Warn("users.create.parseRequest", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	args := db.CreateUserParams{
		OrganizationID: user.OrganizationID,
		Email:          req.Email,
		Role:           req.Role,
	}
	u, err := uc.App.DB.CreateUser(c.Request.Context(), args)
	if err != nil {
		uc.App.Logger.Error("users.create.createUser", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	link := auth.JWT("activation", u.ID, uc.App.Config.AuthKey)
	if err := uc.App.Mailer.SendRegistration(u.Email, user.Email, link); err != nil {
		uc.App.Logger.Error("users.create.sendRegistration", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	res := responses.User{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		Role:      string(u.Role),
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}

	c.JSON(http.StatusOK, res)
}

func (uc UsersController) ResendInvitation(c *gin.Context) {
	user := c.MustGet("user").(*db.User)
	if user.Role != db.RoleAdmin {
		uc.App.Logger.Warn("users.resend.userNotAdmin", zap.String("userID", user.ID.String()))
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	id := uuid.FromStringOrNil(c.Param("id"))
	u, err := uc.App.DB.GetUserForOrganization(c.Request.Context(), user.OrganizationID, id)
	if err != nil {
		uc.App.Logger.Warn("users.resend.getUserForOrganization", zap.Error(err))
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	link := auth.JWT("activation", u.ID, uc.App.Config.AuthKey)
	if err := uc.App.Mailer.SendRegistration(u.Email, user.Email, link); err != nil {
		uc.App.Logger.Error("users.resend.sendRegistration", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

func (uc UsersController) Update(c *gin.Context) {
	user := c.MustGet("user").(*db.User)
	if user.Role != db.RoleAdmin {
		uc.App.Logger.Warn("users.update.userNotAdmin", zap.String("userID", user.ID.String()))
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	var req requests.UserUpdateRole
	if err := c.ShouldBindJSON(&req); err != nil {
		uc.App.Logger.Warn("users.update.parseRequest", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	id := uuid.FromStringOrNil(c.Param("id"))
	if err := uc.App.DB.UpdateUserRole(c.Request.Context(), user.OrganizationID, id, req.Role); err != nil {
		uc.App.Logger.Warn("users.update.updateUsereRole", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

func (uc UsersController) Delete(c *gin.Context) {
	user := c.MustGet("user").(*db.User)
	if user.Role != db.RoleAdmin {
		uc.App.Logger.Warn("users.delete.userNotAdmin", zap.String("userID", user.ID.String()))
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	if c.Param("id") == user.ID.String() {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	id := uuid.FromStringOrNil(c.Param("id"))
	if err := uc.App.DB.DeleteUser(c.Request.Context(), user.OrganizationID, id); err != nil {
		uc.App.Logger.Warn("users.delete.deleteUser", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}
