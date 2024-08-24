package ipc

import (
	"context"
	"encoding/binary"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"os"
	"path/filepath"
	pb "roci/proto"
	"syscall"
	"time"
)

func createPipe(stateDir, pipeName string) error {
	fifoPath := filepath.Join(stateDir, pipeName)
	return syscall.Mkfifo(fifoPath, 0666)
}

func newReader(stateDir, pipeName string) (*os.File, error) {
	fifoPath := filepath.Join(stateDir, pipeName)
	return os.OpenFile(fifoPath, os.O_RDONLY|os.O_CREATE, os.ModeNamedPipe)
}

func newWriter(stateDir, pipeName string) (*os.File, error) {
	fifoPath := filepath.Join(stateDir, pipeName)
	return os.OpenFile(fifoPath, os.O_WRONLY|os.O_CREATE, os.ModeNamedPipe)
}

func write(fifo *os.File, message proto.Message) error {
	payload, err := proto.Marshal(message)
	if err != nil {
		return err
	}

	length := uint32(len(payload))
	if err := binary.Write(fifo, binary.LittleEndian, length); err != nil {
		return fmt.Errorf("failed to write payload length: %w", err)
	}

	// Then, write the actual payload
	if _, err := fifo.Write(payload); err != nil {
		return fmt.Errorf("failed to write payload: %w", err)
	}

	return nil
}

func read(ctx context.Context, fifo *os.File) chan []byte {
	ch := make(chan []byte, 1)
	go func() {
		defer close(ch)
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			zap.L().Named("pipe").With(zap.String("name", fifo.Name())).Debug("reading")
			// First, read the length of the next payload
			var length uint32
			if err := binary.Read(fifo, binary.LittleEndian, &length); err != nil {
				if err.Error() == "EOF" {
					zap.L().Named("pipe").With(zap.String("name", fifo.Name())).Debug("EOF")
					time.Sleep(1 * time.Millisecond)
					continue // End of file
				}
				return
			}

			// Then, read the actual payload
			payload := make([]byte, length)
			if _, err := fifo.Read(payload); err != nil {
				return
			}

			ch <- payload
		}
	}()

	return ch
}

func readFromInit(ctx context.Context, fifo *os.File) chan *pb.FromInit {
	ctx, cancel := context.WithCancel(ctx)
	ch := make(chan *pb.FromInit, 1)
	go func() {
		defer cancel()
		defer close(ch)
		for payload := range read(ctx, fifo) {
			var message pb.FromInit
			err := proto.Unmarshal(payload, &message)
			if err != nil {
				break
			}

			ch <- &message
		}
	}()
	return ch
}

func readFromRuntime(ctx context.Context, fifo *os.File) chan *pb.FromRuntime {
	ctx, cancel := context.WithCancel(ctx)
	ch := make(chan *pb.FromRuntime, 1)
	go func() {
		defer cancel()
		defer close(ch)
		for payload := range read(ctx, fifo) {
			var message pb.FromRuntime
			err := proto.Unmarshal(payload, &message)
			if err != nil {
				break
			}

			ch <- &message
		}
	}()
	return ch
}
