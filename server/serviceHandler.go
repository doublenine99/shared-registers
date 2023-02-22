package main

import (
	"context"
	"log"
	"shared-registers/common"
	"shared-registers/common/proto"
	"shared-registers/server/store"
)

type server struct {
	proto.UnimplementedSharedRegistersServer
}

// GetPhase
// return to client with the <value, ts> associated with the key to let client decide whether which
// replica has the updated value
func (s *server) GetPhase(ctx context.Context, in *proto.GetPhaseReq) (*proto.GetPhaseRsp, error) {
	//log.Printf("GetPhase Received: %v", in)
	v, err := store.Get(in.GetKey())
	if err != nil {
		log.Printf("GetPhase err: %v", err)
		return nil, err
	}
	return &proto.GetPhaseRsp{Value: v}, nil
}

// SetPhase
// Each replica checks if this ts-new is larger than the one it stores
// If yes, replica stores v, ts-new.
// In either case, the storage nodes sends an acknowledgement to the client.
func (s *server) SetPhase(ctx context.Context, in *proto.SetPhaseReq) (*proto.SetPhaseRsp, error) {
	//log.Printf("SetPhase Received: %v", in)
	newTs := in.GetValue().GetTs()
	currValue, err := store.Get(in.GetKey())
	if err != nil {
		return nil, err
	}
	if currValue == nil || common.FindLargestTimeStamp(currValue.Ts, newTs) == newTs {
		store.Set(in.GetKey(), in.GetValue())
	}
	return &proto.SetPhaseRsp{}, nil
}
