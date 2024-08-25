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

func CreateRuntimePipe(stateDir string) error {
	return createPipe(stateDir, runtimePipeFileName)
}

type RuntimePipeR struct {
	fd                 *os.File
	idMapper           procfs.IdMapper
	readyContext       context.Context
	finishReadyContext context.CancelFunc
}

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

func (r *RuntimePipeR) listen(ctx context.Context) (waitForReady <-chan struct{}) {
	r.readyContext, r.finishReadyContext = context.WithCancel(context.Background())
	log := logger.Log().Named("pipe").Named("runtime")
	log.Debug("listening on runtime pipe")

	go func() {
		for msg := range readFromInit(ctx, r.fd) {
			switch msg.Payload.(type) {
			case *pb.FromInit_Ready:
				log.Debug("received ready message")
				r.onReady()
				return
			case *pb.FromInit_MapGid:
				log.Debug("received map gid request")
				req := msg.GetMapGid()
				//TODO fix, pid self is runtime pid not init pid, needs to be init pid
				err := r.idMapper.MapGid(procfs.PidSelf, req.InsideId, req.OutsideId)
				if err != nil {
					log.Error("map gid failed", zap.Error(err))
				}
			case *pb.FromInit_MapUid:
				log.Debug("received map uid request")
				req := msg.GetMapUid()
				//TODO fix, pid self is runtime pid not init pid, needs to be init pid
				err := r.idMapper.MapUid(procfs.PidSelf, req.InsideId, req.OutsideId)
				if err != nil {
					log.Error("map uid failed", zap.Error(err))
				}
			}
		}
	}()

	return r.readyContext.Done()
}

func (r *RuntimePipeR) onReady() {
	r.finishReadyContext()
}

type RuntimePipeWriter interface {
	SendReady() error

	procfs.IdMapper
}

type RuntimePipeW struct {
	fd *os.File
}

func NewRuntimePipeWriter(stateDir string) (RuntimePipeWriter, error) {
	p := new(RuntimePipeW)
	fd, err := openPipeWriter(stateDir, runtimePipeFileName)
	if err != nil {
		return nil, err
	}
	p.fd = fd
	return p, nil
}

func (r RuntimePipeW) SendReady() error {
	return write(r.fd, NewMessageReady())
}

func (r RuntimePipeW) MapUid(pid procfs.Pid, insideId, outsideId uint32) error {
	return write(r.fd, NewMessageUidMapping(insideId, outsideId))
}

func (r RuntimePipeW) MapGid(pid procfs.Pid, insideId, outsideId uint32) error {
	return write(r.fd, NewMessageGidMapping(insideId, outsideId))
}

func (r RuntimePipeW) Close() error {
	if r.fd != nil {
		return r.fd.Close()
	}
	return nil
}
