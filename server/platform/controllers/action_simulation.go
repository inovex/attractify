package platform

import (
	"encoding/json"
	"net/http"

	actions "attractify.io/platform/actions"
	"attractify.io/platform/app"
	"attractify.io/platform/db"
	debugging "attractify.io/platform/debugging/actionsimulation"
	"attractify.io/platform/platform/requests"
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
	c.Router.POST("/action-simulation", c.Simulate)
}

func (ac ActionSimulationController) Simulate(c *gin.Context) {
	user := c.MustGet("user").(*db.User)
	actionss, err := ac.App.DB.GetActions(c.Request.Context(), user.OrganizationID)
	if err != nil {
		ac.App.Logger.Warn("actionSimulation.list.getActions", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var req requests.ActionSimulationUser
	if err := c.ShouldBindJSON(&req); err != nil {
		ac.App.Logger.Warn("actionsimulation.simulate.parseRequest", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	as := debugging.ActionSimulation{
		User: &req,
		Action: &actions.Action{
			Targeting: &db.ActionTargeting{},
			Profile:   &db.Profile{},
			Capping:   []db.ActionCapping{},
			App:       ac.App,
			Ctx:       c,
		},
	}

	res := []responses.CheckedAction{}
	for _, a := range actionss {
		as.Action.Action = &a

		json.Unmarshal(a.Targeting, &as.Action.Targeting)
		json.Unmarshal(a.Capping, &as.Action.Capping)
		json.Unmarshal(a.TestUsers, &as.Action.TestUsers)
		steps, display := as.GetSteps()
		res = append(res, responses.CheckedAction{
			Name:    a.Name,
			Display: display,
			Steps:   steps,
		})
	}

	c.JSON(http.StatusOK, res)
}
