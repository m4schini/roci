/*
Copyright © 2023 github.com/m4schini
*/
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
	Example: `For example, if the container id is "ubuntu01" and runc list currently shows the
status of "ubuntu01" as "stopped" the following will delete resources held for
"ubuntu01" removing "ubuntu01" from the runc list of containers:

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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	deleteCmd.Flags().BoolP("force", "f", false, `Forcibly deletes the container if it is still running (uses SIGKILL)`)
}
