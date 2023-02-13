build:
	go build -o client src/client.go
	go build -o server src/replica.go
	
clean:
	rm client
	rm replica

