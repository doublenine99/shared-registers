package protocol

import (
	"errors"
	"log"
	"shared-registers/client/util"
	"shared-registers/common"
	"shared-registers/common/proto"
	"sync"
	"time"
)

type SharedRegisterClient struct {
	ClientID     string
	PhaseTimeout time.Duration // the max waiting time from all the replicas each phase, default 1s
	replicaConns []*grpcClient
	quorumSize   int        // len(replicaConns) / 2 + 1
	opsLock      sync.Mutex // each SharedRegisterClient should only execute operations sequentially
	DebugMode    bool
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
		c, err := createGrpcClient(addr)
		if err != nil || c == nil {
			log.Printf("did not connect to %s: %v", addr, err)
			continue
		}
		s.replicaConns = append(s.replicaConns, c)
	}
	if len(s.replicaConns) < 3 {
		return nil, errors.New("have to connect to at least 3 replicas to continue")
	}
	//log.Printf("connected to %d servers\n", len(s.replicaConns))
	s.quorumSize = len(s.replicaConns)/2 + 1
	return s, nil
}

func (s *SharedRegisterClient) Write(key string, value string) error {
	s.opsLock.Lock()
	defer s.opsLock.Unlock()
	if s.DebugMode {
		defer util.PrintFuncExeTime("Write", time.Now())
	}

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
	s.opsLock.Lock()
	defer s.opsLock.Unlock()
	if s.DebugMode {
		defer util.PrintFuncExeTime("Read", time.Now())
	}

	latestValue, err := s.completeGetPhase(key)
	if latestValue == nil {
		return "", errors.New("key " + key + " doesn't exist")
	}
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
func (s *SharedRegisterClient) completeGetPhase(key string) (*proto.StoredValue, error) {
	// use a channel with size 1 to compare and store the value with the largest TS among concurrent
	// request goroutine to avoid data racing
	currMaxChan := make(chan *proto.StoredValue, 1)
	requests := make([]func() bool, 0)

	for _, conn := range s.replicaConns {
		conn := conn
		getFromReplica := func() bool {
			resp, err := conn.GetPhase(&proto.GetPhaseReq{Key: key})
			if err != nil {
				return false
			}
			if resp != nil && resp.GetValue() != nil {
				// read from the channel for current largest TS and compare with the current resp
				currLargest, open := <-currMaxChan
				// if the channel is already closed, ignore the response from the replica
				if !open {
					return true
				}
				if resp.GetValue().GetTs() == common.FindLargestTimeStamp(currLargest.GetTs(), resp.GetValue().GetTs()) {
					currMaxChan <- resp.GetValue()
				} else {
					currMaxChan <- currLargest
				}
			}
			return true
		}
		requests = append(requests, getFromReplica)
	}
	currMaxChan <- &proto.StoredValue{Ts: &proto.TimeStamp{}}
	timedOut := util.WaitForMajoritySuccessFromJobs(s.quorumSize, s.PhaseTimeout, requests)
	if timedOut {
		return nil, errors.New("completeGetPhase timeout")
	}
	largestVal := <-currMaxChan
	// have to manually the channel to let the unfinished request goroutine detect and return
	close(currMaxChan)
	return largestVal, nil
}

// client asks storage nodes to store the (v, ts-new).
// Each replica checks if this ts-new is larger than the one it stores
// If yes, replica stores v, ts-new.
// In either case, the storage nodes sends an acknowledgement to the client.
// client then waits for a majority of acknowledgements
func (s *SharedRegisterClient) completeSetPhase(key, value string, timestamp *proto.TimeStamp) error {
	requests := make([]func() bool, 0)
	for _, conn := range s.replicaConns {
		conn := conn
		setToReplica := func() bool {
			err := conn.SetPhase(&proto.SetPhaseReq{
				Key: key,
				Value: &proto.StoredValue{
					Val: value,
					Ts:  timestamp,
				}})
			if err != nil {
				return false
			}
			return true
		}
		requests = append(requests, setToReplica)
	}
	timedOut := util.WaitForMajoritySuccessFromJobs(s.quorumSize, s.PhaseTimeout, requests)
	if timedOut {
		return errors.New("completeGetPhase timeout")
	}
	return nil
}
