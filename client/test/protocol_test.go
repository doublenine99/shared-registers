package test

import (
	"log"
	"shared-registers/client/protocol"
	"shared-registers/client/util"
	"strconv"
	"sync"
	"testing"
	"time"
)

var (
	totalCommandNum = 20000
	maxClientNum    = 32
	serverAddrs     = []string{
		"amd183.utah.cloudlab.us:50051",
		"amd185.utah.cloudlab.us:50051",
		"amd192.utah.cloudlab.us:50051",
		"amd200.utah.cloudlab.us:50051",
		"amd204.utah.cloudlab.us:50051",
	} // set servers' addresses
)

func initKVStore(initNum int) {
	log.Println("Start initKVStore")
	setUpClient, err := protocol.CreateSharedRegisterClient("setUpClient", serverAddrs)
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < initNum; i++ {
		err := setUpClient.Write("k"+strconv.Itoa(i), "v"+strconv.Itoa(i))
		if err != nil {
			log.Fatalf("failed to initialize the k-v store with %d k-v pairs.", initNum)
		}

		if i%(initNum/10) == 0 {
			log.Println(strconv.Itoa(i) + "'s write")
		}
	}
	log.Printf("stored %d k-v pairs\n", initNum)
}

// * detect racing and generating profile for debugging purpose
// go test protocol_test.go -v -race -bench=. -benchmem -memprofile memprofile.out -cpuprofile profile.out &> out_prof.log
// * run the benchmark with verbose log
// go test protocol_test.go -v -bench=. &> out.log
func BenchmarkRunClient(b *testing.B) {
	go util.PrintGoroutineNum(3 * time.Second)
	//initKVStore(1000000)
	for numClients := 1; numClients <= maxClientNum; numClients *= 2 {
		throughPutPerClient := testReadOnly(numClients, totalCommandNum/numClients, b)
		b.Logf("Read Only\t numClient=%d\t throughputPerSecondPerClient=%f\t totalThroughput=%f\n", numClients, throughPutPerClient, float64(numClients)*throughPutPerClient)
		throughPutPerClient = testWriteOnly(numClients, totalCommandNum/numClients, b)
		b.Logf("Write Only\t numClient=%d\t throughputPerSecondPerClient=%f\t totalThroughput=%f\n", numClients, throughPutPerClient, float64(numClients)*throughPutPerClient)
		throughPutPerClient = testReadAndWrite(numClients, totalCommandNum/numClients, b)
		b.Logf("R And W\t numClient=%d\t throughputPerSecondPerClient=%f\t totalThroughput=%f\n", numClients, throughPutPerClient, float64(numClients)*throughPutPerClient)
	}
}

func testReadOnly(numClients int, commandsNumPerClient int, t *testing.B) (throughPutPerClient float64) {
	var wg sync.WaitGroup
	wg.Add(numClients)
	startTime := time.Now()
	for clientId := 0; clientId < numClients; clientId++ {
		go func(clientId int) {
			client, err := protocol.CreateSharedRegisterClient("clientRead"+strconv.Itoa(clientId), serverAddrs)
			if err != nil {
				log.Fatalf("CreateSharedRegisterClient err: %v %d", err, clientId)
			}

			for i := 0; i < commandsNumPerClient; i++ {
				key, expectedValue := "k"+strconv.Itoa(i), "v"+strconv.Itoa(i)
				result, err := client.Read(key)
				if err == nil && result != expectedValue {
					t.Errorf("Incorrect read: key=%s, actualValue=%s, expectedValue=%s", key, result, expectedValue)
				}
			}
			wg.Done()
		}(clientId)
	}
	wg.Wait()
	throughPutPerClient = float64(commandsNumPerClient) / (float64(time.Since(startTime)) / float64(time.Second))
	return throughPutPerClient
}

func testWriteOnly(numClients int, commandsNumPerClient int, t *testing.B) (throughPutPerClient float64) {
	var wg sync.WaitGroup
	wg.Add(numClients)
	startTime := time.Now()
	for clientId := 0; clientId < numClients; clientId++ {
		go func(clientId int) {
			client, err := protocol.CreateSharedRegisterClient("clientWrite"+strconv.Itoa(clientId), serverAddrs)
			if err != nil {
				log.Fatalf("CreateSharedRegisterClient err: %v %d", err, clientId)
			}
			for i := 0; i < commandsNumPerClient; i++ {
				key, value := "k"+strconv.Itoa(i), "v"+strconv.Itoa(i)
				err := client.Write(key, value)
				if err != nil {
					t.Errorf("Failed write: key=%s", key)
				}
			}
			wg.Done()
		}(clientId)
	}
	wg.Wait()
	throughPutPerClient = float64(commandsNumPerClient) / (float64(time.Since(startTime)) / float64(time.Second))
	return throughPutPerClient
}

func testReadAndWrite(numClients int, commandsNumPerClient int, t *testing.B) (throughPutPerClient float64) {
	var wg sync.WaitGroup
	wg.Add(numClients)
	startTime := time.Now()
	for clientId := 0; clientId < numClients; clientId++ {
		go func(clientId int) {
			client, err := protocol.CreateSharedRegisterClient("clientReadAndWrite"+strconv.Itoa(clientId), serverAddrs)
			if err != nil {
				log.Fatalf("CreateSharedRegisterClient err: %v %d", err, clientId)
			}
			for i := 0; i < commandsNumPerClient/2; i++ {
				key, value := "k"+strconv.Itoa(i), "v"+strconv.Itoa(i)
				err := client.Write(key, value)
				if err != nil {
					t.Errorf("Failed write: key=%s", key)
				}
				result, err := client.Read(key)
				if err == nil && result != value {
					t.Errorf("Incorrect read: key=%s, actualValue=%s, expectedValue=%s", key, result, value)
				}
			}
			wg.Done()
		}(clientId)
	}
	wg.Wait()
	throughPutPerClient = float64(commandsNumPerClient) / (float64(time.Since(startTime)) / float64(time.Second))
	return throughPutPerClient
}
