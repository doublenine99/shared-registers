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
	"strings"
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
	wg.Add(len(replicaConns)/2 + 1) // wait for the responses from a majority
	for _, conn := range replicaConns {
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
				log.Printf("replica GET failed: %v", err)
			}
		}()
	}
	wg.Wait()
	return curVal, curMaxTs, curClientId
}

func writeGetPhase(key string) *proto.TimeStamp {
	_, curMaxTs, _ := completeGetPhase(key)
	var ts = &proto.TimeStamp{ClientID: myClientId, RequestNumber: curMaxTs + 1}
	return ts
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

func completeSetPhase(key, value string, timestamp *proto.TimeStamp) {
	var wg sync.WaitGroup
	wg.Add(len(replicaConns)/2 + 1) // wait for the responses from a majority
	for _, conn := range replicaConns {
		go func() {
			err := SetPhase(&proto.SetPhaseReq{Key: key, Value: value, Ts: timestamp}, conn)
			if err != nil {
				// successfully got a response form a replica
				wg.Done()
			} else {
				log.Printf("replica SET failed: %v", err)
			}
		}()
	}
	wg.Wait()
}

func writeSetPhase(key, value string, timestamp *proto.TimeStamp) {
	completeSetPhase(key, value, timestamp)
}

func readSetPhase(key, value string, timestamp *proto.TimeStamp) {
	completeSetPhase(key, value, timestamp)
}

func execWrite(key string, value string) {
	ts := writeGetPhase(key)
	writeSetPhase(key, value, ts)
}

func execRead(key string) string {
	value, ts := readGetPhase(key)
	readSetPhase(key, value, ts)
	return value
}

func execBatchOperations(fileName, resultFilename string) {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	resulFile, err := os.Create(resultFilename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer resulFile.Close()
	// Create a writer
	resultWriter := bufio.NewWriter(resulFile)

	scanner := bufio.NewScanner(file)
	// for each operation in the batch file
	for scanner.Scan() {
		operationFileds := strings.Fields(scanner.Text())
		var result = ""
		switch len(operationFileds) {
		case 2:
			if strings.EqualFold(operationFileds[0], "GET") {
				key := operationFileds[1]
				resultValue := execRead(key)
				result = "READ\tKey=" + key + "\tValue=" + resultValue
			} else {
				fmt.Println(err)
				os.Exit(1)
			}
		case 3:
			if strings.EqualFold(operationFileds[0], "SET") {
				key, value := operationFileds[1], operationFileds[2]
				execWrite(key, value)
				result = "WRITE\tKey=" + key + "\tValue=" + value
			} else {
				fmt.Println(err)
				os.Exit(1)
			}
		default:
			{
				fmt.Println(err)
				os.Exit(1)
			}
		}

		// write result into file
		_, err := resultWriter.WriteString(result + "\n")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if err := scanner.Err(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	resultWriter.Flush()
}

func main() {

}
