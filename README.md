## Design and implementation
We designed our MWMR shared registers protocol based on *Attiya, Bar-Noy, and Dolev algorithm*s, and implemented a client program and replica service using Golang. 
### Communication protocol
Our communication protocol is **gRPC**, which sets up HTTP2 long-live connections between each client and replica and transfers the messages encoded by **Protocol Buffers** to reduce payload size of each TCP packet.
### Timestamp
`	`Timestamp is essential in our protocol to maintain the consistency of the requests from different clients. We define our timestamp for each client request as a *<requestNumber, clientID>* structure, the request number refers to the current Write operation times on each key and clientID would be unique for each client. During comparison, the request number is considered first and then clientID if there is a tie.
### Replica
1. Upon start, each client will initialize a built-in Sync.Map, which is a thread-safe Hash Table structure, to serve as the local Key-value store to handle the **Read**() and **Write**() from the clients. The key is string type, and the value is a <value, timestamp> pair from the client.
1. The replica will set up service on port 50051 to handle requests from clients. We expose two RPC functions to the client: **SetPhase**() and **GetPhase**().
- **SetPhase**() will take the input key and value from the user, read the timestamp and compare with the existing timestamp associated with the key in the store, set the *<newValue, newTS>*  to the store only if the upcoming timestamp is bigger. Send ACK to the client in anycase. 
- **GetPhase**() will simply return the *<value, timestamp>* stored locally associated with the key to the client.
### Client
1. We define a client structure in the program. Creating a client object will try to connect with all the replicas with provided addresses and calculated quorum size *f = c/2+1.*
1. The client structure exposes **Read**() and **Write**() functions to the user. Each function will consist of **completeGetPhase**() followed by **completeSetPhase**().  

- **completeGetPhase**() will send concurrent requests to wait for the stored values from the majority of replicas, finding the value associated with the largest timestamp.
- **completeSetPhase**() will send concurrent requests to set the <key, <value, timestamp>> to each replica and wait for the acks from the majority replicas.
- For **Read**() op, we will get the value with the largest timestamp from completeGetPhase(), broadcast the entry to all replicas in completeSetPhase(), and then return the value to the user. This ensures every **Read**() will get the latest value from the majority.
- For **Write**() op, we will get the largest timestamp from **completeGetPhase**(), make new timestamp as *<preRequestNum+1, currClientID>* and pass to **completeSetPhase**() along with the key, value planning to write to store. This value will be written successfully to the replica which does not have a larger timestamp for this key.
- In the **completeGetPhase()** and **completeSetPhase**()**,** to avoid long blocking in the client because of more than majority replica network delays or failures, we introduced a timeout of 1s for each phase. Also, we set the timeout for each request to 500ms to avoid goroutines accumulating when network delays are high in the client side. The early return from the majority result and timeout exit mechanism are implemented using a shared channel between 5 replicas’ concurrent requests. 


Testing correctness
We test for correctness in the situations where there are no server failures, less than a quorum of failures, and greater than or equal to a quorum failures. We also test the situation where there are multiple clients writing and reading.
----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
- #### Less than a quorum of failures: 
For each test, the structure is the following: 1 client, 5 replicas = 5 servers, we write 10 values to 10 keys, then we read from those 10 keys and confirm the output values are correct.

1. Test 1: We crash 2 servers before the **GetPhase** of the **Write**() operation
1. Test 2: We crash 2 servers after the **GetPhase** but before the **SetPhase** of the **Write**() operation
1. Test 3: We crash 2 servers after the **SetPhase** of the **Write**() operation
1. Test 4: We crash 2 servers before the **GetPhase** of the **Read**() operation
1. Test 5: We crash 2 servers after the **GetPhase** but before the **SetPhase** of the **Read**() operation
1. Test 6: We crash 2 servers after the **SetPhase** of the **Read**() operation
- #### Greater than or equal to a quorum of failures: 
For each test, the structure is the following: 1 client, 5 replicas = 5 servers, we write 10 values to 10 keys, then we read from those 10 keys. Since too many replicas fail, we assert in each test that the Write() and Read() operations after the failures throw errors

1. Test 1: We crash 3 servers during **Write**() operation, then confirm the **Write**() call throws an error
1. Test 2: We call **Write**(), then crash 3 servers while calling **Read**(). We check that each **Read**() call throws an error
- #### Multiple clients with failures:
For this test, the structure is the following: 10 clients, 5 replicas = 5 servers. We continuously write different values to the same keys and read from different clients concurrently. At a random time during this process, we kill 2 servers. We confirm that every client reads the same values from registers throughout the process.
## Evaluation: 
Setup • Hardware, specs, n/w latencies, bw:

