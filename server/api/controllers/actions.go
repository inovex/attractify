package api

import (
	"encoding/json"
	"net/http"
	"time"

	"attractify.io/platform/actions"
	"attractify.io/platform/analytics"
	"attractify.io/platform/api/requests"
	"attractify.io/platform/api/responses"
	"attractify.io/platform/app"
	"attractify.io/platform/db"
	"attractify.io/platform/privacy"
	"attractify.io/platform/profiles"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"

	"go.uber.org/zap"
)

type ActionsController struct {
	Router *gin.RouterGroup
	App    *app.App
}

func InitActions(router *gin.RouterGroup, app *app.App) {
	c := ActionsController{Router: router, App: app}
	c.Router.GET("/actions", c.List)
	c.Router.POST("/actions", c.ListWithCtx)
	c.Router.POST("/actions/act", c.Act)
}

func (ac ActionsController) List(c *gin.Context) {
	var req requests.Actions
	if err := c.ShouldBindQuery(&req); err != nil {
		ac.App.Logger.Warn("api.actions.list.parseRequest", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ac.list(c, req)
}

func (ac ActionsController) ListWithCtx(c *gin.Context) {
	var req requests.Actions
	if err := c.ShouldBindJSON(&req); err != nil {
		ac.App.Logger.Warn("api.actions.list.parseRequest", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ac.list(c, req)
}

func (ac ActionsController) list(c *gin.Context, req requests.Actions) {
	auth := c.MustGet("auth").(*db.OrganizationAuth)
	res := []responses.Action{}

	// Is user locked
	l := privacy.NewLocked(ac.App, auth.OrganizationID, req.UserID)
	if l.IsLocked() {
		c.JSON(http.StatusOK, res)
		return
	}

	// Retrieve profile or create new one.
	params := profiles.Params{
		Time:           time.Now().UTC(),
		OrganizationID: auth.OrganizationID,
		UserID:         req.UserID,
		IsAnonymous:    true,
		Channel:        auth.Channel,
		Type:           "anonymous_id",
	}
	p := profiles.New(c.Request.Context(), ac.App, params)
	profile, identity, err := p.GetOrCreate()
	if err != nil {
		ac.App.Logger.Error("api.actions.list.profile", zap.Error(err))
		c.JSON(http.StatusBadRequest, res)
		return
	}

	// Get all actions
	acts, err := ac.App.DB.GetActions(c.Request.Context(), auth.OrganizationID)
	if err != nil {
		ac.App.Logger.Error("api.actions.list.getActions", zap.Error(err))
		c.JSON(http.StatusOK, res)
		return
	}

	now := time.Now()
	var caArgs []analytics.CreateReactionParams
	for _, ca := range acts {
		a := actions.New(c.Request.Context(), ac.App, auth.OrganizationID, &ca, profile, identity)
		if a.ShouldDisplay(req.TypeName, req.Tags, auth.Channel, req.UserID, req.Context, time.Now().UTC(), auth.Timezone) {
			res = append(res, responses.Action{
				ID:         ca.ID,
				TypeName:   ca.TypeName,
				Version:    ca.TypeVersion,
				Tags:       ca.Tags,
				Properties: a.MapProperties(auth.Channel),
			})

			caArgs = append(caArgs, analytics.CreateReactionParams{
				OrganizationID: auth.OrganizationID,
				ActionID:       ca.ID,
				IdentityID:     identity.ID,
				Channel:        auth.Channel,
				Event:          analytics.ReactionEventDelivered,
				Context:        string(req.Context),
				CreatedAt:      now,
			})
		}
	}

	if err := ac.App.Analytics.CreateReactions(caArgs); err != nil {
		ac.App.Logger.Error("api.actions.list.createReactions", zap.Error(err))
	}

	c.JSON(http.StatusOK, res)
}

func (ac ActionsController) Act(c *gin.Context) {
	auth := c.MustGet("auth").(*db.OrganizationAuth)

	var req requests.ActionsAct
	if err := c.ShouldBindJSON(&req); err != nil {
		ac.App.Logger.Warn("api.actions.act.parseRequest", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// Is user locked
	l := privacy.NewLocked(ac.App, auth.OrganizationID, req.UserID)
	if l.IsLocked() {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	pIdentity, err := ac.App.DB.GetProfileIdentityForUserID(
		c.Request.Context(), auth.OrganizationID, req.UserID,
	)
	if err != nil {
		ac.App.Logger.Warn("api.actions.act.getProfileIdentitiyForUserID", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	profile, err := ac.App.DB.GetProfile(c.Request.Context(), auth.OrganizationID, pIdentity.ProfileID)
	if err != nil {
		ac.App.Logger.Warn("api.actions.act.getProfile", zap.Error(err))
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	actionID := uuid.FromStringOrNil(req.ActionID)
	action, err := ac.App.DB.GetAction(c.Request.Context(), auth.OrganizationID, actionID)
	if err != nil {
		ac.App.Logger.Warn("api.actions.act.getAction", zap.Error(err))
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	var event string
	switch req.Event {
	case "show":
		event = "shown"
	case "hide":
		event = "hidden"
	case "decline":
		event = "declined"
	case "accept":
		event = "accepted"
	}

	if len(event) == 0 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// Make sure the user is allowed to accept the event.
	a := actions.New(c.Request.Context(), ac.App, auth.OrganizationID, &action, &profile, &pIdentity)
	if event == analytics.ReactionEventAccepted {
		err = a.IsAllowedToAccept(auth.Channel, req.UserID, time.Now().UTC(), auth.Timezone)
		if err != nil {
			ac.App.Logger.Warn("api.actions.act.canAccept " + err.Error())
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
	}

	res, isHookSuccessful, err := a.RunHooks(req.UserID, event, auth.Channel, req.Context, req.Properties)
	if err != nil {
		ac.App.Logger.Error("api.actions.act.executeWebhook", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if !isHookSuccessful {
		var hookResult actions.Result
		json.Unmarshal(res, &hookResult)
		c.JSON(hookResult.StatusCode, res) //TODO: test with bad webhook
		return
	}

	caArgs := analytics.CreateReactionParams{
		OrganizationID: auth.OrganizationID,
		ActionID:       actionID,
		IdentityID:     pIdentity.ID,
		Channel:        auth.Channel,
		Event:          event,
		CreatedAt:      time.Now(),
	}
	if req.Context != nil {
		caArgs.Context = string(*req.Context)
	}

	// TODO: Find better way to enable/disable storing of props/results
	// if req.Properties != nil {
	// 	caArgs.Properties = string(*req.Properties)
	// }

	// if res != nil {
	// 	caArgs.Result = string(res)
	// }

	if err := ac.App.Analytics.CreateReaction(caArgs); err != nil {
		ac.App.Logger.Error("api.actions.act.createReactions", zap.Error(err))
	}

	if res == nil {
		c.JSON(http.StatusOK, gin.H{})
		return
	}

	c.JSON(http.StatusOK, res)
}
