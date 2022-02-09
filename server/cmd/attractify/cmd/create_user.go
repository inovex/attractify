package cmd

import (
	"crypto/rand"
	"errors"
	"fmt"
	"strings"
	"syscall"
	"time"

	"attractify.io/platform/auth"
	"attractify.io/platform/db"
	"github.com/spf13/cobra"

	"regexp"

	_ "github.com/lib/pq"
	"golang.org/x/term"
)

var createUserCmd = &cobra.Command{
	Use:   "create-user",
	Short: "Create an Attractify organization and user",
	Long:  `Creating an initial organization and user for Attractify`,
	Run:   handleUserCreateCmd,
}

func init() {
	rootCmd.AddCommand(createUserCmd)

	createUserCmd.PersistentFlags().StringP("user", "u", "", "Full name of the new user")
	createUserCmd.PersistentFlags().StringP("organization", "o", "", "Name of the organization")
	createUserCmd.PersistentFlags().StringP("email", "e", "", "The user's email address")
	createUserCmd.PersistentFlags().StringP("timezone", "t", "UTC", "The organization's timezone in the TZ format (e.g. America/New_York). Defaults to UTC")

	createUserCmd.MarkPersistentFlagRequired("user")
	createUserCmd.MarkPersistentFlagRequired("organization")
	createUserCmd.MarkPersistentFlagRequired("email")
}

func handleUserCreateCmd(cmd *cobra.Command, args []string) {
	// Gather commandline flags
	userName, err := cmd.Flags().GetString("user")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	organizationName, err := cmd.Flags().GetString("organization")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	email, err := cmd.Flags().GetString("email")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	email = strings.ToLower(email)
	if !isEmailValid(email) {
		fmt.Println("Email is not valid")
		return
	}
	timezone, err := cmd.Flags().GetString("timezone")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if !isValidTimezone(timezone) {
		fmt.Println("Time zone is not valid. Please use the IANA time zone format.")
		return
	}

	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		fmt.Println(err.Error())
		return
	}

	org, err := cliApp.DB.CreateOrganization(
		ctx,
		organizationName,
		timezone,
		key,
	)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	password, err := readPassword()
	if err != nil {
		fmt.Println("could not read password from stdin")
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
	user, err := cliApp.DB.CreateUser(ctx, userParams)
	if err != nil {
		fmt.Printf("could not create user: %s", err.Error())
		return
	}

	fmt.Printf("Organization %s has been created.\n", org.Name)
	fmt.Printf("User %s has been created.\n", user.Name)
	fmt.Printf("You can now login with your email address %s and your password.\n", user.Email)
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

func readPassword() (string, error) {
	pw, err := requestPassword("Enter Password: ")
	if err != nil {
		return "", err
	}

	pwConfirmation, err := requestPassword("Confirm Password: ")
	if err != nil {
		return "", err
	}

	if pw != pwConfirmation {
		return "", errors.New("passwords do not match")
	}

	return pw, nil
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}

func isValidTimezone(timezone string) bool {
	_, err := time.LoadLocation(timezone)
	return err == nil
}
