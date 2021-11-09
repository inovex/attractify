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

type ComputedTraitsController struct {
	Router *gin.RouterGroup
	App    *app.App
}

func InitComputedTraits(router *gin.RouterGroup, app *app.App) {
	c := ComputedTraitsController{Router: router, App: app}
	c.Router.GET("/computed-traits", c.List)
	c.Router.GET("/computed-traits/:id", c.Show)
	c.Router.POST("/computed-traits", c.Create)
	c.Router.DELETE("/computed-traits/:id", c.Delete)
	c.Router.PUT("/computed-traits/:id", c.Update)
	c.Router.POST("/computed-traits/:id/refresh", c.Refresh)
}

func (cc ComputedTraitsController) List(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	computedTraits, err := cc.App.DB.GetComputedTraits(c.Request.Context(), user.OrganizationID)
	if err != nil {
		cc.App.Logger.Error("computedTraits.list.getComputedTraits", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	res := []responses.ComputedTrait{}
	for _, c := range computedTraits {
		// FIXME: Extract to standalone methode.
		res = append(res, responses.ComputedTrait{
			ID:          c.ID,
			Name:        c.Name,
			Key:         c.Key,
			Type:        string(c.Type),
			EventID:     c.EventID,
			Conditions:  c.Conditions,
			Properties:  c.Properties,
			CreatedAt:   c.CreatedAt,
			UpdatedAt:   c.UpdatedAt,
			RefreshedAt: c.RefreshedAt,
		})
	}

	c.JSON(http.StatusOK, res)
}

func (cc ComputedTraitsController) Show(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	id := uuid.FromStringOrNil(c.Param("id"))
	computedTrait, err := cc.App.DB.GetComputedTrait(c.Request.Context(), user.OrganizationID, id)
	if err != nil {
		cc.App.Logger.Warn("computedTraits.show.getComputedTrait", zap.Error(err))
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	// FIXME: Extract to standalone methode.
	res := responses.ComputedTrait{
		ID:          computedTrait.ID,
		Name:        computedTrait.Name,
		Key:         computedTrait.Key,
		Type:        string(computedTrait.Type),
		EventID:     computedTrait.EventID,
		Conditions:  computedTrait.Conditions,
		Properties:  computedTrait.Properties,
		CreatedAt:   computedTrait.CreatedAt,
		UpdatedAt:   computedTrait.UpdatedAt,
		RefreshedAt: computedTrait.RefreshedAt,
	}

	c.JSON(http.StatusOK, res)
}

func (cc ComputedTraitsController) Create(c *gin.Context) {
	var req requests.ComputedTraitCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		cc.App.Logger.Warn("computedTraits.create.parseRequest", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user := c.MustGet("user").(*db.User)

	args := db.CreateComputedTraitParams{
		OrganizationID: user.OrganizationID,
		Name:           req.Name,
		Key:            req.Key,
		Type:           req.Type,
		EventID:        req.EventID,
		Conditions:     req.Conditions,
		Properties:     req.Properties,
	}
	computedTrait, err := cc.App.DB.CreateComputedTrait(c.Request.Context(), args)
	if err != nil {
		cc.App.Logger.Error("computedTraits.create.createComputedTrait", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// FIXME: Extract to standalone methode.
	res := responses.ComputedTrait{
		ID:          computedTrait.ID,
		Name:        computedTrait.Name,
		Key:         computedTrait.Key,
		Type:        string(computedTrait.Type),
		EventID:     computedTrait.EventID,
		Conditions:  computedTrait.Conditions,
		Properties:  computedTrait.Properties,
		CreatedAt:   computedTrait.CreatedAt,
		UpdatedAt:   computedTrait.UpdatedAt,
		RefreshedAt: computedTrait.RefreshedAt,
	}

	c.JSON(http.StatusOK, res)
}

func (cc ComputedTraitsController) Update(c *gin.Context) {
	var req requests.ComputedTraitCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		cc.App.Logger.Warn("computedTraits.update.parseRequest", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user := c.MustGet("user").(*db.User)

	args := db.UpdateComputedTraitParams{
		OrganizationID: user.OrganizationID,
		Name:           req.Name,
		Key:            req.Key,
		EventID:        req.EventID,
		Conditions:     req.Conditions,
		Properties:     req.Properties,
		ID:             uuid.FromStringOrNil(c.Param("id")),
	}

	if err := cc.App.DB.UpdateComputedTrait(c.Request.Context(), args); err != nil {
		cc.App.Logger.Error("computedTraits.update.updateComputedTrait", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

func (cc ComputedTraitsController) Delete(c *gin.Context) {
	user := c.MustGet("user").(*db.User)

	id := uuid.FromStringOrNil(c.Param("id"))
	if err := cc.App.DB.DeleteComputedTrait(c.Request.Context(), user.OrganizationID, id); err != nil {
		cc.App.Logger.Warn("computedTraits.delete.deleteComputedTrait", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

func (cc ComputedTraitsController) Refresh(c *gin.Context) {
	// user := c.MustGet("user").(*db.User)

	// args := db.GetComputedTraitParams{
	// 	OrganizationID: user.OrganizationID,
	// 	ID:             uuid.FromStringOrNil(c.Param("id")),
	// }
	// computedTrait, err := cc.App.DB.GetComputedTrait(c.Request.Context(), args)
	// if err != nil {
	// 	cc.App.Logger.Warn("computedTraits.refresh.getComputedTrait", zap.Error(err))
	// 	c.AbortWithStatus(http.StatusNotFound)
	// 	return
	// }

	// a := computedTraits.NewRefresher(c.Request.Context(), cc.App, &audience)
	// setID, count, err := a.Refresh()
	// if err != nil {
	// 	cc.App.Logger.Error("computedTraits.refresh.refresh", zap.Error(err))
	// 	c.AbortWithStatus(http.StatusBadRequest)
	// 	return
	// }

	// params := db.UpdateComputedTraitProfilesParams{
	// 	CurrentSetID:   setID,
	// 	ProfileCount:   count,
	// 	OrganizationID: audience.OrganizationID,
	// 	ID:             audience.ID,
	// }
	// if err := cc.App.DB.UpdateComputedTraitProfiles(c.Request.Context(), params); err != nil {
	// 	cc.App.Logger.Error("computedTraits.refresh.updateComputedTraitProfiles", zap.Error(err))
	// 	c.AbortWithStatus(http.StatusBadRequest)
	// 	return
	// }

	// if audience.CurrentSetID.Valid {
	// 	dParams := db.DeleteComputedTraitProfilesBySetIDParams{
	// 		OrganizationID: audience.OrganizationID,
	// 		SetID:          audience.CurrentSetID.UUID,
	// 	}
	// 	if err := cc.App.DB.DeleteComputedTraitProfilesBySetID(c.Request.Context(), dParams); err != nil {
	// 		cc.App.Logger.Error("computedTraits.refresh.deleteComputedTraitProfilesBySetID", zap.Error(err))
	// 		c.AbortWithStatus(http.StatusBadRequest)
	// 		return
	// 	}
	// }

	// res := responses.ComputedTraitRefresh{Count: count}

	// c.JSON(http.StatusOK, res)
}
