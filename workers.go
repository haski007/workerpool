package workerpool

import (
	"fmt"
	"sync"
	"time"
)

type Job struct {
	ID          int
	WorkFunc    func() error
	HandleError func(err error)
}

type WorkerPool struct {
	workers   int
	jobQueue  chan Job
	waitGroup sync.WaitGroup
	// is info in console needed
	v bool
}

func NewWorkerPool(workers int, info bool) *WorkerPool {
	return &WorkerPool{
		workers:  workers,
		jobQueue: make(chan Job),
		v:        info,
	}
}

func (wp *WorkerPool) Start() {
	for i := 0; i < wp.workers; i++ {
		go wp.worker(i)
	}
}

func (wp *WorkerPool) info(msg string) {
	if wp.v {
		fmt.Print(msg)
	}
}

func (wp *WorkerPool) worker(id int) {
	wp.waitGroup.Add(1)
	defer wp.waitGroup.Done()

	for job := range wp.jobQueue {
		wp.info(fmt.Sprintf("Worker %d started job %d\n", id, job.ID))
		err := job.WorkFunc()
		if err != nil {
			job.HandleError(err)
		}
	}
}

func (wp *WorkerPool) Submit(job Job) {
	wp.jobQueue <- job
}

func (wp *WorkerPool) Stop() {
	close(wp.jobQueue)
	wp.waitGroup.Wait()
}

func Throttle(interval time.Duration) func() {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	return func() {
		<-ticker.C
	}
}

