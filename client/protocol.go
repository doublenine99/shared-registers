package main

import (
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"shared-registers/common"
	"shared-registers/common/proto"
	"sync"
	"time"
)

type SharedRegisterClient struct {
	ClientID     string
	PhaseTimeout time.Duration // the max waiting time from all the replicas each phase, default 1s
	replicaConns []*grpcClient
	quorumSize   int // len(replicaConns) / 2 + 1
}

func CreateSharedRegisterClient(clientID string, serverAddrs []string) (*SharedRegisterClient, error) {
	// could add dedup logic in server as well
	if clientID == "" {
		return nil, errors.New("invalid client ID")
	}
	if len(serverAddrs) == 0 {
		return nil, errors.New("empty server addresses")
	}
	s := &SharedRegisterClient{
		ClientID:     clientID,
		PhaseTimeout: time.Second,
	}
	for _, addr := range serverAddrs {
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Printf("did not connect to %s: %v", addr, err)
			continue
		}
		s.replicaConns = append(s.replicaConns, &grpcClient{
			conn:           conn,
			c:              proto.NewSharedRegistersClient(conn),
			requestTimeOut: time.Second,
		})
	}

	s.quorumSize = len(s.replicaConns)/2 + 1
	return s, nil
}

func (s *SharedRegisterClient) Write(key string, value string) error {
	latestValue, err := s.completeGetPhase(key)
	if err != nil {
		return err
	}
	newTs := &proto.TimeStamp{
		RequestNumber: latestValue.GetTs().GetRequestNumber() + 1,
		ClientID:      s.ClientID,
	}
	return s.completeSetPhase(key, value, newTs)
}

func (s *SharedRegisterClient) Read(key string) (string, error) {
	latestValue, err := s.completeGetPhase(key)
	if err != nil {
		return "", err
	}
	err = s.completeSetPhase(key, latestValue.GetVal(), latestValue.GetTs())
	if err != nil {
		return "", err
	}
	return latestValue.GetVal(), nil
}

// client waits for a majority of responses from replicas for current <v, timestamp> pairs
// client finds largest received timestamp, and then chooses a higher unique timestamp ts-new (max-ts,client-id)
func (s *SharedRegisterClient) completeGetPhase(key string) (latestValue *proto.StoredValue, err error) {
	var majorityGroup sync.WaitGroup
	majorityGroup.Add(s.quorumSize) // wait for the responses from a majority

	replicaValues := make([]*proto.StoredValue, 0)
	replicaTimeStamps := make([]*proto.TimeStamp, 0)
	sliceLock := sync.Mutex{}
	// TODO: timeout if could not get majority responses
	for _, conn := range s.replicaConns {
		go func(conn *grpcClient) {
			resp, _ := GetPhase(&proto.GetPhaseReq{Key: key}, conn)
			if resp != nil {
				if resp.GetValue() != nil {
					sliceLock.Lock()
					replicaValues = append(replicaValues, resp.GetValue())
					replicaTimeStamps = append(replicaTimeStamps, resp.GetValue().GetTs())
					sliceLock.Unlock()
				}
				majorityGroup.Done()
			}
		}(conn)
	}
	// wait until majority or timeout
	majorityGroup.Wait()

	largestTs := common.FindLargestTimeStamp(replicaTimeStamps...)
	for _, v := range replicaValues {
		if v.GetTs() == largestTs {
			return v, nil
		}
	}
	return nil, nil
}

// client asks storage nodes to store the (v, ts-new).
// Each replica checks if this ts-new is larger than the one it stores
// If yes, replica stores v, ts-new.
// In either case, the storage nodes sends an acknowledgement to the client.
// client then waits for a majority of acknowledgements
func (s *SharedRegisterClient) completeSetPhase(key, value string, timestamp *proto.TimeStamp) error {
	var wg sync.WaitGroup
	// TODO: timeout
	wg.Add(s.quorumSize) // wait for the responses from a majority
	for _, conn := range s.replicaConns {
		go func(conn *grpcClient) {
			err := SetPhase(&proto.SetPhaseReq{
				Key: key,
				Value: &proto.StoredValue{
					Val: value,
					Ts:  timestamp,
				}}, conn)
			if err == nil {
				// successfully got a response form a replica
				wg.Done()
			}
		}(conn)
	}
	wg.Wait()
	return nil
}
