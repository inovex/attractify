package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"attractify.io/platform/analytics"
	"attractify.io/platform/app"
	"attractify.io/platform/audiences"
	"attractify.io/platform/config"
	"attractify.io/platform/db"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	_ "github.com/lib/pq"
)

const (
	defaultConfig = "/run/secrets/config.json"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	app := app.App{Logger: logger}

	cfgPath := defaultConfig
	if len(os.Args) > 1 {
		cfgPath = os.Args[1]
	}
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

	ctx, cancel := context.WithCancel(context.Background())

	log.Println("Starting audience update...")

	a, err := app.DB.GetAllAudiences(ctx)
	if err != nil {
		panic(err)
	}

	for _, audience := range a {
		a := audiences.New(ctx, &app, &audience)
		if _, err := a.Refresh(); err != nil {
			panic(err)
		}
	}

	fmt.Println("Audience update completed...")
	app.DB.Close()
	app.Analytics.Close()
	cancel()
}
