package libcontainer

import (
	"github.com/opencontainers/runtime-spec/specs-go"
	"go.uber.org/zap"
	"path"
	"roci/pkg/libcontainer/initp"
	"roci/pkg/logger"
	"roci/pkg/model"
	"roci/pkg/util"
)

// InitFromStateDir initializes the container environment using the state directory.
// It reads the OCI specification from the state directory and then proceeds with initialization.
func InitFromStateDir(stateDir string) (err error) {
	log := logger.Log().Named("init")
	configPath := path.Join(stateDir, model.OciSpecFileName)
	var spec specs.Spec
	log.Debug("reading config", zap.String("path", configPath))
	err = util.ReadJsonFile(configPath, &spec)
	if err != nil {
		log.Warn("failed to read spec", zap.Error(err), zap.String("configPath", configPath))
		return err
	}
	log.Debug("read config", zap.Any("config", spec))
	return initp.Init(stateDir, spec)
}
