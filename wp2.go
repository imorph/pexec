package main

import (
	"log"
	"time"
)

// Job for exec
type Job func() error

// Worker abstraction
type Worker struct {
	ID int
}

// WorkerPool have size (number of workers) "depth" -- how many jobs fits into channels and three channels
type WorkerPool struct {
	poolSize int
	bufSize  int
	jobs     chan Job
	results  chan error
	exit     chan bool
}

// NewWorkerPool constructor
func NewWorkerPool(poolSize int, bufSize int) *WorkerPool {
	if poolSize < 1 {
		poolSize = 1
	}

	if bufSize < 10 {
		bufSize = 10
	}

	wp := &WorkerPool{
		poolSize: poolSize,
		bufSize:  bufSize,
		jobs:     make(chan Job, bufSize),
		results:  make(chan error, bufSize),
		exit:     make(chan bool),
	}
	return wp
}

// NewWorker constructor
func NewWorker(ID int) Worker {
	w := Worker{
		ID: ID,
	}
	return w
}

// Start new worker
func (w Worker) Start(jobs chan Job, results chan error, exit chan bool) {
	for {
		select {
		// non bloking try
		case <-exit:
			log.Println("WORKER ==> got exit signal stoppping")
			return
		default:
			// get jobs or go to sleep
			select {
			case fn := <-jobs:
				log.Println("WORKER ==> worker", w.ID, "started  job")
				err := fn()
				log.Println("WORKER ==> worker", w.ID, "finished job")
				results <- err
			default:
				time.Sleep(time.Millisecond * 1)
			}
		}
	}
}

// Start new worker pool
func (wp *WorkerPool) Start() {
	var wrk Worker
	for w := 1; w <= wp.poolSize; w++ {
		wrk = NewWorker(w)
		go wrk.Start(wp.jobs, wp.results, wp.exit)
	}
}

// SubmitJob submits new job to WorkerPool
func (wp *WorkerPool) SubmitJob(j Job) {
	wp.jobs <- j
}

// WaitResults blocks on getting results, ether all results from fixed batch will be done or WP will be stopped (maxErrNum exceeded)
func (wp *WorkerPool) WaitResults(maxErrNum int, maxJobsNum int) {
	errNum := 0
	jobCount := 0
	for err := range wp.results {
		jobCount++
		log.Println("RESULT ==> Jobs done:", jobCount)
		if jobCount == maxJobsNum {
			log.Println("RESULT ==> All jobs done. Unblocking.")
			return
		}
		if err != nil {
			log.Println("RESULT ==> error detected!")
			errNum++
			if errNum == maxErrNum {
				log.Println("RESULT ==> Too many errors, unblocking. WorkerPool will be stopped.")
				return
			}
		}
	}
}

// Stop worker pool
func (wp *WorkerPool) Stop() {
	for w := 1; w <= wp.poolSize; w++ {
		log.Println("STOPWP ==> Send stop to worker")
		wp.exit <- true
	}
}
