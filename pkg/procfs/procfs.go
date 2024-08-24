package procfs

const (
	defaultPath = "/proc"
)

var Root = &FS{BasePath: defaultPath}

type FS struct {
	BasePath string
}
