package ipc

import (
	"context"
	"encoding/binary"
	"fmt"
	"google.golang.org/protobuf/proto"
	"io"
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

func openPipeReader(stateDir, pipeName string) (*os.File, error) {
	fifoPath := filepath.Join(stateDir, pipeName)
	return os.OpenFile(fifoPath, os.O_RDONLY|os.O_CREATE, os.ModeNamedPipe)
}

func openPipeWriter(stateDir, pipeName string) (*os.File, error) {
	fifoPath := filepath.Join(stateDir, pipeName)
	return os.OpenFile(fifoPath, os.O_WRONLY|os.O_CREATE, os.ModeNamedPipe)
}

func write(pipe io.Writer, message proto.Message) error {
	payload, err := proto.Marshal(message)
	if err != nil {
		return err
	}

	length := uint32(len(payload))
	if err := binary.Write(pipe, binary.LittleEndian, length); err != nil {
		return fmt.Errorf("failed to write payload length: %w", err)
	}

	if _, err := pipe.Write(payload); err != nil {
		return fmt.Errorf("failed to write payload: %w", err)
	}

	return nil
}

func read(pipe io.Reader) (part []byte, read bool, err error) {
	var length uint32
	err = binary.Read(pipe, binary.LittleEndian, &length)
	switch {
	case err == io.EOF:
		return part, false, nil
	case err != nil:
		return part, false, err
	}

	part = make([]byte, length)
	_, err = pipe.Read(part)
	return part, true, nil
}

func Listen[T proto.Message](ctx context.Context, pipe io.Reader, newInstance func() T) chan T {
	ch := make(chan T, 1)
	go func() {
		defer close(ch)
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			part, read, err := read(pipe)
			if err != nil {
				return
			}
			if !read {
				// If the file is not correctly closed, it leads to active waiting.
				// To mitigate this the runtime sleeps to reduce cpu strain
				time.Sleep(1 * time.Millisecond)
				continue
			}

			var message = newInstance()
			err = proto.Unmarshal(part, message)
			if err != nil {
				break
			}

			ch <- message
		}
	}()

	return ch
}

//func listen[T proto.Message](ctx context.Context, pipe *os.File) chan T {
//	ctx, cancel := context.WithCancel(ctx)
//	ch := make(chan T, 1)
//	go func() {
//		defer cancel()
//		defer close(ch)
//		var err error
//		for rawmessage := range Listen(ctx, pipe) {
//			var message T
//			err = proto.Unmarshal(rawmessage, message)
//			if err != nil {
//				break
//			}
//
//			ch <- message
//		}
//	}()
//	return ch
//}

func readFromInit(ctx context.Context, fifo *os.File) chan *pb.FromInit {
	return Listen(ctx, fifo, func() *pb.FromInit {
		return new(pb.FromInit)
	})
	//ctx, cancel := context.WithCancel(ctx)
	//ch := make(chan *pb.FromInit, 1)
	//go func() {
	//	defer cancel()
	//	defer close(ch)
	//	for payload := range Listen(ctx, fifo) {
	//		var message pb.FromInit
	//		err := proto.Unmarshal(payload, &message)
	//		if err != nil {
	//			break
	//		}
	//
	//		ch <- &message
	//	}
	//}()
	//return ch
}

func readFromRuntime(ctx context.Context, fifo *os.File) chan *pb.FromRuntime {
	return Listen(ctx, fifo, func() *pb.FromRuntime {
		return new(pb.FromRuntime)
	})
	//ctx, cancel := context.WithCancel(ctx)
	//ch := make(chan *pb.FromRuntime, 1)
	//go func() {
	//	defer cancel()
	//	defer close(ch)
	//	for payload := range Listen(ctx, fifo) {
	//		var message pb.FromRuntime
	//		err := proto.Unmarshal(payload, &message)
	//		if err != nil {
	//			break
	//		}
	//
	//		ch <- &message
	//	}
	//}()
	//return ch
}
