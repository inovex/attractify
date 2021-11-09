package platform

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"attractify.io/platform/app"
	"attractify.io/platform/db"
	"attractify.io/platform/platform/requests"
	"attractify.io/platform/platform/responses"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
)

var chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_.")

type AuthTokensController struct {
	Router *gin.RouterGroup
	App    *app.App
}

func InitAuthTokens(router *gin.RouterGroup, app *app.App) {
	c := AuthTokensController{Router: router, App: app}
	c.Router.GET("/auth-tokens", c.List)
	c.Router.POST("/auth-tokens", c.Create)
	c.Router.DELETE("/auth-tokens/:id", c.Delete)
}

func (ac AuthTokensController) List(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	authTokens, err := ac.App.DB.GetAuthTokens(c.Request.Context(), user.OrganizationID)
	if err != nil {
		ac.App.Logger.Error("authTokens.list.getAuthTokens", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	res := []responses.AuthToken{}
	for _, t := range authTokens {
		res = append(res, responses.AuthToken{
			ID:        t.ID,
			Channel:   t.Channel,
			Token:     t.Token,
			CreatedAt: t.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, res)
}

func (ac AuthTokensController) Create(c *gin.Context) {
	var req requests.AuthToken
	if err := c.ShouldBindJSON(&req); err != nil {
		ac.App.Logger.Warn("authTokens.create.parseRequest", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user := c.MustGet("user").(*db.User)

	rand.Seed(time.Now().UnixNano())
	token := fmt.Sprintf("%s-%s", req.Channel, ac.randSeq(64))

	args := db.CreateAuthTokenParams{
		OrganizationID: user.OrganizationID,
		Channel:        req.Channel,
		Token:          token,
	}
	t, err := ac.App.DB.CreateAuthToken(c.Request.Context(), args)
	if err != nil {
		ac.App.Logger.Error("authTokens.create.createAuthToken", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	res := responses.AuthToken{
		ID:        t.ID,
		Channel:   t.Channel,
		Token:     t.Token,
		CreatedAt: t.CreatedAt,
	}

	c.JSON(http.StatusOK, res)
}

func (ac AuthTokensController) Delete(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	id := uuid.FromStringOrNil(c.Param("id"))
	if err := ac.App.DB.DeleteAuthToken(c.Request.Context(), user.OrganizationID, id); err != nil {
		ac.App.Logger.Warn("authTokens.delete.deleteAuthToken", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

func (ac AuthTokensController) randSeq(n int) string {
	b := make([]rune, n)
	cl := len(chars)
	for i := range b {
		b[i] = chars[rand.Intn(cl)]
	}
	return string(b)
}
