package ipc

import (
	"context"
	"fmt"
	"os"
	pb "roci/proto"
)

const initPipeFileName = "init.pipe"

func CreateInitPipe(stateDir string) error {
	return createPipe(stateDir, initPipeFileName)
}

type InitPipeReader interface {
	WaitForStart() error
}

func NewInitPipeReader(stateDir string) (InitPipeReader, error) {
	p := new(InitPipe)
	fd, err := newReader(stateDir, initPipeFileName)
	if err != nil {
		return nil, err
	}
	p.fd = fd
	return p, nil
}

type InitPipeWriter interface {
	SendStart() error
}

func NewInitPipeWriter(stateDir string) (InitPipeWriter, error) {
	p := new(InitPipe)
	fd, err := newWriter(stateDir, initPipeFileName)
	if err != nil {
		return nil, err
	}
	p.fd = fd
	return p, nil
}

type InitPipe struct {
	fd *os.File
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

func (i *InitPipe) SendStart() error {
	return write(i.fd, NewMessageStart())
}
