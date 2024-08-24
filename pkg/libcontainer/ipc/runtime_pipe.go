package ipc

import (
	"context"
	"go.uber.org/zap"
	"os"
	"roci/pkg/logger"
	"roci/pkg/procfs"
	pb "roci/proto"
)

const runtimePipeFileName = "runtime.pipe"

func CreateRuntimePipe(stateDir string) error {
	return createPipe(stateDir, runtimePipeFileName)
}

func NewRuntimePipeReader(stateDir string, idMapper procfs.IdMapper) (ready <-chan struct{}, err error) {
	p := new(RuntimePipeR)
	fd, err := newReader(stateDir, runtimePipeFileName)
	if err != nil {
		return nil, err
	}
	p.fd = fd
	p.idMapper = idMapper

	return p.listen(), nil
}

type RuntimePipeWriter interface {
	SendReady() error

	procfs.IdMapper
}

func NewRuntimePipeWriter(stateDir string) (RuntimePipeWriter, error) {
	p := new(RuntimePipeW)
	fd, err := newWriter(stateDir, runtimePipeFileName)
	if err != nil {
		return nil, err
	}
	p.fd = fd
	return p, nil
}

type RuntimePipeR struct {
	fd                 *os.File
	idMapper           procfs.IdMapper
	readyContext       context.Context
	finishReadyContext context.CancelFunc
}

func (r *RuntimePipeR) listen() (waitForReady <-chan struct{}) {
	r.readyContext, r.finishReadyContext = context.WithCancel(context.TODO())
	log := logger.Log().Named("pipe").Named("runtime")
	log.Debug("listening on runtime pipe")

	go func() {
		for msg := range readFromInit(context.TODO(), r.fd) {
			switch msg.Payload.(type) {
			case *pb.FromInit_Ready:
				log.Debug("received ready message")
				r.onReady()
				return
			case *pb.FromInit_MapUid:
				log.Debug("received map uid request")
				req := msg.GetMapUid()
				err := r.idMapper.MapUid(procfs.PidSelf, req.InsideId, req.OutsideId)
				if err != nil {
					log.Error("map uid failed", zap.Error(err))
				}
			case *pb.FromInit_MapGid:
				log.Debug("received map gid request")
				req := msg.GetMapGid()
				err := r.idMapper.MapGid(procfs.PidSelf, req.InsideId, req.OutsideId)
				if err != nil {
					log.Error("map gid failed", zap.Error(err))
				}
			}
		}
	}()

	return r.readyContext.Done()
}

func (r *RuntimePipeR) onReady() {
	r.finishReadyContext()
}

type RuntimePipeW struct {
	fd *os.File
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
