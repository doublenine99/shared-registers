package main

import (
	"context"
	"log"
	"shared-registers/common/proto"
	"shared-registers/server/store"
)

type server struct {
	proto.UnimplementedSharedRegistersServer
}

func (s *server) GetPhase(ctx context.Context, in *proto.GetPhaseReq) (*proto.GetPhaseRsp, error) {
	log.Printf("GetPhase Received: %v", in)

	val, ts, success := store.Get(in.GetKey())
	if success == true {
		return &proto.GetPhaseRsp{
			Val: val,
			Ts:  ts,
		}, nil
	}
	return &proto.GetPhaseRsp{
		Val: "1",
		Ts:  nil,
	}, nil // TODO: change error from nil to meaningful
}
func (s *server) SetPhase(ctx context.Context, in *proto.SetPhaseReq) (*proto.SetPhaseRsp, error) {
	log.Printf("SetPhase Received: %v", in)

	// Make sure the current value in the map has an older timestamp
	_, ts, success := store.Get(in.GetKey())
	if success == true && ts.RequestNumber > in.GetTs().RequestNumber {
		return &proto.SetPhaseRsp{}, nil
	}

	store.Set(in.GetKey(), in.GetValue(), in.GetTs())
	return &proto.SetPhaseRsp{}, nil
}
