package main

import (
	"bufio"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	"shared-registers/common/proto"
	"strconv"
	"sync"
	"time"
)

var (
	//client         *grpcClient // client instance
	quorum       int
	numReplicas  int
	myClientId   string
	mutex        sync.Mutex
	replicaConns []*grpcClient
)

func init() {
	hostname, _ := os.Hostname()
	myClientId = hostname + strconv.Itoa(os.Getpid())

	file, err := os.Open("../common/config.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		addr := scanner.Text()
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		client := &grpcClient{
			conn:           conn,
			c:              proto.NewSharedRegistersClient(conn),
			requestTimeOut: time.Second,
		}
		replicaConns = append(replicaConns, client)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//addr := flag.String("addr", "localhost:50051", "the address to connect to")
	// Set up a connection to the server.
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second)

}

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

func compareTs(curMaxTs *uint64, curVal *string, curClientId *string, resp *proto.GetPhaseRsp) {
	val := resp.Val
	logicalTs := resp.Ts.RequestNumber
	clientId := resp.Ts.ClientID
	if *curMaxTs == 0 {
		*curVal = val
		*curMaxTs = logicalTs
		*curClientId = clientId
		return
	}

	if logicalTs > *curMaxTs {
		*curVal = val
		*curMaxTs = logicalTs
		*curClientId = clientId
	} else if *curMaxTs == logicalTs {
		if clientId > *curClientId {
			*curVal = val
			*curMaxTs = logicalTs
			*curClientId = clientId
		}
	}
}

func completeGetPhase(key string) (string, uint64, string) {
	var wg sync.WaitGroup
	curVal := ""
	var curMaxTs uint64 = 0
	var curClientId string = ""
	for _, conn := range replicaConns {
		wg.Add(1)
		go func() {
			resp, err := GetPhase(&proto.GetPhaseReq{Key: key}, conn)
			if resp != nil {
				mutex.Lock()
				if resp.Ts != nil {
					compareTs(&curMaxTs, &curVal, &curClientId, resp)
				}
				mutex.Unlock()
				wg.Done()
			} else {
				log.Printf("replica get failed: %v", err)
			}
		}()
	}
	wg.Wait()
	return curVal, curMaxTs, curClientId
}

func writeGetPhase(key string) (string, *proto.TimeStamp) {
	curVal, curMaxTs, _ := completeGetPhase(key)
	var ts = &proto.TimeStamp{ClientID: myClientId, RequestNumber: curMaxTs + 1}
	return curVal, ts
}

func readGetPhase(key string) (string, *proto.TimeStamp) {
	curVal, curMaxTs, curClientId := completeGetPhase(key)
	if curMaxTs == 0 {
		panic("key does not exist")
	} else {
		ts := &proto.TimeStamp{ClientID: curClientId, RequestNumber: curMaxTs}
		return curVal, ts
	}
}

func doWrite(key string, val string) {
	//val, ts := writeGetPhase(key)
	//TODO: implement set phase
	//_ = SetPhase(&proto.SetPhaseReq{
	//	Key:   key,
	//	Value: val,
	//	Ts:    ts,
	//})
}
func main() {

}
