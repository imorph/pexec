package main

import (
	"errors"
	"fmt"
	"sync"
	"testing"
)

func TestSimple(t *testing.T) {

	requests := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
	wp := NewWorkerPool(10, len(requests))
	wp.Start()
	defer wp.Stop()
	rspChan := make(chan string, len(requests))
	for _, r := range requests {
		r := r
		wp.SubmitJob(func() error {
			rspChan <- r
			return nil
		})
	}

	wp.WaitResults(1, len(requests))

	close(rspChan)
	rspSet := map[string]struct{}{}
	for rsp := range rspChan {
		rspSet[rsp] = struct{}{}
	}
	if len(rspSet) < len(requests) {
		t.Fatal("Did not handle all requests")
	}
	for _, req := range requests {
		if _, ok := rspSet[req]; !ok {
			t.Fatal("Missing expected values:", req)
		}
	}
}

func TestWithErr(t *testing.T) {

	requests := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
	wp := NewWorkerPool(2, len(requests))
	wp.Start()

	rspChan := make(chan string, len(requests))
	for _, r := range requests {
		r := r
		wp.SubmitJob(func() error {
			if r == "5" {
				return errors.New("ERR")
			}
			rspChan <- r
			return nil
		})
	}

	wp.WaitResults(1, len(requests))
	wp.Stop()

	close(rspChan)
	rspSet := map[string]struct{}{}
	for rsp := range rspChan {
		rspSet[rsp] = struct{}{}
	}
	if len(rspSet) == len(requests) {
		fmt.Println(len(rspSet), len(requests))
		t.Fatal("Did not stop processing on error")
	}
}

func BenchmarkExecute1Worker(b *testing.B) {
	wp := NewWorkerPool(1, b.N)
	wp.Start()
	defer wp.Stop()

	var allDone sync.WaitGroup
	allDone.Add(b.N)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wp.SubmitJob(func() error {
			allDone.Done()
			return nil
		})
	}
	wp.WaitResults(1, b.N)
	allDone.Wait()
}

func BenchmarkExecute2Worker(b *testing.B) {
	wp := NewWorkerPool(2, b.N)
	wp.Start()
	defer wp.Stop()

	var allDone sync.WaitGroup
	allDone.Add(b.N)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wp.SubmitJob(func() error {
			allDone.Done()
			return nil
		})
	}
	wp.WaitResults(1, b.N)
	allDone.Wait()
}

func BenchmarkExecute10Worker(b *testing.B) {
	wp := NewWorkerPool(10, b.N)
	wp.Start()
	defer wp.Stop()

	var allDone sync.WaitGroup
	allDone.Add(b.N)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wp.SubmitJob(func() error {
			allDone.Done()
			return nil
		})
	}
	wp.WaitResults(1, b.N)
	allDone.Wait()
}

func BenchmarkExecute100Worker(b *testing.B) {
	wp := NewWorkerPool(100, b.N)
	wp.Start()
	defer wp.Stop()

	var allDone sync.WaitGroup
	allDone.Add(b.N)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wp.SubmitJob(func() error {
			allDone.Done()
			return nil
		})
	}
	wp.WaitResults(1, b.N)
	allDone.Wait()
}

func BenchmarkExecute1000Worker(b *testing.B) {
	wp := NewWorkerPool(1000, b.N)
	wp.Start()
	defer wp.Stop()

	var allDone sync.WaitGroup
	allDone.Add(b.N)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wp.SubmitJob(func() error {
			allDone.Done()
			return nil
		})
	}
	wp.WaitResults(1, b.N)
	allDone.Wait()
}

func BenchmarkExecute10000Worker(b *testing.B) {
	wp := NewWorkerPool(10000, b.N)
	wp.Start()
	defer wp.Stop()

	var allDone sync.WaitGroup
	allDone.Add(b.N)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wp.SubmitJob(func() error {
			allDone.Done()
			return nil
		})
	}
	wp.WaitResults(1, b.N)
	allDone.Wait()
}
