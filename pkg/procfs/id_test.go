package procfs

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

const testPid = 42
const testPidStr = "42"

func newTestFs() (fs *FS, deleteFS func()) {
	testFSPath, err := os.MkdirTemp("", "procfs")
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll(filepath.Join(testFSPath, testPidStr), 0766)
	if err != nil {
		panic(err)
	}

	return &FS{procfsPath: testFSPath}, func() {
		os.RemoveAll(testFSPath)
	}
}

func TestFS_MapGid(t *testing.T) {
	fs, deleteFs := newTestFs()
	defer deleteFs()

	err := fs.MapGid(testPid, 0, 2)
	if err != nil {
		t.Error(err)
		return
	}

	mapFilePath := filepath.Join(fs.procfsPath, testPidStr, "gid.map")
	_, err = os.Stat(mapFilePath)
	if errors.Is(err, os.ErrNotExist) {
		t.Error(err)
		return
	}

	mapFile, err := os.ReadFile(mapFilePath)
	if err != nil {
		t.Error(err)
		return
	}

	if string(mapFile) != "0 2 1" {
		t.Log(string(mapFile))
		t.Fail()
	}
}

func TestFS_MapUid(t *testing.T) {
	fs, deleteFs := newTestFs()
	defer deleteFs()

	err := fs.MapUid(testPid, 0, 2)
	if err != nil {
		t.Error(err)
		return
	}

	mapFilePath := filepath.Join(fs.procfsPath, testPidStr, "uid.map")
	_, err = os.Stat(mapFilePath)
	if errors.Is(err, os.ErrNotExist) {
		t.Error(err)
		return
	}

	mapFile, err := os.ReadFile(mapFilePath)
	if err != nil {
		t.Error(err)
		return
	}

	if string(mapFile) != "0 2 1" {
		t.Log(string(mapFile))
		t.Fail()
	}
}
