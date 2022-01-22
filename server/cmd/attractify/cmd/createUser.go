package cmd

import (
	"fmt"
	"os"

	"attractify.io/platform/app"
	"attractify.io/platform/auth"
	"attractify.io/platform/config"
	"attractify.io/platform/db"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	_ "github.com/lib/pq"
	"regexp"
)

var createUserCmd = &cobra.Command{
	Use:   "create-user",
	Short: "Create user for Attractify",
	Long:  `Creating a new user for the Attractify backend.`,
	Run: func(cmd *cobra.Command, args []string) {
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			fmt.Println(err.Error())
		}

		email, err := cmd.Flags().GetString("email")
		if err != nil {
			fmt.Println(err.Error())
		}

		if !isEmailValid(email) {
			fmt.Println("E-mail not valid!")
			return
		}

		logger, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}
		app := app.App{Logger: logger}

		if len(os.Args) <= 1 {
			panic("missing config")
		}

		cfgPath := "/app/config.json"
		app.Config, err = config.Parse(cfgPath)
		if err != nil {
			panic(err)
		}

		dbConn, err := sqlx.Open("postgres", app.Config.DB)
		if err != nil {
			panic(err)
		}
		app.DB = db.New(dbConn)

		pw := auth.NewPassword("admin")
		ua := db.CreateUserParams{
			OrganizationID: uuid.FromStringOrNil("bc70b33d-c77f-4fe3-813d-a2605c0915cb"),
			Email:          email,
			Password:       pw.Password,
			Salt:           pw.Salt,
			Name:           name,
			Role:           db.RoleAdmin,
		}

		user, row := app.DB.CreateCLIUser(ua)
		if row != nil {
			panic(row.Error())
		}

		fmt.Println("User " + user.Name + " created!")

		app.DB.Close()
	},
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}

func init() {
	rootCmd.AddCommand(createUserCmd)

	createUserCmd.PersistentFlags().StringP("name", "n", "", "Specify the users name")
	createUserCmd.PersistentFlags().StringP("email", "e", "", "Specify the users email")
	createUserCmd.MarkPersistentFlagRequired("name")
	createUserCmd.MarkPersistentFlagRequired("email")
}
