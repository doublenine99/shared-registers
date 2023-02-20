package util

import "time"

// WaitForMajoritySuccessFromJobs
// hard to figure out a way using waitgroup to return early without gorotine leaks
// use the following mechanism instead
// 1. create a buffered channel to run all jobs concurrently
// 2. write the result of each job to the buffer when finished asynchronously
// 3. wait and read the result for each job from the buffer until the majorityNum times of success reached
// if any result is false, wait for another job from the buffer,
// if timeout happens, stop blocking and return to the caller
func WaitForMajoritySuccessFromJobs(majorityNum int, timeout time.Duration, jobs []func() bool) (timedOut bool) {
	resChan := make(chan bool, len(jobs))

	for _, j := range jobs {
		go func(j func() bool) {
			jobRes := j()
			resChan <- jobRes
		}(j)
	}

	// set a max wait time for all the jobs
	maxTimer := time.AfterFunc(timeout, func() {
		for i := 0; i < majorityNum; i++ {
			resChan <- true
		}
		timedOut = true
	})
	// block until received majorityNum of success from the job or timeout
	for i := 0; i < majorityNum; i++ {
		ok := <-resChan // wait for one task to complete
		if !ok {
			i -= 1 // one job failed, decrease the counter to wait for another job from the channel
		}
	}
	// stop the background timer if reached majorNum early
	if !timedOut {
		_ = maxTimer.Stop()
	}
	return timedOut
	// don't need to close the channel
}
