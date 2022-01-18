package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
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

		fmt.Println(name)
		fmt.Println(email)
	},
}

func init() {
	rootCmd.AddCommand(createUserCmd)

	createUserCmd.PersistentFlags().StringP("name", "n", "", "Specify the users name")
	createUserCmd.PersistentFlags().StringP("email", "e", "", "Specify the users email")
	createUserCmd.MarkPersistentFlagRequired("name")
	createUserCmd.MarkPersistentFlagRequired("email")
}
