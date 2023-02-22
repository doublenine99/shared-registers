package protocol

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"shared-registers/client/util"
	"shared-registers/common/proto"
	"time"
)

type grpcClient struct {
	conn                               *grpc.ClientConn
	c                                  proto.SharedRegistersClient
	requestTimeOut                     time.Duration
	DebugMode                          bool
	SetPhaseMockFail, GetPhaseMockFail bool
}

func (g *grpcClient) SetPhase(req *proto.SetPhaseReq) error {
	if g.DebugMode {
		defer util.PrintFuncExeTime("SetPhase", time.Now())
	}
	if g.SetPhaseMockFail {
		return errors.New(fmt.Sprintf("%s GetPhase failed: MockError", g.conn.Target()))
	}
	ctx, cancel := context.WithTimeout(context.Background(), g.requestTimeOut)
	defer cancel()
	_, err := g.c.SetPhase(ctx, req)
	if err != nil {
		log.Printf("%s SetPhase failed: %v", g.conn.Target(), err)
		return err
	}

	return nil
}

func (g *grpcClient) GetPhase(req *proto.GetPhaseReq) (*proto.GetPhaseRsp, error) {
	if g.DebugMode {
		defer util.PrintFuncExeTime("GetPhase", time.Now())
	}
	if g.GetPhaseMockFail {
		return nil, errors.New(fmt.Sprintf("%s GetPhase failed: MockError", g.conn.Target()))
	}
	ctx, cancel := context.WithTimeout(context.Background(), g.requestTimeOut)
	defer cancel()
	rsp, err := g.c.GetPhase(ctx, req)
	if err != nil {
		log.Printf("%s GetPhase failed: %v", g.conn.Target(), err)
		return nil, err
	}

	return rsp, nil
}
