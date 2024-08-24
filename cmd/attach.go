/*
Copyright © 2023 github.com/m4schini
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"roci/pkg/procfs"
)

// attachCmd represents the attach command
var attachCmd = &cobra.Command{
	Use:     "attach <container-id>",
	Short:   "Attach to stdin, stdout & stderr container",
	Args:    cobra.ExactArgs(1),
	PreRunE: ContainerPreRunE,
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			containerId = args[0]
		)

		c, err := confs.Load(containerId)
		if err != nil {
			return err
		}
		pid := c.State().Pid
		fmt.Println(c.State().Pid)

		//TODO attach to file descriptor in rootfs?
		go func() {
			err = procfs.AttachReader(pid, 1, os.Stdout)
			if err != nil {
				panic(err)
			}
		}()

		go func() {
			err = procfs.AttachReader(pid, 2, os.Stderr)
			if err != nil {
				panic(err)
			}
		}()

		return procfs.AttachWriter(pid, 0, os.Stdin)
	},
}

func init() {
	rootCmd.AddCommand(attachCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// attachCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// attachCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
