package cmd

import (
	"crypto/rand"
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	"attractify.io/platform/app"
	"attractify.io/platform/auth"
	"attractify.io/platform/config"
	"attractify.io/platform/db"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"regexp"

	_ "github.com/lib/pq"
	"golang.org/x/term"
)

var createUserCmd = &cobra.Command{
	Use:   "create-user",
	Short: "Create user for Attractify",
	Long:  `Creating an initial user for the Attractify backend`,
	Run:   handleUserCmd,
}

func handleUserCmd(cmd *cobra.Command, args []string) {
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	userName, err := cmd.Flags().GetString("user-name")
	if err != nil {
		logger.Fatal(err.Error())
		return
	}

	organizationName, err := cmd.Flags().GetString("organization-name")
	if err != nil {
		logger.Fatal(err.Error())
		return
	}

	email, err := cmd.Flags().GetString("email")
	if err != nil {
		logger.Fatal(err.Error())
		return
	}
	email = strings.ToLower(email)

	if !isEmailValid(email) {
		logger.Fatal("Email is not valid")
		return
	}

	timeZone, err := cmd.Flags().GetString("time-zone")
	if err != nil {
		logger.Fatal(err.Error())
		return
	}

	timeZone, err = validateTimeZone(timeZone)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	app := app.App{Logger: logger}

	if len(os.Args) <= 1 {
		logger.Fatal("missing config")
		return
	}

	cfgPath := "/app/config.json"
	app.Config, err = config.Parse(cfgPath)
	if err != nil {
		logger.Fatal(err.Error())
		return
	}

	dbConn, err := sqlx.Open("postgres", app.Config.DB)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	app.DB = db.New(dbConn)
	defer app.DB.Close()

	password, err := checkPasswordFromCLI()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		logger.Fatal(err.Error())
		return
	}

	org, err := app.DB.CreateCLIOrganization(
		organizationName,
		timeZone,
		key,
	)
	if err != nil {
		logger.Fatal(err.Error())
		return
	}

	pw := auth.NewPassword(password)
	userParams := db.CreateUserParams{
		OrganizationID: org.ID,
		Email:          email,
		Password:       pw.Password,
		Salt:           pw.Salt,
		Name:           userName,
		Role:           db.RoleAdmin,
	}

	user, err := app.DB.CreateCLIUser(userParams)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Could not create user: %s\n", err.Error()))
		return
	}

	fmt.Printf("Organization %s has been created\n", org.Name)
	fmt.Printf("User %s has been created\n", user.Name)

}

func requestPassword(msg string) (string, error) {
	fmt.Print(msg)
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}

	fmt.Println("")
	if len(bytePassword) < 8 {
		return "", errors.New("password is too short")
	}

	return strings.TrimSpace(string(bytePassword)), nil
}

func checkPasswordFromCLI() (string, error) {
	pw1, err := requestPassword("Enter Password: ")
	if err != nil {
		return "", err
	}

	pw2, err := requestPassword("Confirm Password: ")
	if err != nil {
		return "", err
	}

	if pw1 != pw2 {
		return "", errors.New("passwords do not match")
	}

	return pw1, nil
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}

func validateTimeZone(timeZone string) (string, error) {
	if timeZone == "" {
		return "Europe/Berlin", nil
	}
	_, err := time.LoadLocation(timeZone)
	if err != nil {
		fmt.Println("Time zone is not valid")
		fmt.Println("Please use the IANA time zone format")
		return "", errors.New("")
	}
	return timeZone, nil
}

func init() {
	rootCmd.AddCommand(createUserCmd)

	createUserCmd.PersistentFlags().StringP("user-name", "u", "", "Specify the users name")
	createUserCmd.PersistentFlags().StringP("organization-name", "o", "", "Specify the users name")
	createUserCmd.PersistentFlags().StringP("email", "e", "", "Specify the users email")
	createUserCmd.PersistentFlags().StringP("time-zone", "t", "", "Specify the organizations time zone. Default: Europa/Berlin")

	createUserCmd.MarkPersistentFlagRequired("user-name")
	createUserCmd.MarkPersistentFlagRequired("organization-name")
	createUserCmd.MarkPersistentFlagRequired("email")
}
