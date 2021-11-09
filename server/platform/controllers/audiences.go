package platform

import (
	"net/http"
	"time"

	"attractify.io/platform/app"
	"attractify.io/platform/audiences"
	"attractify.io/platform/db"
	"attractify.io/platform/platform/requests"
	"attractify.io/platform/platform/responses"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
)

type AudienceController struct {
	Router *gin.RouterGroup
	App    *app.App
}

func InitAudiences(router *gin.RouterGroup, app *app.App) {
	a := AudienceController{Router: router, App: app}
	a.Router.GET("/audiences", a.List)
	a.Router.GET("/audiences/:id", a.Show)
	a.Router.POST("/audiences", a.Create)
	a.Router.DELETE("/audiences/:id", a.Delete)
	a.Router.PUT("/audiences/:id", a.Update)
	a.Router.POST("/audiences/preview", a.Preview)
	a.Router.PUT("/audiences/:id/refresh", a.Refresh)
}

func (ac AudienceController) List(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	audiences, err := ac.App.DB.GetAudiences(c.Request.Context(), user.OrganizationID)
	if err != nil {
		ac.App.Logger.Error("audiences.list.getAudiences", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	res := []responses.Audience{}
	for _, a := range audiences {
		// FIXME: Extract to standalone methode.
		res = append(res, responses.Audience{
			ID:               a.ID,
			Name:             a.Name,
			Description:      a.Description,
			IncludeAnonymous: a.IncludeAnonymous,
			Events:           a.Events,
			Traits:           a.Traits,
			ProfileCount:     a.ProfileCount,
			CreatedAt:        a.CreatedAt,
			UpdatedAt:        a.UpdatedAt,
			RefreshedAt:      a.RefreshedAt,
		})
	}

	c.JSON(http.StatusOK, res)
}

func (ac AudienceController) Show(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	id := uuid.FromStringOrNil(c.Param("id"))
	audience, err := ac.App.DB.GetAudience(c.Request.Context(), user.OrganizationID, id)
	if err != nil {
		ac.App.Logger.Warn("audiences.show.getAudience", zap.Error(err))
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	// FIXME: Extract to standalone methode.
	res := responses.Audience{
		ID:               audience.ID,
		Name:             audience.Name,
		Description:      audience.Description,
		IncludeAnonymous: audience.IncludeAnonymous,
		Events:           audience.Events,
		Traits:           audience.Traits,
		ProfileCount:     audience.ProfileCount,
		CreatedAt:        audience.CreatedAt,
		UpdatedAt:        audience.UpdatedAt,
		RefreshedAt:      audience.RefreshedAt,
	}

	c.JSON(http.StatusOK, res)
}

func (ac AudienceController) Create(c *gin.Context) {
	var req requests.AudienceCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		ac.App.Logger.Warn("audiences.create.parseRequest", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user := c.MustGet("user").(*db.User)

	args := db.CreateAudienceParams{
		OrganizationID:   user.OrganizationID,
		Name:             req.Name,
		Description:      req.Description,
		IncludeAnonymous: req.IncludeAnonymous,
		Events:           req.Events,
		Traits:           req.Traits,
	}
	audience, err := ac.App.DB.CreateAudience(c.Request.Context(), args)
	if err != nil {
		ac.App.Logger.Error("audiences.create.createAudience", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// FIXME: Extract to standalone methode.
	res := responses.Audience{
		ID:               audience.ID,
		Name:             audience.Name,
		Description:      audience.Description,
		IncludeAnonymous: audience.IncludeAnonymous,
		Events:           audience.Events,
		Traits:           audience.Traits,
		CreatedAt:        audience.CreatedAt,
		UpdatedAt:        audience.UpdatedAt,
		RefreshedAt:      audience.RefreshedAt,
	}

	c.JSON(http.StatusOK, res)
}

func (ac AudienceController) Update(c *gin.Context) {
	var req requests.AudienceCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		ac.App.Logger.Warn("audiences.update.parseRequest", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user := c.MustGet("user").(*db.User)

	args := db.UpdateAudienceParams{
		OrganizationID:   user.OrganizationID,
		Name:             req.Name,
		Description:      req.Description,
		IncludeAnonymous: req.IncludeAnonymous,
		Events:           req.Events,
		Traits:           req.Traits,
		ID:               uuid.FromStringOrNil(c.Param("id")),
	}

	if err := ac.App.DB.UpdateAudience(c.Request.Context(), args); err != nil {
		ac.App.Logger.Error("audiences.update.updateAudience", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

func (ac AudienceController) Delete(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	id := uuid.FromStringOrNil(c.Param("id"))
	if err := ac.App.DB.DeleteAudience(c.Request.Context(), user.OrganizationID, id); err != nil {
		ac.App.Logger.Warn("audiences.delete.deleteAudience", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

func (ac AudienceController) Preview(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	var req requests.AudienceCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		ac.App.Logger.Warn("audiences.preview.parseRequest", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	start := time.Now()

	audience := db.Audience{
		OrganizationID:   user.OrganizationID,
		IncludeAnonymous: req.IncludeAnonymous,
		Traits:           req.Traits,
		Events:           req.Events,
	}
	a := audiences.New(c.Request.Context(), ac.App, &audience)
	profiles, err := a.Preview()
	if err != nil {
		ac.App.Logger.Error("audiences.refresh.preview", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	end := time.Now()

	pl := []responses.AudienceProfile{}
	for _, p := range profiles {
		pl = append(pl, responses.AudienceProfile{
			ID: p.ID,
		})
	}
	res := responses.AudiencePreview{
		Profiles:      pl,
		ExecutionTime: end.Sub(start).Milliseconds(),
	}

	c.JSON(http.StatusOK, res)
}

func (ac AudienceController) Refresh(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	id := uuid.FromStringOrNil(c.Param("id"))
	audience, err := ac.App.DB.GetAudience(c.Request.Context(), user.OrganizationID, id)
	if err != nil {
		ac.App.Logger.Warn("audiences.refresh.getAudience", zap.Error(err))
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	a := audiences.New(c.Request.Context(), ac.App, &audience)
	count, err := a.Refresh()
	if err != nil {
		ac.App.Logger.Error("audiences.refresh.refresh", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	res := responses.AudienceRefresh{Count: count}

	c.JSON(http.StatusOK, res)
}
