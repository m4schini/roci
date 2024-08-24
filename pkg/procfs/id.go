package procfs

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	uidMapFileName     = "uid.map"
	gidMapFileName     = "gid.map"
	mapFilePermissions = 0666
)

type IdMapper interface {
	MapUid(pid Pid, insideId, outsideId uint32) error
	MapGid(pid Pid, insideId, outsideId uint32) error
}

func (F *FS) MapUid(pid Pid, insideId, outsideId uint32) error {
	return mapId(F.BasePath, pid, uidMapFileName, insideId, outsideId)
}

func (F *FS) MapGid(pid Pid, insideId, outsideId uint32) error {
	return mapId(F.BasePath, pid, gidMapFileName, insideId, outsideId)
}

func mapId(basePath string, pid Pid, mapFile string, insideId, outsideId uint32) error {
	path := filepath.Join(basePath, pid.String(), mapFile)
	mapping := fmt.Sprintf("%d %d %d", insideId, outsideId, 1)
	return os.WriteFile(path, []byte(mapping), mapFilePermissions)
}
