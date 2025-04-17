package main

import "testing"

func TestWorker(t *testing.T) {
    jobs := make(chan int, 1)
    results := make(chan int, 1)
    go worker(1, jobs, results)
    jobs <- 3
    close(jobs)
    if <-results != 6 {
        t.Error("Expected result to be 6")
    }
}
