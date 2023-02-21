package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"shared-registers/client/protocol"
	"strconv"
	"strings"
)

var client *protocol.SharedRegisterClient

func init() {
	hostname, _ := os.Hostname()
	myClientId := hostname + strconv.Itoa(os.Getpid())

	file, err := os.Open("./config.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	addrs := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		addr := scanner.Text()
		addrs = append(addrs, addr)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	client, err = protocol.CreateSharedRegisterClient(myClientId, addrs)
	if err != nil {
		log.Fatal("CreateSharedRegisterClient: ", err)
	}
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
	lineNumber := 0
	// for each operation in the batch file
	for scanner.Scan() {
		inputLine := scanner.Text()
		operationFileds := strings.Fields(inputLine)
		var result = ""
		switch len(operationFileds) {
		case 2:
			if strings.EqualFold(operationFileds[0], "R") {
				key := operationFileds[1]
				resultValue, err := client.Read(key)
				if err != nil {
					result = err.Error()
				} else {
					result = "READ\tKey=" + key + "\tValue=" + resultValue
				}
			} else {
				fmt.Printf("Error when parsing line %d\n", lineNumber)
				fmt.Println(inputLine)
				os.Exit(1)
			}
		case 3:
			if strings.EqualFold(operationFileds[0], "W") {
				key, value := operationFileds[1], operationFileds[2]
				err = client.Write(key, value)
				if err != nil {
					log.Fatal("Client Write err: ", err)
				}
				result = "WRITE\tKey=" + key + "\tValue=" + value
			} else {
				fmt.Printf("Error when parsing line %d\n", lineNumber)
				fmt.Println(inputLine)
				os.Exit(1)
			}
		default:
			{
				fmt.Printf("Error when parsing line %d\n", lineNumber)
				fmt.Println(inputLine)
				os.Exit(1)
			}
		}

		// write result into file
		_, err := resultWriter.WriteString(result + "\n")
		if err != nil {
			log.Fatal(err)
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		lineNumber++
	}
	err = resultWriter.Flush()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	fmt.Println("Usage:")
	fmt.Println("R [key]")
	fmt.Println("W [key] [value]")
	fmt.Println("EXEC [filepath] [resultFilepath]")

	// read commands from the console
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		operationFileds := strings.Fields(scanner.Text())
		var result = ""
		switch len(operationFileds) {
		case 2:
			if strings.EqualFold(operationFileds[0], "R") {
				key := operationFileds[1]
				resultValue, err := client.Read(key)
				if err != nil {
					result = err.Error()
				} else {
					result = "READ\tKey=" + key + "\tValue=" + resultValue
				}
			} else {
				fmt.Println("Invalid Operation!")
			}
		case 3:
			if strings.EqualFold(operationFileds[0], "W") {
				key, value := operationFileds[1], operationFileds[2]
				err := client.Write(key, value)
				if err != nil {
					log.Fatal("Client Write err: ", err)
				}
				result = "WRITE\tKey=" + key + "\tValue=" + value
			} else if strings.EqualFold(operationFileds[0], "EXEC") {
				execBatchOperations(operationFileds[1], operationFileds[2])
			} else {
				fmt.Println("Invalid Operation!")
			}
		default:
			{
				fmt.Println("Invalid Operation!")
			}
		}
		fmt.Println(result)
	}

	if scanner.Err() != nil {
		// Handle error.
	}
}
