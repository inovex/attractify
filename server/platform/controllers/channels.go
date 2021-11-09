package platform

import (
	"net/http"

	"attractify.io/platform/app"
	"attractify.io/platform/db"
	"attractify.io/platform/platform/requests"
	"attractify.io/platform/platform/responses"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
)

type ChannelsController struct {
	Router *gin.RouterGroup
	App    *app.App
}

func InitChannels(router *gin.RouterGroup, app *app.App) {
	c := ChannelsController{Router: router, App: app}
	c.Router.GET("/channels", c.List)
	c.Router.GET("/channels/:id", c.Show)
	c.Router.POST("/channels", c.Create)
	c.Router.DELETE("/channels/:id", c.Delete)
	c.Router.PUT("/channels/:id", c.Update)
}

func (cc ChannelsController) List(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	channels, err := cc.App.DB.GetChannels(c.Request.Context(), user.OrganizationID)
	if err != nil {
		cc.App.Logger.Error("channels.list.getChannels", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	res := []responses.Channel{}
	for _, t := range channels {
		res = append(res, responses.Channel{
			ID:        t.ID,
			Name:      t.Name,
			Key:       t.Key,
			CreatedAt: t.CreatedAt,
			UpdatedAt: t.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, res)
}

func (cc ChannelsController) Show(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	id := uuid.FromStringOrNil(c.Param("id"))
	channel, err := cc.App.DB.GetChannel(c.Request.Context(), user.OrganizationID, id)
	if err != nil {
		cc.App.Logger.Warn("channels.show.getChannel", zap.Error(err))
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, channel)
}

func (cc ChannelsController) Create(c *gin.Context) {
	var req requests.Channel
	if err := c.ShouldBindJSON(&req); err != nil {
		cc.App.Logger.Warn("channels.create.parseRequest", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user := c.MustGet("user").(*db.User)

	args := db.CreateChannelParams{
		OrganizationID: user.OrganizationID,
		Name:           req.Name,
		Key:            req.Key,
	}
	channel, err := cc.App.DB.CreateChannel(c.Request.Context(), args)
	if err != nil {
		cc.App.Logger.Error("channels.create.createChannel", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	res := responses.Channel{
		ID:        channel.ID,
		Name:      channel.Name,
		Key:       channel.Key,
		CreatedAt: channel.CreatedAt,
		UpdatedAt: channel.UpdatedAt,
	}

	c.JSON(http.StatusOK, res)
}

func (cc ChannelsController) Update(c *gin.Context) {
	var req requests.Channel
	if err := c.ShouldBindJSON(&req); err != nil {
		cc.App.Logger.Warn("channels.update.parseRequest", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user := c.MustGet("user").(*db.User)

	args := db.UpdateChannelParams{
		OrganizationID: user.OrganizationID,
		Name:           req.Name,
		Key:            req.Key,
		ID:             uuid.FromStringOrNil(c.Param("id")),
	}

	if err := cc.App.DB.UpdateChannel(c.Request.Context(), args); err != nil {
		cc.App.Logger.Error("channels.update.updateChannel", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

func (cc ChannelsController) Delete(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	id := uuid.FromStringOrNil(c.Param("id"))
	if err := cc.App.DB.DeleteChannel(c.Request.Context(), user.OrganizationID, id); err != nil {
		cc.App.Logger.Warn("channels.delete.deleteChannel", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}
