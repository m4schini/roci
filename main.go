/*
Copyright © 2023 github.com/m4schini
*/
package main

import (
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"roci/cmd"
	"roci/pkg/logger"
)

func main() {
	//s := strings.Join(os.Args, " ")
	//fmt.Println(s)
	//util.DebugLogToFile(s)
	defer logger.Sync()

	if isInitProcess() {
		executeInit()
	} else {
		executeCLI()
	}
}

func isInitProcess() bool {
	return len(os.Args) == 3 && os.Args[1] == cmd.InitCommandName
}

func executeInit() {
	logger.Set(logger.Log().Named("init").With(zap.String("cid", filepath.Base(os.Args[2]))))
	log := logger.Log()

	log.Debug("running container init")
	err := cmd.ExecuteInit(os.Args[2])
	if err != nil {
		log.Fatal("failed to run container init", zap.Error(err))
	}
}

func executeCLI() {
	logger.Set(logger.Log().Named("main"))
	log := logger.Log()

	log.Debug("running runtime cli")
	cmd.ExecuteCLI()
}
