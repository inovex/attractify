package main

import (
	"net/http"
	"os"
	"strings"

	"attractify.io/platform/analytics"
	api "attractify.io/platform/api/controllers"
	"attractify.io/platform/app"
	"attractify.io/platform/config"
	"attractify.io/platform/db"
	"attractify.io/platform/mailer"
	"attractify.io/platform/middlewares"
	platform "attractify.io/platform/platform/controllers"
	"attractify.io/platform/stream"
	"go.uber.org/zap"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

var corsConfig = cors.Config{
	AllowOrigins:     []string{"*"},
	AllowMethods:     []string{"*"},
	AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
	ExposeHeaders:    []string{"Content-Length"},
	AllowCredentials: true,
}

func apiHandler(app *app.App) http.Handler {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(cors.New(corsConfig))
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	p := r.Group("/v1/platform")
	p.Use(middlewares.PlatformTokenAuth(app))
	{
		platform.InitDashboard(p, app)
		platform.InitOrganization(p, app)
		platform.InitUser(p, app)
		platform.InitUsers(p, app)
		platform.InitActions(p, app)
		platform.InitReactions(p, app)
		platform.InitPrivacy(p, app)
		platform.InitAnalyze(p, app)
		platform.InitEvents(p, app)
		platform.InitData(p, app)
		platform.InitEventLog(p, app)
		platform.InitContexts(p, app)
		platform.InitAudiences(p, app)
		platform.InitProfiles(p, app)
		platform.InitChannels(p, app)
		platform.InitAuthTokens(p, app)
		platform.InitCustomTraits(p, app)
		platform.InitComputedTraits(p, app)
	}

	t := r.Group("/v1")
	t.Use(middlewares.APITokenAuth(app))
	{
		api.InitIdentify(t, app)
		api.InitTrack(t, app)
		api.InitActions(t, app)
	}

	return r
}

func frontendHandler(app *app.App) http.Handler {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	r.Static("/", "./dist")
	r.NoRoute(func(c *gin.Context) {
		c.File("./dist/index.html")
	})

	return r
}

type Gateway struct {
	api      http.Handler
	frontend http.Handler
}

func NewGateway(app *app.App) *Gateway {
	gw := Gateway{}
	gw.api = apiHandler(app)
	gw.frontend = frontendHandler(app)

	return &gw
}

func (g Gateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/v1") {
		g.api.ServeHTTP(w, r)
	} else {
		g.frontend.ServeHTTP(w, r)
	}
}

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	app := app.App{Logger: logger}

	if len(os.Args) <= 1 {
		panic("missing config")
	}

	cfgPath := os.Args[1]
	app.Config, err = config.Parse(cfgPath)
	if err != nil {
		panic(err)
	}

	dbConn, err := sqlx.Open("postgres", app.Config.DB)
	if err != nil {
		panic(err)
	}
	app.DB = db.New(dbConn)

	app.Analytics, err = analytics.OpenDB(app.Config.Analytics)
	if err != nil {
		panic(err)
	}

	app.Stream = stream.New(app.Config.Stream.Brokers, app.Config.Stream.Topic)

	if !app.Config.Server.DebugMode {
		gin.SetMode(gin.ReleaseMode)
	}

	app.Mailer = mailer.New(app.Config)

	gw := NewGateway(&app)
	server := &http.Server{
		Addr:    app.Config.Server.Bind,
		Handler: gw,
	}

	if len(app.Config.Server.Cert) > 0 {
		err = server.ListenAndServeTLS(app.Config.Server.Cert, app.Config.Server.Key)
	} else {
		err = server.ListenAndServe()
	}
	if err != nil {
		panic(err)
	}
}
