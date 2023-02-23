package protocol

import (
	"strconv"
	"testing"
)

var (
	_testServiceAddrs = []string{
		"amd183.utah.cloudlab.us:50051",
		"amd185.utah.cloudlab.us:50051",
		"amd192.utah.cloudlab.us:50051",
		"amd200.utah.cloudlab.us:50051",
		"amd204.utah.cloudlab.us:50051",
	}
)

func TestWriteFailBeforeGetPhase(t *testing.T) {
	commandNum := 10
	testClient, err := CreateSharedRegisterClient("testClient", _testServiceAddrs)
	if err != nil {
		t.Error(err)
	}
	if len(testClient.replicaConns) != 5 {
		t.Error("fail to connect to 5 replicas")
	}
	// simulate less than quorumSize replicas fail to process request before GetPhase
	for i := 0; i < testClient.quorumSize-1; i++ {
		testClient.replicaConns[i].GetPhaseMockFail = true
		testClient.replicaConns[i].SetPhaseMockFail = true
	}

	for i := 0; i < commandNum; i++ {
		key, value := "CK"+strconv.Itoa(i), "CV"+strconv.Itoa(i)
		err := testClient.Write(key, value)
		if err != nil {
			t.Errorf("Failed write: key=%s", key)
		}
		result, err := testClient.Read(key)
		if err == nil && result != value {
			t.Errorf("Incorrect read: key=%s, actualValue=%s, expectedValue=%s", key, result, value)
		}
	}
}

func TestWriteFailBetweenTwoPhases(t *testing.T) {
	commandNum := 10
	testClient, err := CreateSharedRegisterClient("testClient", _testServiceAddrs)
	if err != nil {
		t.Error(err)
	}
	if len(testClient.replicaConns) != 5 {
		t.Error("fail to connect to 5 replicas")
	}
	// simulate less than quorumSize replicas fail to process request after GetPhase but before SetPhase
	for i := 0; i < testClient.quorumSize-1; i++ {
		testClient.replicaConns[i].GetPhaseMockFail = false
		testClient.replicaConns[i].SetPhaseMockFail = true
	}

	for i := 0; i < commandNum; i++ {
		key, value := "DK"+strconv.Itoa(i), "DV"+strconv.Itoa(i)
		err := testClient.Write(key, value)
		if err != nil {
			t.Errorf("Failed write: key=%s", key)
		}
		result, err := testClient.Read(key)
		if err == nil && result != value {
			t.Errorf("Incorrect read: key=%s, actualValue=%s, expectedValue=%s", key, result, value)
		}
	}
}

func TestWriteAfterTwoPhases(t *testing.T) {
	commandNum := 10
	testClient, err := CreateSharedRegisterClient("testClient", _testServiceAddrs)
	if err != nil {
		t.Error(err)
	}
	if len(testClient.replicaConns) != 5 {
		t.Error("fail to connect to 5 replicas")
	}
	// use microsecond timeout to simulate the error could not receive the ack from replica
	for i := 0; i < testClient.quorumSize-1; i++ {
		testClient.replicaConns[i].RespMockFail = true
	}

	for i := 0; i < commandNum; i++ {
		key, value := "EK"+strconv.Itoa(i), "EV"+strconv.Itoa(i)
		err := testClient.Write(key, value)
		if err != nil {
			t.Errorf("Failed write: key=%s", key)
		}
		result, err := testClient.Read(key)
		if err == nil && result != value {
			t.Errorf("Incorrect read: key=%s, actualValue=%s, expectedValue=%s", key, result, value)
		}
	}
}
