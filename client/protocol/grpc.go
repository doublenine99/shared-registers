package protocol

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"shared-registers/client/util"
	"shared-registers/common/proto"
	"time"
)

type grpcClient struct {
	conn                                             *grpc.ClientConn
	c                                                proto.SharedRegistersClient
	requestTimeOut                                   time.Duration
	DebugMode                                        bool
	SetPhaseMockFail, GetPhaseMockFail, RespMockFail bool
}

func createGrpcClient(addr string) (*grpcClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil || conn == nil {
		log.Printf("did not connect to %s: %v", addr, err)
		return nil, err
	}

	return &grpcClient{
		conn:           conn,
		c:              proto.NewSharedRegistersClient(conn),
		requestTimeOut: 500 * time.Millisecond,
	}, nil
}

func (g *grpcClient) SetPhase(req *proto.SetPhaseReq) error {
	if g.DebugMode {
		defer util.PrintFuncExeTime("SetPhase", time.Now())
	}
	if g.SetPhaseMockFail {
		log.Printf("%s SetPhaseMockFail\n", g.conn.Target())
		return errors.New(fmt.Sprintf("%s GetPhase failed: MockError", g.conn.Target()))
	}
	ctx, cancel := context.WithTimeout(context.Background(), g.requestTimeOut)
	defer cancel()
	_, err := g.c.SetPhase(ctx, req)
	if err != nil {
		//log.Printf("%s SetPhase failed: %v", g.conn.Target(), err)
		return err
	}
	if g.RespMockFail {
		log.Printf("%s RespMockFail\n", g.conn.Target())
		return errors.New(fmt.Sprintf("%s GetPhase failed: MockError", g.conn.Target()))
	}

	return nil
}

func (g *grpcClient) GetPhase(req *proto.GetPhaseReq) (*proto.GetPhaseRsp, error) {
	if g.DebugMode {
		defer util.PrintFuncExeTime("GetPhase", time.Now())
	}
	if g.GetPhaseMockFail {
		log.Printf("%s GetPhaseMockFail\n", g.conn.Target())
		return nil, errors.New(fmt.Sprintf("%s GetPhase failed: MockError", g.conn.Target()))
	}
	ctx, cancel := context.WithTimeout(context.Background(), g.requestTimeOut)
	defer cancel()
	rsp, err := g.c.GetPhase(ctx, req)
	if err != nil {
		//log.Printf("%s GetPhase failed: %v", g.conn.Target(), err)
		return nil, err
	}
	if g.RespMockFail {
		log.Printf("%s RespMockFail\n", g.conn.Target())
		return nil, errors.New(fmt.Sprintf("%s GetPhase failed: MockError", g.conn.Target()))
	}

	return rsp, nil
}
