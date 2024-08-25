package ipc

import (
	"context"
	"fmt"
	"os"
	pb "roci/proto"
)

const initPipeFileName = "init.pipe"

type InitPipeWriter interface {
	SendStart() error
}

type InitPipeReader interface {
	WaitForStart() error
}

type InitPipe struct {
	fd *os.File
}

func CreateInitPipe(stateDir string) error {
	return createPipe(stateDir, initPipeFileName)
}

func NewInitPipeReader(stateDir string) (InitPipeReader, error) {
	p := new(InitPipe)
	fd, err := openPipeReader(stateDir, initPipeFileName)
	if err != nil {
		return nil, err
	}
	p.fd = fd
	return p, nil
}

func (i *InitPipe) WaitForStart() error {
	ch := readFromRuntime(context.TODO(), i.fd)
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

func NewInitPipeWriter(stateDir string) (InitPipeWriter, error) {
	p := new(InitPipe)
	fd, err := openPipeWriter(stateDir, initPipeFileName)
	if err != nil {
		return nil, err
	}
	p.fd = fd
	return p, nil
}

func (i *InitPipe) SendStart() error {
	return write(i.fd, NewMessageStart())
}
