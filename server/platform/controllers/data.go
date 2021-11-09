package platform

import (
	"encoding/json"
	"net/http"

	"attractify.io/platform/app"
	"attractify.io/platform/db"
	"attractify.io/platform/platform/requests"
	"attractify.io/platform/platform/responses"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gopkg.in/square/go-jose.v2"
)

type DataController struct {
	Router *gin.RouterGroup
	App    *app.App
}

func InitData(router *gin.RouterGroup, app *app.App) {
	d := DataController{Router: router, App: app}
	d.Router.POST("/data", d.Encrypt)
}

func (dc DataController) Encrypt(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	var req requests.Data
	if err := c.ShouldBindJSON(&req); err != nil {
		dc.App.Logger.Warn("data.encrypt.parseRequest", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	key, err := dc.App.DB.GetKeyForOrganization(c.Request.Context(), user.OrganizationID)
	if err != nil {
		dc.App.Logger.Error("data.encrypt.getKeyForOrganization", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	encrypter, err := jose.NewEncrypter(jose.A256GCM,
		jose.Recipient{Algorithm: jose.DIRECT, Key: key}, nil)
	if err != nil {
		dc.App.Logger.Error("data.encrypt.newEncrypter", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	plaintext, err := json.Marshal(req.Data)
	if err != nil {
		dc.App.Logger.Error("data.encrypt.marshal", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	object, err := encrypter.Encrypt(plaintext)
	if err != nil {
		dc.App.Logger.Error("data.encrypt.encrypt", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	serialized, err := object.CompactSerialize()
	if err != nil {
		dc.App.Logger.Error("data.encrypt.compactSerialize", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	res := responses.Data{Data: serialized}
	c.JSON(http.StatusOK, res)
}
