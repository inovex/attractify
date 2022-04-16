package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"attractify.io/platform/analytics"
	"attractify.io/platform/app"
	"attractify.io/platform/config"
	"attractify.io/platform/consumers"
	"attractify.io/platform/db"
	"attractify.io/platform/stream"
	"go.uber.org/zap"

	"github.com/jmoiron/sqlx"

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

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	wg := sync.WaitGroup{}

	go func() {
		<-sigs
		log.Println("Shutting consumer down")
		cancel()
		wg.Done()
		app.DB.Close()
		app.Analytics.Close()
	}()

	app.Stream = stream.New(app.Config.Stream.Brokers, app.Config.Stream.Topic)
	defer app.Stream.Close()

	log.Println("Consumer running...")

	consumer := consumers.New(ctx, &app)
	app.Stream.Consume(ctx, app.Config.Stream.Topic, "0", func(m *stream.Msg) error {
		return consumer.ProcessMsg(m)
	})

	wg.Add(1)
	wg.Wait()
}
