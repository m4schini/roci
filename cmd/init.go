package cmd

import (
	"roci/pkg/libcontainer"
)

const (
	InitCommandName = "init"
)

func ExecuteInit(stateDir string) (err error) {
	return libcontainer.InitFromStateDir(stateDir)
}
