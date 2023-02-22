package test

import (
	"log"
	"runtime"
	"shared-registers/client/util"
	"testing"
	"time"
)

// change the test parameter to check whether the timeout or blocking work as expected
// have to observe the after goroutine num as well to assure no gorotine leaks
func TestWaitForMajoritySuccessFromJobs(t *testing.T) {
	quorum := 1
	t.Log("before goroutine: ", runtime.NumGoroutine())

	jobs := make([]func() bool, 0)
	jobs = append(jobs, func() bool {
		time.Sleep(5 * time.Second)
		log.Println("1")
		return true
	})
	jobs = append(jobs, func() bool {
		log.Println("2")
		return false
	})
	jobs = append(jobs, func() bool {
		log.Println("3")
		return false
	})
	util.WaitForMajoritySuccessFromJobs(quorum, time.Second, jobs)
	log.Println("stop blocking: ")
	ticker := time.NewTicker(time.Second)
	done := make(chan bool)
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			log.Println("after: ", runtime.NumGoroutine())
		}
	}
}
func TestPrintFuncExeTime(t *testing.T) {
	defer util.PrintFuncExeTime("test1", time.Now())
	time.Sleep(time.Second)
}
