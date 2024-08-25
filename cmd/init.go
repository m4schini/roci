package cmd

import (
	"roci/pkg/libcontainer"
)

const (
	InitCommandName = "init"
)

// ExecuteInit starts the init process with the data from the specified state directory.
// This is called by main.main().
func ExecuteInit(stateDir string) (err error) {
	return libcontainer.InitFromStateDir(stateDir)
}
