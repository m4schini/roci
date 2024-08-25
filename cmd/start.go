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
}
