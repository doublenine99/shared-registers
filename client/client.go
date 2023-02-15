package main

import (
	"context"
	"flag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"shared-registers/common/proto"
	"time"
)

var (
	client *grpcClient // client instance
)

func init() {
	flag.Parse()
	addr := flag.String("addr", "localhost:50051", "the address to connect to")
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	client = &grpcClient{
		conn:           conn,
		c:              proto.NewSharedRegistersClient(conn),
		requestTimeOut: time.Second,
	}
}

type grpcClient struct {
	conn           *grpc.ClientConn
	c              proto.SharedRegistersClient
	requestTimeOut time.Duration
}

func SetPhase(req *proto.SetPhaseReq) error {
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

func GetPhase(req *proto.GetPhaseReq) (*proto.GetPhaseRsp, error) {
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

func main() {
	_ = SetPhase(&proto.SetPhaseReq{
		Key:   "1",
		Value: "1",
		Ts:    nil,
	})
	_, _ = GetPhase(&proto.GetPhaseReq{Key: "1"})

}
