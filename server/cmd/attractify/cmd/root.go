package cmd

import (
	"context"
	"fmt"
	"os"

	"attractify.io/platform/app"
	"attractify.io/platform/config"
	"attractify.io/platform/db"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	cfgFile string
	ctx     context.Context
	cliApp  *app.App
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "attractify",
	Short: "Attractify CLI tool",
	Long:  `The Attractify CLI tool helps you to perform basic admin operations on your Attractify installation.`,
}

func init() {
	cobra.OnInitialize(initApp)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func initApp() {
	ctx = context.Background()

	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	cliApp = &app.App{Logger: logger}
	cliApp.Config, err = config.Parse(cfgFile)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	dbConn, err := sqlx.Open("postgres", cliApp.Config.DB)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	cliApp.DB = db.New(dbConn)
}
