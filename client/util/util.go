package util

import (
	"fmt"
	"log"
	"runtime"
	"time"
)

// WaitForMajoritySuccessFromJobs
// hard to figure out a way using waitgroup to return early without gorotine leaks
// use the following mechanism instead
// 1. create a buffered channel to run all jobs concurrently
// 2. write the result of each job to the buffer when finished asynchronously
// 3. wait and read the result for each job from the buffer until the majorityNum times of success reached
// if any result is false, wait for another job from the buffer,
// if timeout happens, stop blocking and return to the caller
func WaitForMajoritySuccessFromJobs(majorityNum int, timeout time.Duration, jobs []func() bool) bool {
	resChan := make(chan bool, 2*len(jobs)) // prevent sender blocking
	timedOut := false
	// after timeout, send majorityNum of success to the channel to stop blocking
	maxTimer := time.AfterFunc(timeout, func() {
		for i := 0; i < majorityNum; i++ {
			resChan <- true
		}
		timedOut = true
	})
	// run each job concurrently and send the result to the channel after finished
	for _, j := range jobs {
		j := j
		go func() {
			jobRes := j()
			//log.Println("job finished")
			resChan <- jobRes
		}()
	}

	// block til received majorityNum of success from the jobs or timeout func
	for i := 0; i < majorityNum; i++ {
		ok := <-resChan // wait for one task to complete
		if !ok {
			i -= 1 // one job failed, decrease the counter to wait for another job from the channel
		}
	}
	// stop the background timer if reached majorNum early
	if !timedOut {
		maxTimer.Stop()
	}
	return timedOut
	// don't need to close the channel
}

func PrintFuncExeTime(funcName string, startTime time.Time) {
	log.Printf("%s took %v to finish\n", funcName, time.Since(startTime))
}

func PrintGoroutineNum(interval time.Duration) {
	ticker := time.NewTicker(interval)

	for range ticker.C {
		fmt.Println("current goroutine cnt: ", runtime.NumGoroutine())
	}
}
