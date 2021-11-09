package middlewares

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"attractify.io/platform/app"
	"github.com/gin-gonic/gin"
	uuid "github.com/gofrs/uuid"
	"gopkg.in/square/go-jose.v2/jwt"
)

const (
	platformAudience      = "platform"
	activationAudience    = "activation"
	resetPasswordAudience = "reset_password"
)

func APITokenAuth(app *app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if len(auth) == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		token, err := parseAPIToken(auth)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		organization, err := app.DB.GetOrganizationByToken(c.Request.Context(), token)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("auth", &organization)
	}
}

func PlatformTokenAuth(app *app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "POST" &&
			(c.FullPath() == "/v1/platform/user/session" ||
				c.FullPath() == "/v1/platform/organization" ||
				c.FullPath() == "/v1/platform/user/reset-password") {
			return
		}
		auth := c.GetHeader("Authorization")
		if len(auth) == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims, err := parseJWTToken(auth, app.Config.AuthKey)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if c.Request.Method == "POST" && c.FullPath() == "/v1/platform/user" {
			if !claims.Audience.Contains(activationAudience) {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			if claims.IssuedAt.Time().Add(time.Hour * 24).Before(time.Now().UTC()) {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
		} else if c.Request.Method == "PUT" && c.FullPath() == "/v1/platform/user/reset-password" {
			if !claims.Audience.Contains(resetPasswordAudience) {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			if claims.IssuedAt.Time().Add(time.Minute * 15).Before(time.Now().UTC()) {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
		} else {
			if !claims.Audience.Contains(platformAudience) {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
		}

		id := uuid.FromStringOrNil(claims.Subject)
		user, err := app.DB.GetUser(c.Request.Context(), id)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if user.LoggedOutAt != nil && claims.IssuedAt.Time().Before(*user.LoggedOutAt) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("user", &user)
	}
}

func parseJWTToken(header string, key []byte) (*jwt.Claims, error) {
	token := strings.Split(header, "Bearer ")
	if len(token) != 2 {
		return nil, errTokenNotFound
	}

	tok, err := jwt.ParseEncrypted(token[1])
	if err != nil {
		return nil, err
	}

	claims := jwt.Claims{}
	if err := tok.Claims(key, &claims); err != nil {
		return nil, err
	}

	return &claims, nil
}

func parseAPIToken(header string) (string, error) {
	token := strings.Split(header, "Bearer ")
	if len(token) != 2 {
		return "", errTokenNotFound
	}

	return token[1], nil
}

var (
	errTokenNotFound = errors.New("could not find token in header")
)
