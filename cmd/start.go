/*
Copyright © 2023 github.com/m4schini
*/
package cmd

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"roci/pkg/logger"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:     "start <container-id>",
	Short:   "executes the user defined process in a created container",
	Long:    `The start command executes the user defined process in a created container.`,
	Args:    cobra.ExactArgs(1),
	PreRunE: ContainerPreRunE,
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			containerId = args[0]
			log         = logger.Log().Named("start")
		)
		log.Debug("start called", zap.String("containerId", containerId))

		return confs.Start(containerId)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
