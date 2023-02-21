package test

import (
	"fmt"
	"log"
	"os"
	"shared-registers/client/protocol"
	"strconv"
	"sync"
	"testing"
	"time"
)

var (
	serverAddrs []string
)

func setup() {
	fmt.Println("Before all tests")
	serverAddrs = []string{"localhost:50051"} // set servers' addresses
	setUpClient, _ := protocol.CreateSharedRegisterClient("setUpClient", serverAddrs)
	for i := 0; i < 100000; i++ {
		err := setUpClient.Write("k"+strconv.Itoa(i), "v"+strconv.Itoa(i))
		if err != nil {
			log.Fatalln("failed to initialize the k-v store with 100K k-v pairs.")
		}
		if i%10000 == 0 {
			log.Println(strconv.Itoa(i) + "'s write")
		}
	}
	fmt.Println("stored 100K k-v pairs")
}

func testReadOnly(numClients int, commandsNumPerClient int, t *testing.T) (throughPutPerClient float64) {
	var wg sync.WaitGroup
	wg.Add(numClients)
	startTime := time.Now()
	for clientId := 0; clientId < numClients; clientId++ {
		go func(clientId int) {
			client, _ := protocol.CreateSharedRegisterClient("clientRead"+strconv.Itoa(clientId), serverAddrs)
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

func testWriteOnly(numClients int, commandsNumPerClient int, t *testing.T) (throughPutPerClient float64) {
	var wg sync.WaitGroup
	wg.Add(numClients)
	startTime := time.Now()
	for clientId := 0; clientId < numClients; clientId++ {
		go func(clientId int) {
			client, _ := protocol.CreateSharedRegisterClient("clientWrite"+strconv.Itoa(clientId), serverAddrs)
			for i := 0; i < commandsNumPerClient; i++ {
				key, value := "k"+strconv.Itoa(i), "v"+strconv.Itoa(i)
				err := client.Write(key, value)
				if err != nil {
					t.Errorf("Incorrect write: key=%s", key)
				}
			}
			wg.Done()
		}(clientId)
	}
	wg.Wait()
	throughPutPerClient = float64(commandsNumPerClient) / (float64(time.Since(startTime)) / float64(time.Second))
	return throughPutPerClient
}

func testReadAndWrite(numClients int, commandsNumPerClient int, t *testing.T) (throughPutPerClient float64) {
	var wg sync.WaitGroup
	wg.Add(numClients)
	startTime := time.Now()
	for clientId := 0; clientId < numClients; clientId++ {
		go func(clientId int) {
			client, _ := protocol.CreateSharedRegisterClient("clientReadAndWrite"+strconv.Itoa(clientId), serverAddrs)
			for i := 0; i < commandsNumPerClient/2; i++ {
				key, value := "k"+strconv.Itoa(i), "v"+strconv.Itoa(i)
				err := client.Write(key, value)
				if err != nil {
					t.Errorf("Incorrect read: key=%s", key)
				}
				_, err = client.Read(key)
			}
			wg.Done()
		}(clientId)
	}
	wg.Wait()
	throughPutPerClient = float64(commandsNumPerClient) / (float64(time.Since(startTime)) / float64(time.Second))
	return throughPutPerClient
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func TestRunClient(t *testing.T) {
	commandNumPerClient := 10000
	for numClients := 1; numClients <= 32; numClients *= 2 {
		throughPutPerClient := testReadOnly(numClients, commandNumPerClient, t)
		fmt.Printf("Read Only, numClient=%d, throughputPerSecondPerClient=%f\n", numClients, throughPutPerClient)
		throughPutPerClient = testWriteOnly(numClients, commandNumPerClient, t)
		fmt.Printf("Write Only, numClient=%d, throughputPerSecondPerClient=%f\n", numClients, throughPutPerClient)
		throughPutPerClient = testReadAndWrite(numClients, commandNumPerClient, t)
		fmt.Printf("Read And Write, numClient=%d, throughputPerSecondPerClient=%f\n", numClients, throughPutPerClient)
	}
}
