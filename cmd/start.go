package cmd

import (
	"time"

	"github.com/spf13/cobra"
	api "github.com/spothala/tdapp/api"
)

var port int

// githubCmd represents the github command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the server",
	Long:  `Starts the hooks server on the port specified`,
	Run: func(cmd *cobra.Command, args []string) {
		server, _ := api.Start(port, clientID, debug)
		api.GracefulShutdown(server, 10*time.Second)
	},
}

func init() {
	RootCmd.AddCommand(startCmd)

	// Add global flags for github command
	startCmd.PersistentFlags().IntVarP(&port, "port", "p", 8080, "Port to start the server")
}
