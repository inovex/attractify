package platform

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"attractify.io/platform/analytics"
	"attractify.io/platform/app"
	"attractify.io/platform/computedtraits"
	"attractify.io/platform/db"
	"attractify.io/platform/platform/responses"
	"attractify.io/platform/privacy"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
)

type ProfilesController struct {
	Router *gin.RouterGroup
	App    *app.App
}

func InitProfiles(router *gin.RouterGroup, app *app.App) {
	p := ProfilesController{Router: router, App: app}
	p.Router.GET("/profiles", p.List)
	p.Router.GET("/profiles/:id", p.Show)
	p.Router.GET("/profiles/searchbyuserid/:id/traits", p.SearchByUserIDWithTraits)
	p.Router.GET("/profiles/searchbyuserid/:id", p.SearchByUserID)
	p.Router.DELETE("/profiles/:id", p.Delete)
	p.Router.DELETE("/profiles/:id/identity", p.DeleteIdentity)
	p.Router.GET("/profiles/:id/identities", p.ListIdentities)
	p.Router.GET("/profiles/:id/events", p.ListEvents)
	p.Router.POST("/profiles/:id/refresh-computed-traits", p.RefreshComputedTraits)
}

func (pc ProfilesController) List(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	profiles, err := pc.App.DB.GetProfilesForOrganization(c.Request.Context(), user.OrganizationID, 10, 0)
	if err != nil {
		pc.App.Logger.Error("profiles.list.getProfileForOrganization", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	res := []responses.Profile{}
	for _, p := range profiles {
		res = append(res, responses.Profile{
			ID:             p.ID,
			CustomTraits:   p.CustomTraits,
			ComputedTraits: p.ComputedTraits,
			CreatedAt:      p.CreatedAt,
			UpdatedAt:      p.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, res)
}

func (pc ProfilesController) Show(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	id := uuid.FromStringOrNil(c.Param("id"))
	p, err := pc.App.DB.GetProfile(c.Request.Context(), user.OrganizationID, id)
	if err != nil {
		pc.App.Logger.Warn("profile.show.getProfileForOrganization", zap.Error(err))
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	res := responses.Profile{
		ID:             p.ID,
		CustomTraits:   p.CustomTraits,
		ComputedTraits: p.ComputedTraits,
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
	}

	c.JSON(http.StatusOK, res)
}

func (pc ProfilesController) SearchByUserID(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	id := escape(c.Param("id"))
	profiles, err := pc.App.DB.SearchByUserID(c.Request.Context(), user.OrganizationID, id)
	if err != nil {
		pc.App.Logger.Warn("profile.search.searchForUserID", zap.Error(err))
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	res := []responses.Profile{}
	for _, p := range profiles {
		res = append(res, responses.Profile{
			ID:             p.ProfileID,
			CustomTraits:   p.CustomTraits,
			ComputedTraits: p.ComputedTraits,
			CreatedAt:      p.CreatedAt,
			UpdatedAt:      p.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, res)
}

func (pc ProfilesController) SearchByUserIDWithTraits(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	id := escape(c.Param("id"))
	profiles, err := pc.App.DB.SearchByUserID(c.Request.Context(), user.OrganizationID, id)
	if err != nil {
		pc.App.Logger.Warn("profile.search.searchForUserID", zap.Error(err))
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	res := []responses.IdentityWithTraits{}
	for _, p := range profiles {
		res = append(res, responses.IdentityWithTraits{
			ID:             p.ID,
			Channel:        p.Channel,
			Type:           p.Type,
			UserID:         p.UserID,
			CustomTraits:   p.CustomTraits,
			ComputedTraits: p.ComputedTraits,
			IsAnonymous:    p.IsAnonymous,
			CreatedAt:      p.CreatedAt,
			UpdatedAt:      p.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, res)
}

func (pc ProfilesController) Delete(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	p := privacy.NewDeletionByProfileID(pc.App, user.OrganizationID, uuid.FromStringOrNil(c.Param("id")))
	if err := p.Run(); err != nil {
		pc.App.Logger.Warn("profiles.delete.run", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

func (pc ProfilesController) DeleteIdentity(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	p := privacy.NewDeletionByUserID(pc.App, user.OrganizationID, c.Param("id"))
	if err := p.DeleteSingleIdentity(); err != nil {
		pc.App.Logger.Warn("profiles.delete.DeleteSingleIdentity", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

func (pc ProfilesController) ListIdentities(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	id := uuid.FromStringOrNil(c.Param("id"))
	identities, err := pc.App.DB.GetProfileIdentitiesForProfile(c.Request.Context(), user.OrganizationID, id)
	if err != nil {
		pc.App.Logger.Warn("profile.listIdentities.getProfileIdentitiesForProfile", zap.Error(err))
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	var res []responses.Identity
	for _, i := range identities {
		idt := responses.Identity{
			ID:          i.ID,
			Channel:     i.Channel,
			Type:        i.Type,
			UserID:      i.UserID,
			IsAnonymous: i.IsAnonymous,
			CreatedAt:   i.CreatedAt,
			UpdatedAt:   i.UpdatedAt,
		}
		res = append(res, idt)
	}

	c.JSON(http.StatusOK, res)
}

func (pc ProfilesController) ListEvents(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	id := uuid.FromStringOrNil(c.Param("id"))
	identities, err := pc.App.DB.GetProfileIdentitiesForProfile(c.Request.Context(), user.OrganizationID, id)
	if err != nil {
		pc.App.Logger.Warn("profile.listEvents.getProfileIdentitiesForProfile", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	identityIDs := []uuid.UUID{}
	for _, i := range identities {
		identityIDs = append(identityIDs, i.ID)
	}

	start, _ := time.Parse("2006-01-02 15:04:05", "2000-01-01 00:00:00")
	end, _ := time.Parse("2006-01-02 15:04:05", "2099-01-01 00:00:00")
	eArgs := analytics.GetEventsParams{
		OrganizationID: user.OrganizationID,
		IdentityIDs:    identityIDs,
		Start:          start,
		End:            end,
		Offset:         0,
		Limit:          10,
	}
	events, err := pc.App.Analytics.GetEvents(eArgs)
	if err != nil {
		pc.App.Logger.Warn("profile.listEvents.getEventsForProfile", zap.Error(err))
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	res := []responses.EventLog{}
	for _, e := range events {
		evt := responses.EventLog{
			ID:         e.ID,
			EventID:    e.EventID,
			Channel:    e.Channel,
			Context:    json.RawMessage(e.Context),
			Properties: json.RawMessage(e.Properties),
			CreatedAt:  e.CreatedAt,
		}
		res = append(res, evt)
	}

	c.JSON(http.StatusOK, res)
}

func (pc ProfilesController) RefreshComputedTraits(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	id := uuid.FromStringOrNil(c.Param("id"))
	p, err := pc.App.DB.GetProfile(c.Request.Context(), user.OrganizationID, id)
	if err != nil {
		pc.App.Logger.Warn("profile.refreshComputedTraits.getProfileForOrganization", zap.Error(err))
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	computedTraits, err := pc.App.DB.GetComputedTraits(c.Request.Context(), user.OrganizationID)
	if err != nil {
		pc.App.Logger.Warn("profile.refreshComputedTraits.getComputedTraits", zap.Error(err))
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	for _, ct := range computedTraits {
		ctr := computedtraits.New(c.Request.Context(), pc.App, &ct)
		if err := ctr.Refresh(&p); err != nil {
			pc.App.Logger.Warn("profile.refreshComputedTraits.refresh", zap.Error(err))
		}
	}

	c.AbortWithStatus(http.StatusNoContent)
}

func escape(v string) string {
	return strings.NewReplacer(`\`, `\\`, `'`, `\'`).Replace(v)
}
