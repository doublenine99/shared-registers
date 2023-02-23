package protocol

import (
	"log"
	"shared-registers/client/util"
	"strconv"
	"sync"
	"testing"
	"time"
)

func initKVStore(initNum int) {
	log.Println("Start initKVStore")
	setUpClient, err := CreateSharedRegisterClient("setUpClient", _testServiceAddrs)
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
// go benchmark pressure_test.go -v -race -bench=. -benchmem -memprofile memprofile.out -cpuprofile profile.out &> out_prof.log
// * run the benchmark with verbose log
// go benchmark pressure_test.go -v -bench=. &> out.log
func BenchmarkRunClient(b *testing.B) {
	_totalCommandNum := 20000
	_testMaxClientNum := 32
	go util.PrintGoroutineNum(3 * time.Second)
	//initKVStore(1000000)
	for numClients := 1; numClients <= _testMaxClientNum; numClients *= 2 {
		throughPutPerClient := testReadOnly(numClients, _totalCommandNum/numClients, b)
		b.Logf("Read Only\t numClient=%d\t throughputPerSecondPerClient=%f\t totalThroughput=%f\n", numClients, throughPutPerClient, float64(numClients)*throughPutPerClient)
		throughPutPerClient = testWriteOnly(numClients, _totalCommandNum/numClients, b)
		b.Logf("Write Only\t numClient=%d\t throughputPerSecondPerClient=%f\t totalThroughput=%f\n", numClients, throughPutPerClient, float64(numClients)*throughPutPerClient)
		throughPutPerClient = testReadAndWrite(numClients, _totalCommandNum/numClients, b)
		b.Logf("R And W\t numClient=%d\t throughputPerSecondPerClient=%f\t totalThroughput=%f\n", numClients, throughPutPerClient, float64(numClients)*throughPutPerClient)
	}
}

func testReadOnly(numClients int, commandsNumPerClient int, t *testing.B) (throughPutPerClient float64) {
	var wg sync.WaitGroup
	wg.Add(numClients)
	startTime := time.Now()
	for clientId := 0; clientId < numClients; clientId++ {
		go func(clientId int) {
			client, err := CreateSharedRegisterClient("clientRead"+strconv.Itoa(clientId), _testServiceAddrs)
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
			client, err := CreateSharedRegisterClient("clientWrite"+strconv.Itoa(clientId), _testServiceAddrs)
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
			client, err := CreateSharedRegisterClient("clientReadAndWrite"+strconv.Itoa(clientId), _testServiceAddrs)
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
