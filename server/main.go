package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"shared-registers/common/proto"
)

var (
	port = 50001
)

func parseArgs() {
	defer flag.Parse()
	port = *flag.Int("port", 50051, "the port to start the service")
	flag.Parse()
}
func main() {
	parseArgs()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	proto.RegisterSharedRegistersServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
