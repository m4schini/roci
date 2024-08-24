/*
Copyright © 2023 github.com/m4schini
*/
package cmd

import (
	"fmt"
	"os"
	"roci/pkg/model"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	configDir        = "/etc/roci"
	configDirFlag    = "configDir"
	containerDir     = "/run/roci/container"
	containerDirFlag = "containerDir"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "roci",
	Short: "Reduced Open Container Initiative runtime",
	Long: `roci is a experimental runtime only intended for testing purposes that implements only a subset of the oci
runtime specification. This means while the most basic container operations should be compatible enough
to run, the results and processes might behave different`,
	RunE: CommonPreRunE,
	FParseErrWhitelist: cobra.FParseErrWhitelist{
		UnknownFlags: true,
	},
	SilenceUsage: true,
}

// ExecuteCLI adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func ExecuteCLI() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(model.ExitCode(err))
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	defaultConfig()

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is /etc/roci/config.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().Bool("systemd-cgroup", false, "This only exists for runc-like cli compatibilty. Does nothing")
	rootCmd.PersistentFlags().MarkDeprecated("systemd-cgroup", "This only exists for runc-like cli compatibilty")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".roci" (without extension).
		viper.AddConfigPath("/etc/roci")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func defaultConfig() {
	viper.SetDefault(configDirFlag, configDir)
	viper.SetDefault(containerDirFlag, containerDir)
}
