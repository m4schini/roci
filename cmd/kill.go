package cmd

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"roci/pkg/logger"
	"roci/pkg/model"
)

const (
	defaultSignal = "SIGTERM"
)

// killCmd represents the kill command
var killCmd = &cobra.Command{
	Use:   "kill [command options] <container-id> [signal]",
	Short: "kill sends the specified signal (default: SIGTERM) to the container's init process",
	Example: `For example, if the container id is "ubuntu01" the following will send a "KILL"
signal to the init process of the "ubuntu01" container:

       # runc kill ubuntu01 SIGKILL`,
	Args:    cobra.RangeArgs(1, 2),
	PreRunE: ContainerPreRunE,
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			containerId = args[0]
			signalName  = getSignal(args, 1)
			log         = logger.Log().Named("kill")
		)
		log.Debug("kill called", zap.String("containerId", containerId), zap.String("signal", signalName))

		signal, err := model.SyscallSignal(signalName)
		if err != nil {
			return err
		}

		return confs.Kill(containerId, signal)
	},
}

func getSignal(args []string, position int) string {
	if position-1 >= len(args) {
		return args[position]
	} else {
		return defaultSignal
	}
}

func init() {
	rootCmd.AddCommand(killCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// killCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	killCmd.Flags().BoolP("all", "a", false, "send the specified signal to all processes inside the container")
	killCmd.Flags().MarkDeprecated("all", "does nothing, only exists for runc compatibility")
}
