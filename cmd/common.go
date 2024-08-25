package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"roci/pkg/libcontainer"
	"roci/pkg/model"
	"roci/pkg/util"
)

var confs *libcontainer.FS

func CommonPreRunE(cmd *cobra.Command, args []string) (err error) {
	if !util.HasSudo() {
		return model.ErrNoSudo
	}

	var allDirs = []string{viper.GetString(containerDirFlag), viper.GetString(configDirFlag)}
	for _, dir := range allDirs {
		err = os.MkdirAll(dir, 0o755)
		if err != nil {
			return err
		}
	}

	return nil
}

func ContainerPreRunE(cmd *cobra.Command, args []string) (err error) {
	if err = CommonPreRunE(cmd, args); err != nil {
		return err
	}

	if confs, err = libcontainer.NewContainerFS(viper.GetString(containerDirFlag)); err != nil {
		return err
	}

	return nil
}

func MustGetString(cmd *cobra.Command, key string) string {
	v, err := cmd.Flags().GetString(key)
	if err != nil {
		panic(err)
	}
	return v
}

func MustGetBool(cmd *cobra.Command, key string) bool {
	v, err := cmd.Flags().GetBool(key)
	if err != nil {
		panic(err)
	}
	return v
}
