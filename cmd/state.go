/*
Copyright © 2023 github.com/m4schini
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"roci/pkg/logger"
)

// stateCmd represents the state command
var stateCmd = &cobra.Command{
	Use:   "state <container-id>",
	Short: "output the state of a container",
	Long: `The state command outputs current state information for the
instance of a container.`,
	Args:    cobra.ExactArgs(1),
	PreRunE: ContainerPreRunE,
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			containerId = args[0]
			log         = logger.Log().Named("state")
		)
		log.Debug("state called", zap.String("containerId", containerId))

		state, err := confs.State(containerId)
		if err != nil {
			return err
		}

		output, err := json.MarshalIndent(state, "", "  ")
		if err != nil {
			return err
		}

		fmt.Println(string(output))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(stateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// stateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// stateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
