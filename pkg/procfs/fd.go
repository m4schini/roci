package procfs

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

const (
	base = "/proc"
)

func AttachReader(pid, fd int, reader io.Reader) error {
	f, err := Open(pid, fmt.Sprintf("fd/%v", fd))
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, reader)
	return err
}

func AttachWriter(pid, fd int, writer io.Writer) error {
	f, err := Open(pid, fmt.Sprintf("fd/%v", fd))
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(writer, f)
	return err
}

func Open(pid int, path string) (*os.File, error) {
	pidDir := filepath.Join(base, strconv.Itoa(pid))
	_, err := os.Stat(pidDir)
	if err != nil {
		return nil, err
	}

	targetPath := filepath.Join(pidDir, path)
	return os.Open(targetPath)
}
