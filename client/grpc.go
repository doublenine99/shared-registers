package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"shared-registers/common/proto"
	"time"
)

type grpcClient struct {
	conn           *grpc.ClientConn
	c              proto.SharedRegistersClient
	requestTimeOut time.Duration
}

func (g *grpcClient) SetPhase(req *proto.SetPhaseReq) error {
	ctx, cancel := context.WithTimeout(context.Background(), g.requestTimeOut)
	defer cancel()
	_, err := g.c.SetPhase(ctx, req)
	if err != nil {
		log.Printf("SetPhase failed: %v", err)
		return err
	}
	log.Printf("SetPhase succ")
	return nil
}

func (g *grpcClient) GetPhase(req *proto.GetPhaseReq) (*proto.GetPhaseRsp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), g.requestTimeOut)
	defer cancel()
	rsp, err := g.c.GetPhase(ctx, req)
	if err != nil {
		log.Printf("GetPhase failed: %v", err)
		return nil, err
	}
	log.Printf("GetPhase succ: %s", rsp.String())
	return rsp, nil
}