- 5 node cluster
- Location: cloudlab, Utah
- Hardware: 10 cores, x86- 64 architecture, 8192M memory
- 10000 Mb/s link
### Performance Test A:
1. Describe experiment setup – workload, clients,server location etc. 
- Clients: 1, 2, 4, 8, 16, 32, 64
- Workload: 
  - Read Only
  - Write Only
  - 50% Read 50% Write
- Test duration: 3 minutes
1. Hypothesis

We hypothesize that the system's performance will improve with an increase in the number of clients up to a certain point, after which the throughput will level off. We also expect that the read and write workload will have similar latency and throughput, write heavy workload will have slightly worse performance due to the increased number of K-V pairs in the system.

As the number of clients increases, the throughput is expected to initially increase due to increased parallelism, but will eventually level off when the server's capacity is reached. At the same time, the latency is expected to increase as the server has to handle more concurrent requests.

The initial increase in throughput without increased latency much is due to the network being the bottleneck. However, as the number of clients increases, the server's ability to handle concurrent requests becomes the bottleneck, leading to increased latency.

1. Observation/Result
- Latency vs. Throughput: As predicted, throughput is expected to initially increase without increased latency, but will eventually level off.

![](./images/Aspose.Words.16310126-a02c-437e-8c34-bdf75d358f6e.001.png)

- Latency vs. numClients and Throughput vs. numClients: As predicted, the latency will initially hold but increase later, and the throughput will increase but eventually level off.
  ![](./images/Aspose.Words.16310126-a02c-437e-8c34-bdf75d358f6e.002.png)![](./images/Aspose.Words.16310126-a02c-437e-8c34-bdf75d358f6e.003.png)
1. Draw conclusions

The Shared Register system has similar performance on Reads and Writes. The system's performance is affected by the number of clients, with an initial increase in throughput, but eventually leading to saturation and increased latency. The system's bottleneck shifts from the network to the server's capacity as the number of clients increases, leading to increased latency.
### Performance Test B:
1. Setup: 
   - Clients: 1, 2, 4, 8, 16, 32
   - Workload: 50% Read, 50% Write
   - Test duration: 3 minutes	
   - Leader failure / slowdown at ~1.5 minutes
1. Hypothesis:
- Etcd:
- Leader fail will cause minor dip in throughput since clients must wait for new leader to be elected before it can continue processing requests
- Leader slowdown will cause major drop in throughput since the leader is the bottleneck
- Shared Register:
- Replica fail should not cause any variation since with only 1 failure, there are still a majority of nodes up and the protocol can continue to commit operations as usual
- Replica slow down should also not cause any variation in throughput since a slowdown in one replica should not cause any delays if the remaining replicas are still healthy. The client can continue to commit operations after receiving messages from a majority of nodes.
1. Observation/Results
- Etcd:
  - As predicted, leader fail causes a minor dip in the throughput and then stabilizes again
  - As predicted, leader slowdown causes major drop in throughput that remains since the leader is the bottleneck
- Shared Register:
  - As predicted, a single replica fail in shared registers does not alter the throughput and we continue to see a linearly increasing line throughout the test duration
  - As predicted, a single replica slowdown does not cause a decrease in throughput
  - Fluctuations in throughput when a single replica slowdown are not expected. The reason is that the new RPC go routine is being blocked when there are too many go routine RPCs in flight to prevent issues. However, when one server slows down, each client will have to wait longer for the in-flight RPCs to complete. This leads to fluctuations in throughput, which can be more severe when there are more clients.
1. Conclusions
- In leader-based protocols, the performance of the system is limited by the leader
- With shared registers, the leader-less implementation permits failures/slowdowns without seeing any performance degradations as long as a majority of nodes are still up and healthy

![](./images/Aspose.Words.16310126-a02c-437e-8c34-bdf75d358f6e.004.png)![](./images/Aspose.Words.16310126-a02c-437e-8c34-bdf75d358f6e.005.png)

![](./images/Aspose.Words.16310126-a02c-437e-8c34-bdf75d358f6e.006.png)![](./images/Aspose.Words.16310126-a02c-437e-8c34-bdf75d358f6e.007.png)

![](./images/Aspose.Words.16310126-a02c-437e-8c34-bdf75d358f6e.008.png)![](./images/Aspose.Words.16310126-a02c-437e-8c34-bdf75d358f6e.009.png)

![](./images/Aspose.Words.16310126-a02c-437e-8c34-bdf75d358f6e.010.png)![](./images/Aspose.Words.16310126-a02c-437e-8c34-bdf75d358f6e.011.png)

![](./images/Aspose.Words.16310126-a02c-437e-8c34-bdf75d358f6e.012.png)![](./images/Aspose.Words.16310126-a02c-437e-8c34-bdf75d358f6e.013.png)

![](./images/Aspose.Words.16310126-a02c-437e-8c34-bdf75d358f6e.014.png)![](./images/Aspose.Words.16310126-a02c-437e-8c34-bdf75d358f6e.015.png)

