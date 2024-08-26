package ipc

import (
	"context"
	"go.uber.org/zap"
	"io"
	"os"
	"roci/pkg/logger"
	"roci/pkg/procfs"
	pb "roci/proto"
)

const runtimePipeFileName = "runtime.pipe"

// RuntimePipeWriter is the runtime.pipe interface for the init process
type RuntimePipeWriter interface {
	SendReady() error

	procfs.IdMapper
}

// RuntimePipeR is the cmd process implementation for the runtime.pipe message handlers
type RuntimePipeR struct {
	fd                 *os.File
	idMapper           procfs.IdMapper
	readyContext       context.Context
	finishReadyContext context.CancelFunc
}

// CreateRuntimePipe creates runtime.pipe fifo special file inside stateDir
func CreateRuntimePipe(stateDir string) error {
	return createPipe(stateDir, runtimePipeFileName)
}

// NewRuntimePipeReader opens the runtime pipe with read access and returns InitPipeReader
func NewRuntimePipeReader(ctx context.Context, stateDir string, idMapper procfs.IdMapper) (ready <-chan struct{}, closer io.Closer, err error) {
	p := new(RuntimePipeR)
	fd, err := openPipeReader(stateDir, runtimePipeFileName)
	if err != nil {
		return nil, nil, err
	}
	p.fd = fd
	p.idMapper = idMapper

	return p.listen(ctx), fd, nil
}

// listen reads messages from the runtime.pipe and handles them
func (r *RuntimePipeR) listen(ctx context.Context) (waitForReady <-chan struct{}) {
	r.readyContext, r.finishReadyContext = context.WithCancel(context.Background())
	log := logger.Log().Named("pipe").Named("runtime")
	log.Debug("listening on runtime pipe")

	go func() {
		for msg := range listenRuntimePipe(ctx, r.fd) {
			switch msg.Payload.(type) {
			case *pb.FromInit_Ready:
				log.Debug("received ready message")
				r.onReady()
				return
			case *pb.FromInit_MapGid:
				log.Debug("received map gid request")
				req := msg.GetMapGid()
				//TODO fix, pid self is runtime pid not init pid, needs to be init pid
				// (no performance impact expected)
				err := r.idMapper.MapGid(procfs.PidSelf, req.InsideId, req.OutsideId)
				if err != nil {
					log.Error("map gid failed", zap.Error(err))
				}
			case *pb.FromInit_MapUid:
				log.Debug("received map uid request")
				req := msg.GetMapUid()
				//TODO fix, pid self is runtime pid not init pid, needs to be init pid
				// (no performance impact expected)
				err := r.idMapper.MapUid(procfs.PidSelf, req.InsideId, req.OutsideId)
				if err != nil {
					log.Error("map uid failed", zap.Error(err))
				}
			}
		}
	}()

	return r.readyContext.Done()
}

// onReady implements the handling of the ready message
func (r *RuntimePipeR) onReady() {
	r.finishReadyContext()
}

// RuntimePipeW is the init process implementation of the RuntimePipeWriter
type RuntimePipeW struct {
	fd *os.File
}

// NewRuntimePipeWriter opens the runtime.pipe with write access and returns RuntimePipeWriter
func NewRuntimePipeWriter(stateDir string) (RuntimePipeWriter, error) {
	p := new(RuntimePipeW)
	fd, err := openPipeWriter(stateDir, runtimePipeFileName)
	if err != nil {
		return nil, err
	}
	p.fd = fd
	return p, nil
}

// SendReady sends the ready message
func (r RuntimePipeW) SendReady() error {
	return write(r.fd, NewMessageReady())
}

// MapUid sends the UidMapping message
func (r RuntimePipeW) MapUid(pid procfs.Pid, insideId, outsideId uint32) error {
	return write(r.fd, NewMessageUidMapping(insideId, outsideId))
}

// MapGid sends the GidMapping message
func (r RuntimePipeW) MapGid(pid procfs.Pid, insideId, outsideId uint32) error {
	return write(r.fd, NewMessageGidMapping(insideId, outsideId))
}

func (r RuntimePipeW) Close() error {
	if r.fd != nil {
		return r.fd.Close()
	}
	return nil
}
