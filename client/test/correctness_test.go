package test

import (
	"shared-registers/client/protocol"
	"testing"
)

func Test(t *testing.T) {
	// TODO: add mock fail grpc client
	protocol.CreateClientWithCustomizedGrpcClient("1")
}
