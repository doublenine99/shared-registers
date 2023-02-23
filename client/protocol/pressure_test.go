package protocol

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

var numKeys = 100000

func generateRandomIntString() string {
	rand.Seed(time.Now().UnixNano())
	min := 1
	max := numKeys
	return strconv.Itoa(rand.Intn(max-min+1) + min)
}

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

func TestInitKVStore(t *testing.T) {
	initKVStore(10)
}

// * detect racing and generating profile for debugging purpose
// go test -run BenchmarkRunClient -v -race -bench=. -benchmem -memprofile memprofile.out -cpuprofile profile.out &> out_prof.log
// * run the benchmark with verbose log
// go test -run BenchmarkRunClient -v -bench=. &> out.log
func BenchmarkRunClient(b *testing.B) {
	_testMaxClientNum := 32
	//go util.PrintGoroutineNum(3 * time.Second)

	dirName := "results" // store results in this dir
	// Check if the directory already exists
	if _, err := os.Stat(dirName); err == nil {
		// If the directory exists, remove it and its contents
		err := os.RemoveAll(dirName)
		if err != nil {
			fmt.Printf("Error removing directory: %v\n", err)
			return
		}
	}
	// Create the directory
	err := os.Mkdir(dirName, os.ModePerm)
	if err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return
	}

	for numClients := 1; numClients <= _testMaxClientNum; numClients *= 2 {
		throughPutPerSec := testReadOnly(numClients, b)
		b.Logf("Read Only\t numClient=%d\t totalThroughput=%f\n", numClients, throughPutPerSec)
		throughPutPerSec = testWriteOnly(numClients, b)
		b.Logf("Write Only\t numClient=%d\t totalThroughput=%f\n", numClients, throughPutPerSec)
		throughPutPerSec = testReadAndWrite(numClients, b)
		b.Logf("R And W\t numClient=%d\t totalThroughput=%f\n", numClients, throughPutPerSec)
	}
}

func testReadOnly(numClients int, t *testing.B) float64 {
	var wg sync.WaitGroup
	wg.Add(numClients)
	var totalCommandCount uint32 = 0
	startTime := time.Now()
	for clientId := 1; clientId <= numClients; clientId++ {
		go func(clientId int) {
			client, err := CreateSharedRegisterClient("clientRead"+strconv.Itoa(clientId), _testServiceAddrs)
			if err != nil {
				log.Fatalf("CreateSharedRegisterClient err: %v %d", err, clientId)
			}

			for start := time.Now(); time.Since(start) < time.Second*100; {
				randInt := generateRandomIntString()
				key, expectedValue := "k"+randInt, "v"+randInt
				result, err := client.Read(key)
				if err == nil && result != expectedValue {
					t.Errorf("Incorrect read: key=%s, actualValue=%s, expectedValue=%s", key, result, expectedValue)
				}
				atomic.AddUint32(&totalCommandCount, 1)
			}
			wg.Done()
		}(clientId)
	}
	wg.Wait()
	throughPutPerSec := float64(totalCommandCount) / (float64(time.Since(startTime)) / float64(time.Second))
	return throughPutPerSec
}

func testWriteOnly(numClients int, t *testing.B) float64 {
	var wg sync.WaitGroup
	var totalCommandCount uint32 = 0
	wg.Add(numClients)
	startTime := time.Now()
	for clientId := 1; clientId <= numClients; clientId++ {
		go func(clientId int) {

			client, err := CreateSharedRegisterClient("clientWrite"+strconv.Itoa(clientId), _testServiceAddrs)
			if err != nil {
				log.Fatalf("CreateSharedRegisterClient err: %v %d", err, clientId)
			}

			for start := time.Now(); time.Since(start) < time.Second*100; {
				randInt := generateRandomIntString()
				key, value := "k"+randInt, "v"+randInt
				err := client.Write(key, value)
				if err != nil {
					t.Errorf("Failed write: key=%s", key)
				}
				atomic.AddUint32(&totalCommandCount, 1)
			}
			wg.Done()
		}(clientId)
	}
	wg.Wait()
	throughPutPerSec := float64(totalCommandCount) / (float64(time.Since(startTime)) / float64(time.Second))
	return throughPutPerSec
}

func testReadAndWrite(numClients int, t *testing.B) float64 {
	var wg sync.WaitGroup
	var totalCommandCount uint32 = 0
	wg.Add(numClients)
	startTime := time.Now()
	for clientId := 1; clientId <= numClients; clientId++ {
		go func(clientId int) {
			resultFile, err := os.Create("results/" + strconv.Itoa(numClients) + "clients_" + strconv.Itoa(clientId) + ".txt")
			defer resultFile.Close()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			resultWriter := bufio.NewWriter(resultFile)
			defer resultWriter.Flush()

			client, err := CreateSharedRegisterClient("clientReadAndWrite"+strconv.Itoa(clientId), _testServiceAddrs)
			if err != nil {
				log.Fatalf("CreateSharedRegisterClient err: %v %d", err, clientId)
			}

			for start := time.Now(); time.Since(start) < time.Second*100; {
				randInt := generateRandomIntString()
				key, value := "k"+randInt, "v"+randInt
				err := client.Write(key, value)
				if err != nil {
					t.Errorf("Failed write: key=%s", key)
				}
				resultToLog := "W, " + time.Now().UTC().String() + "\n"
				_, err = resultWriter.WriteString(resultToLog)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				result, err := client.Read(key)
				if err == nil && result != value {
					t.Errorf("Incorrect read: key=%s, actualValue=%s, expectedValue=%s", key, result, value)
				}
				resultToLog = "R, " + time.Now().UTC().String() + "\n"
				_, err = resultWriter.WriteString(resultToLog)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				atomic.AddUint32(&totalCommandCount, 2)
			}

			wg.Done()
		}(clientId)
	}
	wg.Wait()
	throughPutPerSec := float64(totalCommandCount) / (float64(time.Since(startTime)) / float64(time.Second))
	return throughPutPerSec
}
