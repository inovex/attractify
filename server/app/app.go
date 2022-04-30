package app

import (
	"attractify.io/platform/analytics"
	"attractify.io/platform/config"
	"attractify.io/platform/db"
	"attractify.io/platform/mailer"
	"go.uber.org/zap"
)

type App struct {
	Config    *config.Config
	DB        *db.DB
	Analytics *analytics.Analytics
	Mailer    *mailer.Mailer
	Logger    *zap.Logger
}
