package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spothala/go-http-api/utils"
)

// githubCmd represents the github command
var verCmd = &cobra.Command{
	Use:   "version",
	Short: "Version of the Hook",
	Long:  `Version of the Hook Server`,
	Run: func(cmd *cobra.Command, args []string) {
		respJSON, err := utils.WriteJson(map[string]interface{}{"version": configYaml.Version,
			"description": configYaml.Description})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(utils.ReturnPrettyPrintJson(respJSON))
	},
}

func init() {
	RootCmd.AddCommand(verCmd)
}
