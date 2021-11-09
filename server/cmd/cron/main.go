package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"attractify.io/platform/analytics"
	"attractify.io/platform/app"
	"attractify.io/platform/config"
	"attractify.io/platform/cron"
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

	cr := cron.New(ctx, &app)
	cr.Init()

	log.Println("Cron deamon running...")

	wg := sync.WaitGroup{}
	wg.Add(1)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		fmt.Println("Shutting cron deamon down")
		app.DB.Close()
		app.Analytics.Close()
		cancel()
		wg.Done()
	}()

	wg.Wait()
}
