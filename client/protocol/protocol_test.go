package protocol

import (
	"fmt"
	"go.uber.org/goleak"
	"testing"
	"time"
)

func request300ms() bool {
	res := true
	time.Sleep(300 * time.Millisecond)
	fmt.Println("request300ms", res)
	return res
}

func request500ms() bool {
	res := true
	time.Sleep(500 * time.Millisecond)
	fmt.Println("request500ms", res)
	return res
}

func request1s() bool {
	res := true
	time.Sleep(1 * time.Second)
	fmt.Println("request1s", res)
	return res
}

func TestGoroutine(t *testing.T) {
	go request1s()
	go request500ms()
	go request300ms()
	go func() {
		time.Sleep(300 * time.Millisecond)
		fmt.Println("anonymous function")
	}()
	fmt.Println("main flow")
	time.Sleep(5 * time.Second)
}

func TestBufferedChannel(t *testing.T) {
	chBufferedOne := make(chan int, 1)

	send := func(i int) {
		t.Log("about to send", i)
		chBufferedOne <- i // This will block if the channel is full
		t.Log("sent to queue", i)
	}

	receive := func() {
		t.Log("waiting for value from queue")
		j := <-chBufferedOne // This will block if no data in the channel
		t.Log("received from the queue", j)
	}

	go send(10)
	go send(20)

	time.Sleep(5 * time.Second)
	go receive()
	time.Sleep(5 * time.Second)
	go receive()

	// block main goroutine from exiting
	time.Sleep(999 * time.Hour)
}

func TestUnbufferedChannel(t *testing.T) {
	chUnbuffered := make(chan int) // make(chan int, 0)

	send := func(i int) {
		t.Log("about to send", i)
		chUnbuffered <- i // This will block if the channel is full
		t.Log("sent to queue", i)
	}

	receive := func() {
		t.Log("waiting for value from queue")
		j := <-chUnbuffered // This will block if no data in the channel
		t.Log("received from the queue", j)
	}

	go send(1)

	time.Sleep(5 * time.Second)
	go receive()

	time.Sleep(5 * time.Second)
	send(2) // FIXME: This will block forever as there is no active receiver
}
func TestSelect(t *testing.T) {
	timer1 := time.After(10 * time.Second) // return a channel that will send the timeste after 1s
	timer2 := time.After(3 * time.Second)  // return a channel that will send the current time after 3s

	select {
	case msg := <-timer1:
		t.Log("received message 1", msg)
	case msg := <-timer2:
		t.Log("received message 2", msg)
	default:
		fmt.Println("no message received")
	}
}

// go test -race -run TestRace ./...
func TestRace(t *testing.T) {
	m := make(map[string]string)
	go func() {
		m["1"] = "a" // First conflicting access.
	}()
	m["1"] = "b" // Second conflicting access.

	time.Sleep(time.Second)
}

func TestSync(t *testing.T) {
	defer goleak.VerifyNone(t)

	start := time.Now()
	t.Log("start at", start)
	defer func() { t.Log("end after", time.Since(start)) }()

	followers := 4
	majority := 2
	resChan := make(chan bool, followers) // buffered channel to prevent sender blocking

	replicate := func(d time.Duration, res bool) {
		time.Sleep(d)
		resChan <- res
	}
	// send concurrent request to followers
	go replicate(500*time.Millisecond, true)
	go replicate(300*time.Millisecond, true)
	go replicate(1*time.Second, true)
	go replicate(2*time.Second, true)

	// start a timer to wait for majorityNum of success from the followers
	timeout := time.After(10 * time.Second)
	for i := 0; i < majority; i++ {
		select {
		case ok := <-resChan:
			if !ok {
				i -= 1 // decrease the counter to wait for another job res from the channel
			}
		case <-timeout:
			t.Log("replicate failed")
			return
		}
	}
	t.Log("replicate success")
	time.Sleep(5 * time.Second)
}

func TestSyncGoleak(t *testing.T) {
	defer goleak.VerifyNone(t)

	for i := 0; i < 10; i++ {
		timedOut := WaitForMajoritySuccessFromJobs(2, 100*time.Millisecond, []func() bool{
			request500ms,
			request300ms,
			request1s,
		})
		t.Log("timedOut: ", timedOut)
	}

	time.Sleep(time.Second)
}

// wait for majorityNum of success from the jobs for at most timeout duration
// run all jobs concurrently
// wait for more job result if some job result is false
// return immediately if majorityNum of success is reached
// return false if timeout

// 1. create a buffered channel to store all jobs concurrently
// 2. write the result of each job to the buffer when finished asynchronously
// 3. wait and read the result for each job from the buffer until the majorityNum times of success reached
// if any result is false, wait for another job from the channel,
// if timeout happens, stop waiting and return to the caller
func WaitForMajoritySuccessFromJobs(majorityNum int, timeout time.Duration, followersRequests []func() bool) (succ bool) {
	resChan := make(chan bool, len(followersRequests)) // buffered channel to prevent sender blocking

	// run each job concurrently and send the result to the channel
	for _, f := range followersRequests {
		f := f
		go func() {
			res := f()
			resChan <- res
		}()
	}
	timer := time.After(timeout)
	// block til received majorityNum of success from the jobs or timeout func
	for i := 0; i < majorityNum; i++ {
		select {
		case ok := <-resChan:
			if !ok {
				i -= 1 // decrease the counter to wait for another job res from the channel
			}
		case <-timer:
			fmt.Println("timeout")
			return false
		}
	}
	return true
}
