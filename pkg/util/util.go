package util

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"os"
	"os/user"
	"path/filepath"
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

func ContainerIdFromStateDirPath(p string) string {
	return filepath.Base(p)
}

// HasSudo checks if the current user has sudo (root) rights.
func HasSudo() bool {
	currentUser, err := user.Current()
	if err != nil {
		fmt.Println("Error getting current user:", err)
		return false
	}

	uid, err := strconv.Atoi(currentUser.Uid)
	if err != nil {
		fmt.Println("Error converting UID to integer:", err)
		return false
	}

	return uid == 0
}

func DebugLogToFile(s string) {
	//f, err := os.OpenFile("/tmp/rocilog", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
	//if err != nil {
	//	fmt.Println(err)
	//} else {
	//	_, err = fmt.Fprintf(f, s+"\n")
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//	f.Sync()
	//	f.Close()
	//}
}
