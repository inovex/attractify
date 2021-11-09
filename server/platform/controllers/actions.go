package platform

import (
	"encoding/json"
	"fmt"
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

type WebhookPayload struct {
	Event            string          `json:"event"`
	ActionID         uuid.UUID       `json:"actionId"`
	UserID           string          `json:"userId"`
	Channel          string          `json:"channel"`
	UserProperties   json.RawMessage `json:"userProperties,omitempty"`
	CustomProperties json.RawMessage `json:"customProperties,omitempty"`
	Timestamp        time.Time       `json:"timestamp"`
}

type ActionsController struct {
	Router *gin.RouterGroup
	App    *app.App
}

func InitActions(router *gin.RouterGroup, app *app.App) {
	c := ActionsController{Router: router, App: app}
	c.Router.GET("/actions", c.List)
	c.Router.GET("/actions/:id", c.Show)
	c.Router.POST("/actions", c.Create)
	c.Router.POST("/actions/:id/duplicate", c.Duplicate)
	c.Router.DELETE("/actions/:id", c.Delete)
	c.Router.PUT("/actions/:id", c.Update)
	c.Router.PUT("/actions/:id/state", c.UpdateState)
	c.Router.POST("/actions/:id/test-webhook", c.TestWebhook)
}

func (ac ActionsController) List(c *gin.Context) {
	user := c.MustGet("user").(*db.User)
	actions, err := ac.App.DB.GetActions(c.Request.Context(), user.OrganizationID)
	if err != nil {
		ac.App.Logger.Warn("actions.list.getActions", zap.Error(err))
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	res := []responses.Action{}
	for _, a := range actions {
		res = append(res, responses.Action{
			ID:             a.ID,
			OrganizationID: a.OrganizationID,
			Name:           a.Name,
			Type:           a.Type,
			State:          string(a.State),
			Tags:           a.Tags,
			Properties:     a.Properties,
			Targeting:      a.Targeting,
			Capping:        a.Capping,
			Hooks:          a.Hooks,
			TestUsers:      a.TestUsers,
			CreatedAt:      a.CreatedAt,
			UpdatedAt:      a.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, res)
}

func (ac ActionsController) Show(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	id := uuid.FromStringOrNil(c.Param("id"))
	action, err := ac.App.DB.GetAction(c.Request.Context(), user.OrganizationID, id)
	if err != nil {
		ac.App.Logger.Warn("actions.show.getAction", zap.Error(err))
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	res := responses.Action{
		ID:             action.ID,
		OrganizationID: action.OrganizationID,
		Name:           action.Name,
		Type:           action.Type,
		State:          string(action.State),
		Tags:           action.Tags,
		Properties:     action.Properties,
		Targeting:      action.Targeting,
		Capping:        action.Capping,
		Hooks:          action.Hooks,
		TestUsers:      action.TestUsers,
		CreatedAt:      action.CreatedAt,
		UpdatedAt:      action.UpdatedAt,
	}

	c.JSON(http.StatusOK, res)
}

func (ac ActionsController) Create(c *gin.Context) {
	var req requests.ActionCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		ac.App.Logger.Warn("actions.create.parseRequest", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user := c.MustGet("user").(*db.User)
	tags, _ := json.Marshal(req.Tags)
	properties, _ := json.Marshal(req.Properties)
	targeting, _ := json.Marshal(req.Targeting)
	capping, _ := json.Marshal(req.Capping)
	hooks, _ := json.Marshal(req.Hooks)
	testUsers, _ := json.Marshal(req.TestUsers)

	if len(req.Targeting.Audiences) > 0 {
		if err := ac.App.DB.ValidateAudience(
			c.Request.Context(),
			user.OrganizationID,
			req.Targeting.Audiences,
		); err != nil {
			ac.App.Logger.Warn("actions.create.validateAudience", zap.Error(err))
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
	}

	args := db.CreateActionParams{
		OrganizationID: user.OrganizationID,
		Type:           req.Type,
		Name:           req.Name,
		State:          db.ActionState(req.State),
		Tags:           tags,
		Properties:     properties,
		Targeting:      targeting,
		Capping:        capping,
		Hooks:          hooks,
		TestUsers:      testUsers,
	}
	action, err := ac.App.DB.CreateAction(c.Request.Context(), args)
	if err != nil {
		ac.App.Logger.Error("actions.create.createAction", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	res := responses.Action{
		ID:             action.ID,
		OrganizationID: action.OrganizationID,
		Name:           action.Name,
		Type:           action.Type,
		Properties:     action.Properties,
		Targeting:      action.Targeting,
		Capping:        action.Capping,
		Hooks:          action.Hooks,
		State:          string(action.State),
		CreatedAt:      action.CreatedAt,
		UpdatedAt:      action.UpdatedAt,
	}

	c.JSON(http.StatusOK, res)
}

func (ac ActionsController) Duplicate(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	id := uuid.FromStringOrNil(c.Param("id"))
	action, err := ac.App.DB.GetAction(c.Request.Context(), user.OrganizationID, id)
	if err != nil {
		ac.App.Logger.Warn("actions.duplicate.getAction", zap.Error(err))
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if action.OrganizationID != user.OrganizationID {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	args := db.CreateActionParams{
		OrganizationID: user.OrganizationID,
		Type:           action.Type,
		Name:           fmt.Sprintf("%s (copy)", action.Name),
		State:          db.ActionState(db.StateInactive),
		Tags:           action.Tags,
		Properties:     action.Properties,
		Targeting:      action.Targeting,
		Capping:        action.Capping,
		Hooks:          action.Hooks,
		TestUsers:      action.TestUsers,
	}
	if _, err = ac.App.DB.CreateAction(c.Request.Context(), args); err != nil {
		ac.App.Logger.Error("actions.duplicate.createAction", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

func (ac ActionsController) Update(c *gin.Context) {
	var req requests.ActionCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		ac.App.Logger.Warn("actions.update.parseRequest", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user := c.MustGet("user").(*db.User)
	tags, _ := json.Marshal(req.Tags)
	properties, _ := json.Marshal(req.Properties)
	targeting, _ := json.Marshal(req.Targeting)
	capping, _ := json.Marshal(req.Capping)
	hooks, _ := json.Marshal(req.Hooks)
	testUsers, _ := json.Marshal(req.TestUsers)

	if len(req.Targeting.Audiences) > 0 {
		if err := ac.App.DB.ValidateAudience(
			c.Request.Context(),
			user.OrganizationID,
			req.Targeting.Audiences,
		); err != nil {
			ac.App.Logger.Warn("actions.create.validateAudience", zap.Error(err))
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
	}

	args := db.UpdateActionParams{
		OrganizationID: user.OrganizationID,
		Type:           req.Type,
		Name:           req.Name,
		Tags:           tags,
		State:          db.ActionState(req.State),
		Properties:     properties,
		Targeting:      targeting,
		Capping:        capping,
		Hooks:          hooks,
		TestUsers:      testUsers,
		ID:             uuid.FromStringOrNil(c.Param("id")),
	}
	if err := ac.App.DB.UpdateAction(c.Request.Context(), args); err != nil {
		ac.App.Logger.Error("actions.update.updateAction", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

func (ac ActionsController) UpdateState(c *gin.Context) {
	var req requests.ActionState
	if err := c.ShouldBindJSON(&req); err != nil {
		ac.App.Logger.Warn("actions.updateState.parseRequest", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user := c.MustGet("user").(*db.User)
	id := uuid.FromStringOrNil(c.Param("id"))
	if err := ac.App.DB.UpdateActionState(c.Request.Context(),
		user.OrganizationID,
		id,
		db.ActionState(req.State),
	); err != nil {
		ac.App.Logger.Error("actions.update.updateActionState", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

func (ac ActionsController) Delete(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	id := uuid.FromStringOrNil(c.Param("id"))
	if err := ac.App.DB.DeleteAction(c.Request.Context(), user.OrganizationID, id); err != nil {
		ac.App.Logger.Error("actions.delete.deleteAction", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

func (ac ActionsController) TestWebhook(c *gin.Context) {
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
