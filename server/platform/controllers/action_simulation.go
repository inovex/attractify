package platform

import (
	"net/http"

	"attractify.io/platform/app"
	"attractify.io/platform/db"
	"attractify.io/platform/platform/responses"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ActionSimulationController struct {
	Router *gin.RouterGroup
	App    *app.App
}

func InitActionSimulation(router *gin.RouterGroup, app *app.App) {
	c := ActionSimulationController{Router: router, App: app}
	c.Router.GET("/action-simulation", c.Simulate)
}

func (ac ActionSimulationController) Simulate(c *gin.Context) {
	user := c.MustGet("user").(*db.User)
	actions, err := ac.App.DB.GetActions(c.Request.Context(), user.OrganizationID)
	if err != nil {
		ac.App.Logger.Warn("actionSimulation.list.getActions", zap.Error(err))
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	res := []responses.CheckedAction{}
	for _, a := range actions {
		state := "a"
		res = append(res, responses.CheckedAction{
			Name:  a.Name,
			State: state,
		})
	}

	c.JSON(http.StatusOK, res)
}

func (ac ActionsController) TestWebhookA(c *gin.Context) {
	// var req requests.ActionWebhookTest
	// if err := c.ShouldBindJSON(&req); err != nil {
	// 	ac.App.Logger.Warn("actions.testWebhook.parseRequest", zap.Error(err))
	// 	c.AbortWithStatus(http.StatusBadRequest)
	// 	return
	// }

	// user := c.MustGet("user").(*db.User)
	// args := db.GetActionParams{
	// 	OrganizationID: user.OrganizationID,
	// 	ID:             uuid.FromStringOrNil(c.Param("id")),
	// }
	// campaign, err := ac.App.DB.GetAction(c.Request.Context(), args)
	// if err != nil {
	// 	ac.App.Logger.Warn("actions.testWebhook.getAction", zap.Error(err))
	// 	c.AbortWithStatus(http.StatusBadRequest)
	// 	return
	// }

	// var actions []db.Action
	// if err := json.Unmarshal(campaign.Actions, &actions); err != nil {
	// 	ac.App.Logger.Error("actions.testWebhook.unmarshalActions", zap.Error(err))
	// 	c.AbortWithStatus(http.StatusInternalServerError)
	// 	return
	// }

	// if len(actions.Webhook.URL) == 0 {
	// 	c.AbortWithStatus(http.StatusBadRequest)
	// 	return
	// }

	// payload := webhook.Payload{
	// 	Event:          req.Event,
	// 	ActionID:       campaign.ID,
	// 	UserID:         req.UserID,
	// 	Channel:        req.Channel,
	// 	UserProperties: json.RawMessage(req.Properties),
	// 	Timestamp:      time.Now(),
	// }

	// if len(actions.Webhook.Properties) > 0 {
	// 	payload.CustomProperties = json.RawMessage(actions.Webhook.Properties)
	// }

	// key, err := ac.App.DB.GetKey(c.Request.Context(), user.OrganizationID)
	// if err != nil {
	// 	ac.App.Logger.Error("actions.testWebhook.getKey", zap.Error(err))
	// 	c.AbortWithStatus(http.StatusInternalServerError)
	// 	return
	// }

	// res, err := webhook.Execute(actions.Webhook.URL, payload, key)
	// if err != nil {
	// 	ac.App.Logger.Warn("actions.testWebhook.execute", zap.Error(err))
	// 	c.AbortWithStatus(http.StatusBadRequest)
	// 	return
	// }

	// c.JSON(http.StatusOK, res)
}
