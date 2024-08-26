package ipc

import (
	"google.golang.org/protobuf/proto"
	pb "roci/proto"
)

// NewMessageReady creates new 'ready' proto message
func NewMessageReady() proto.Message {
	return &pb.FromInit{Payload: &pb.FromInit_Ready{Ready: &pb.Ready{}}}
}

// NewMessageStart creates new 'start' proto message
func NewMessageStart() proto.Message {
	return &pb.FromRuntime{Payload: &pb.FromRuntime_Start{Start: &pb.Start{}}}
}

// NewMessageUidMapping creates new UidMapping call message
func NewMessageUidMapping(insideId, outsideId uint32) proto.Message {
	return &pb.FromInit{Payload: &pb.FromInit_MapUid{MapUid: &pb.IdMapping{
		InsideId:  insideId,
		OutsideId: outsideId,
	}}}
}

// NewMessageGidMapping creates new GidMapping call message
func NewMessageGidMapping(insideId, outsideId uint32) proto.Message {
	return &pb.FromInit{Payload: &pb.FromInit_MapGid{MapGid: &pb.IdMapping{
		InsideId:  insideId,
		OutsideId: outsideId,
	}}}
}
