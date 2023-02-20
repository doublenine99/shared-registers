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

func SetPhase(req *proto.SetPhaseReq, client *grpcClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), client.requestTimeOut)
	defer cancel()
	_, err := client.c.SetPhase(ctx, req)
	if err != nil {
		log.Printf("SetPhase failed: %v", err)
		return err
	}
	log.Printf("SetPhase succ")
	return nil
}

func GetPhase(req *proto.GetPhaseReq, client *grpcClient) (*proto.GetPhaseRsp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), client.requestTimeOut)
	defer cancel()
	rsp, err := client.c.GetPhase(ctx, req)
	if err != nil {
		log.Printf("GetPhase failed: %v", err)
		return nil, err
	}
	log.Printf("GetPhase succ: %s", rsp.String())
	return rsp, nil
}
