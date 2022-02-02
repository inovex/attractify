package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"

	"attractify.io/platform/app"
	_ "attractify.io/platform/auth"
	"attractify.io/platform/config"
	"attractify.io/platform/db"
	"github.com/gofrs/uuid"
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

		email = strings.ToLower(email)

		if !isEmailValid(email) {
			fmt.Println("Email is not valid")
			return
		}

		logger, err := zap.NewProduction()
		handleErrorPanic(err)

		app := app.App{Logger: logger}

		if len(os.Args) <= 1 {
			panic("missing config")
		}

		cfgPath := "/app/config.json"
		app.Config, err = config.Parse(cfgPath)
		handleErrorPanic(err)

		dbConn, err := sqlx.Open("postgres", app.Config.DB)
		handleErrorPanic(err)

		app.DB = db.New(dbConn)

		organizationID := uuid.Must(uuid.NewV4())

		password, err := checkPasswordFromCLI()
		handleErrorPanic(err)

		//Test-Prints
		fmt.Println("pw: " + password)
		fmt.Println("orgID: " + organizationID.String())
		fmt.Println("name: " + name)
		fmt.Println("email: " + email)

		//TODO: create Organization

		/*pw := auth.NewPassword(password)
		ua := db.CreateUserParams{
			OrganizationID: organizationID,
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

		fmt.Sprintf("User %s has been created", user.Name)*/

		app.DB.Close()
	},
}

func requestPassword(msg string) (string, error) {
	fmt.Print(msg)
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	handleErrorPanic(err)

	if len(bytePassword) < 8 {
		fmt.Print("\nPassword is too short\n")
		return "", errors.New("password too short")
	}

	fmt.Println("")
	return strings.TrimSpace(string(bytePassword)), nil
}

func checkPasswordFromCLI() (string, error) {
	pw1, err := requestPassword("Enter Password: ")
	handleErrorPanic(err)

	pw2, err := requestPassword("Confirm Password: ")
	handleErrorPanic(err)

	if pw1 != pw2 {
		fmt.Println("Passwords do not match")
		return "", errors.New("passwords do not match")
	}

	return pw1, nil
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}

func handleErrorPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func CreateOrganization(name string, email string) {
	/* var req requests.OrganizationCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		return
	}

	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		oc.App.Logger.Error("organizations.create.genKey", zap.Error(err))
		return
	}

	org, err := oc.App.DB.CreateOrganization(
		c.Request.Context(),
		req.OrganizationName,
		req.Timezone,
		key,
	)
	if err != nil {
		oc.App.Logger.Error("organizations.create.createOrganization", zap.Error(err))
		return
	} */

	/* println(name)

	pw := auth.NewPassword(req.Password)
	ua := db.CreateUserParams{
		OrganizationID: org.ID,
		Email:          req.Email,
		Password:       pw.Password,
		Salt:           pw.Salt,
		Name:           req.Name,
		Role:           db.RoleAdmin,
	}
	user, err := oc.App.DB.CreateUser(c.Request.Context(), ua)
	if err != nil {
		oc.App.Logger.Error("organizations.create.createUser", zap.Error(err))
		return
	} */

}

func init() {
	rootCmd.AddCommand(createUserCmd)

	createUserCmd.PersistentFlags().StringP("name", "n", "", "Specify the users name")
	createUserCmd.PersistentFlags().StringP("email", "e", "", "Specify the users email")
	createUserCmd.MarkPersistentFlagRequired("name")
	createUserCmd.MarkPersistentFlagRequired("email")
}
