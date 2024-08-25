package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"roci/pkg/libcontainer"
	"roci/pkg/logger"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create [command options] <container−id>",
	Short: "Create a container",
	Long: `The create command creates an instance of a container for a bundle. The bundle
is a directory with a specification file named "config.json" and a root
filesystem.

The specification file includes an args parameter. The args parameter is used
to specify command(s) that get run when the container is started. To change the
command(s) that get executed on start, edit the args parameter of the spec.
`,
	Args:    cobra.ExactArgs(1),
	PreRunE: ContainerPreRunE,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var (
			containerId  = args[0]
			bundle       = MustGetString(cmd, "bundle")
			pidFile      = MustGetString(cmd, "pid-file")
			writePidFile = pidFile != ""
			log          = logger.Log().Named("create")
		)
		log.Debug("create called", zap.String("containerId", containerId), zap.String("bundle", bundle))
		defer log.Debug("create call handled")

		bundleAbs, err := filepath.Abs(bundle)
		if err != nil {
			return err
		}

		log.Debug("creating container")
		c, err := libcontainer.CreateContainer(confs, containerId, bundleAbs)
		if err != nil {
			return err
		}

		if writePidFile {
			pid := c.State().Pid
			log.Debug("writing pid file", zap.Int("pid", pid))
			err = writePid(pidFile, pid)
			if err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringP("bundle", "b", ".", `path to the root of the bundle directory, defaults to the current directory`)
	createCmd.Flags().String("pid-file", "", `specify the file to write the process id to`)
}

func writePid(pidFile string, pid int) error {
	return os.WriteFile(pidFile, []byte(fmt.Sprintf("%v\n", pid)), 0666)
}
