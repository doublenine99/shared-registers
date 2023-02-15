package main

import (
	"context"
	"log"
	"shared-registers/common/proto"
)

type server struct {
	proto.UnimplementedSharedRegistersServer
}

func (s *server) GetPhase(ctx context.Context, in *proto.GetPhaseReq) (*proto.GetPhaseRsp, error) {
	// TODO: GetLogic
	log.Printf("GetPhase Received: %v", in)
	return &proto.GetPhaseRsp{
		Val: "1",
		Ts:  nil,
	}, nil
}
func (s *server) SetPhase(ctx context.Context, in *proto.SetPhaseReq) (*proto.SetPhaseRsp, error) {
	// TODO: SetLogic
	log.Printf("SetPhase Received: %v", in)
	return &proto.SetPhaseRsp{}, nil
}
