package ipc

import (
	"context"
	"fmt"
	"os"
	pb "roci/proto"
)

const initPipeFileName = "init.pipe"

// InitPipeWriter is the init.pipe interface for the cmd process
type InitPipeWriter interface {
	SendStart() error
}

// InitPipeReader is the init.pipe interface for the init process
type InitPipeReader interface {
	WaitForStart() error
}

// InitPipe contains init.pipe fd
type InitPipe struct {
	fd *os.File
}

// CreateInitPipe creates init.pipe fifo special file inside container statedir
func CreateInitPipe(stateDir string) error {
	return createPipe(stateDir, initPipeFileName)
}

// NewInitPipeReader opens the init pipe with read access and returns InitPipeReader
func NewInitPipeReader(stateDir string) (InitPipeReader, error) {
	p := new(InitPipe)
	fd, err := openPipeReader(stateDir, initPipeFileName)
	if err != nil {
		return nil, err
	}
	p.fd = fd
	return p, nil
}

// WaitForStart waits for the start message from the cmd process.
// If an unknown message is received it returns an error
func (i *InitPipe) WaitForStart() error {
	ch := listenInitPipe(context.TODO(), i.fd)
	for msg := range ch {
		switch msg.Payload.(type) {
		case *pb.FromRuntime_Start:
			return nil
		default:
			return fmt.Errorf("unknown message")
		}
	}
	return fmt.Errorf("pipe closed without start signal")
}

// NewInitPipeWriter opens the init pipe with write access and returns InitPipeWriter
func NewInitPipeWriter(stateDir string) (InitPipeWriter, error) {
	p := new(InitPipe)
	fd, err := openPipeWriter(stateDir, initPipeFileName)
	if err != nil {
		return nil, err
	}
	p.fd = fd
	return p, nil
}

// SendStart sends the start message over the init pipe
func (i *InitPipe) SendStart() error {
	return write(i.fd, NewMessageStart())
}
