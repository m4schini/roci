package procfs

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	// uidMapFileName is the name of the file that stores UID mapping information.
	uidMapFileName = "uid.map"

	// gidMapFileName is the name of the file that stores GID mapping information.
	gidMapFileName = "gid.map"

	// mapFilePermissions defines the file permissions for the UID/GID map files.
	mapFilePermissions = 0666
)

// IdMapper writes mappings for UIDs and GIDs for a given process ID.
type IdMapper interface {
	// MapUid maps a UID from inside a container to an outside UID for the given process ID (pid).
	MapUid(pid Pid, insideId, outsideId uint32) error

	// MapGid maps a GID from inside a container to an outside GID for the given process ID (pid).
	MapGid(pid Pid, insideId, outsideId uint32) error
}

// MapUid maps a UID from inside a container to an outside UID using the uid.map file.
// It writes the mapping to the appropriate file in the specified procfs path for the process ID (pid).
func (F *FS) MapUid(pid Pid, insideId, outsideId uint32) error {
	return mapId(F.procfsPath, pid, uidMapFileName, insideId, outsideId)
}

// MapGid maps a GID from inside a container to an outside GID using the gid.map file.
// It writes the mapping to the appropriate file in the specified procfs path for the process ID (pid).
func (F *FS) MapGid(pid Pid, insideId, outsideId uint32) error {
	return mapId(F.procfsPath, pid, gidMapFileName, insideId, outsideId)
}

// mapId is a helper function that performs the actual ID mapping.
// It writes the mapping of IDs to the appropriate map file (uid.map or gid.map) in the specified base path.
func mapId(basePath string, pid Pid, mapFile string, insideId, outsideId uint32) error {
	path := filepath.Join(basePath, pid.String(), mapFile)
	mapping := fmt.Sprintf("%d %d %d", insideId, outsideId, 1)
	return os.WriteFile(path, []byte(mapping), mapFilePermissions)
}
