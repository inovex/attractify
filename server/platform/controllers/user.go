package platform

import (
	"net/http"
	"time"

	"attractify.io/platform/app"
	"attractify.io/platform/auth"
	"attractify.io/platform/db"
	"attractify.io/platform/platform/requests"
	"attractify.io/platform/platform/responses"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserController struct {
	Router *gin.RouterGroup
	App    *app.App
}

func InitUser(router *gin.RouterGroup, app *app.App) {
	c := UserController{Router: router, App: app}
	c.Router.GET("/user", c.Show)
	c.Router.POST("/user", c.Activate)
	c.Router.POST("/user/session", c.SignIn)
	c.Router.DELETE("/user/session", c.SignOut)
	c.Router.POST("/user/reset-password", c.ResetPassword)
	c.Router.PUT("/user/reset-password", c.SetNewPassword)
	c.Router.PUT("/user", c.Update)
	c.Router.PUT("/user/password", c.UpdatePassword)
}

func (uc UserController) Show(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	res := responses.User{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
		Role:  string(user.Role),
	}

	c.JSON(http.StatusOK, res)
}

func (uc UserController) SignIn(c *gin.Context) {
	var req requests.UserSession
	if err := c.ShouldBindJSON(&req); err != nil {
		uc.App.Logger.Warn("users.signIn.parseRequest", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user, err := uc.App.DB.GetUserByEmail(c.Request.Context(), req.Email)
	if err != nil {
		uc.App.Logger.Warn("users.signIn.getUserByEmail", zap.Error(err))
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	pw := auth.Password{Password: user.Password, Salt: user.Salt}
	if !pw.Compare(req.Password) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	res := responses.UserSession{
		Email: user.Email,
		Name:  user.Name,
		Role:  string(user.Role),
		Token: auth.JWT("platform", user.ID, uc.App.Config.AuthKey),
	}

	c.JSON(http.StatusOK, res)
}

func (uc UserController) SignOut(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	if err := uc.App.DB.UpdateUserLoggedOutAt(c.Request.Context(), user.OrganizationID, user.ID, time.Now().UTC()); err != nil {
		uc.App.Logger.Warn("users.signOut.updateLoggedOutAt", zap.Error(err))
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

func (uc UserController) Activate(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	var req requests.UserActivate
	if err := c.ShouldBindJSON(&req); err != nil {
		uc.App.Logger.Warn("users.activate.parseRequest", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	pw := auth.NewPassword(req.Password)
	args := db.UpdateUserPasswordAndNameParams{
		ID:       user.ID,
		Password: pw.Password,
		Salt:     pw.Salt,
		Name:     req.Name,
	}
	if err := uc.App.DB.UpdateUserPasswordAndName(c.Request.Context(), args); err != nil {
		uc.App.Logger.Error("users.activate.updateUserPasswordAndName", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

func (uc UserController) ResetPassword(c *gin.Context) {
	var req requests.UserResetPassword
	if err := c.ShouldBindJSON(&req); err != nil {
		uc.App.Logger.Warn("users.resetPassword.parseRequest", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user, err := uc.App.DB.GetUserByEmail(c.Request.Context(), req.Email)
	if err != nil {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}

	link := auth.JWT("reset_password", user.ID, uc.App.Config.AuthKey)
	if err := uc.App.Mailer.SendPasswordReset(user.Email, link); err != nil {
		uc.App.Logger.Error("users.resetPassword.sendPasswordReset", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

func (uc UserController) SetNewPassword(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	var req requests.UserUpdatePassword
	if err := c.ShouldBindJSON(&req); err != nil {
		uc.App.Logger.Warn("users.updatePassword.parseRequest", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	pw := auth.NewPassword(req.Password)
	if err := uc.App.DB.UpdateUserPassword(c.Request.Context(), user.ID, pw.Password, pw.Salt); err != nil {
		uc.App.Logger.Error("users.updatePassword.updateUserPassword", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

func (uc UserController) Update(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	var req requests.UserUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		uc.App.Logger.Warn("users.update.parseRequest", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	args := db.UpdateUserPropertiesParams{
		OrganizationID: user.OrganizationID,
		ID:             user.ID,
		Email:          req.Email,
		Name:           req.Name,
	}
	if err := uc.App.DB.UpdateUserProperties(c.Request.Context(), args); err != nil {
		uc.App.Logger.Error("users.update.updateUserProperties", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

func (uc UserController) UpdatePassword(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	var req requests.UserUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		uc.App.Logger.Warn("users.updatePassword.parseRequest", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	oldPw := auth.Password{Password: user.Password, Salt: user.Salt}
	if !oldPw.Compare(req.OldPassword) {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	newPw := auth.NewPassword(req.NewPassword)
	if err := uc.App.DB.UpdateUserPassword(c.Request.Context(), user.ID, newPw.Password, newPw.Salt); err != nil {
		uc.App.Logger.Error("users.updatePassword.updateUserPassword", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}
