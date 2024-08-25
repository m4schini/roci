package cmd

import (
	"go.uber.org/zap"
	"roci/pkg/logger"
	"roci/pkg/model"
	"syscall"

	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete [command options] <container-id>",
	Short: "delete any resources held by the container often used with detached container",
	Example: `For example, if the container id is "ubuntu01" and roci list currently shows the
status of "ubuntu01" as "stopped" the following will delete resources held for
"ubuntu01" removing "ubuntu01" from the roci list of containers:

       # runc delete ubuntu01`,
	Args:    cobra.ExactArgs(1),
	PreRunE: ContainerPreRunE,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var (
			containerId = args[0]
			forceFlag   = MustGetBool(cmd, "force")
			log         = logger.Log().Named("delete")
		)
		log.Debug("delete called", zap.String("containerId", containerId), zap.Bool("force", forceFlag))
		if forceFlag {
			log.Debug("force delete called")
			err = confs.Kill(containerId, syscall.SIGKILL)
			log.Debug("tried to kill process", zap.Error(err))
		}

		err = confs.Remove(containerId)
		if err == model.ErrNotExist {
			log.Debug("tried to delete container that doesn't exist")
			return nil
		}
		return err
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().BoolP("force", "f", false, `Forcibly deletes the container if it is still running (uses SIGKILL)`)
}
