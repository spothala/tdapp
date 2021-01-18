package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var debug bool
var clientID string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "hook",
	Short: "A Server to support webhooks",
	Long:  `Currently Supporting Git Webhooks`,
}

// Execute  method execution
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVarP(&clientID, "clientid", "i", "", "Consumer  Key of the TD App")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

}
