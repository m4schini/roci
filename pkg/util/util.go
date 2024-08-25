package util

import (
	"encoding/json"
	"go.uber.org/zap"
	"os"
	"os/user"
	"strconv"
)

func ReadJsonFile(path string, v any) error {
	zap.L().Debug("reading json file", zap.String("path", path))
	defer zap.L().Debug("done reading json file", zap.String("path", path))
	f, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer func() {
		err = f.Close()
		if err != nil {
			zap.L().Warn("failed to close file", zap.Error(err))
		}
		zap.L().Debug("closed file")
	}()

	zap.L().Debug("decoding json")
	err = json.NewDecoder(f).Decode(v)
	if err != nil {
		return err
	}
	zap.L().Debug("decoded json")
	return nil
}

func WriteJsonFile(path string, v any) error {
	zap.L().Debug("writing json file", zap.String("path", path))
	defer zap.L().Debug("done writing json file", zap.String("path", path))
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0o711)
	if err != nil {
		return err
	}
	defer func() {
		err = f.Close()
		if err != nil {
			zap.L().Warn("failed to close file", zap.Error(err))
		}
		zap.L().Debug("closed file")
	}()

	return json.NewEncoder(f).Encode(v)
}

// HasSudo checks if the current user has sudo (root) privileges.
func HasSudo() bool {
	currentUser, err := user.Current()
	if err != nil {
		return false
	}

	uid, err := strconv.Atoi(currentUser.Uid)
	if err != nil {
		return false
	}

	return uid == 0
}
