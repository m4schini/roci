package ipc

import (
	"google.golang.org/protobuf/proto"
	pb "roci/proto"
)

func NewMessageReady() proto.Message {
	return &pb.FromInit{Payload: &pb.FromInit_Ready{Ready: &pb.Ready{}}}
}

func NewMessageStart() proto.Message {
	return &pb.FromRuntime{Payload: &pb.FromRuntime_Start{Start: &pb.Start{}}}
}

func NewMessageUidMapping(insideId, outsideId uint32) proto.Message {
	return &pb.FromInit{Payload: &pb.FromInit_MapUid{MapUid: &pb.IdMapping{
		InsideId:  insideId,
		OutsideId: outsideId,
	}}}
}

func NewMessageGidMapping(insideId, outsideId uint32) proto.Message {
	return &pb.FromInit{Payload: &pb.FromInit_MapGid{MapGid: &pb.IdMapping{
		InsideId:  insideId,
		OutsideId: outsideId,
	}}}
}
