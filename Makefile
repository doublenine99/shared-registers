build:
	go build -o client src/client.go
	go build -o replica src/replica.go
	
clean:
	rm client
	rm replica

