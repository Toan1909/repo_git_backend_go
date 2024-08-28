package helper

import (
	"fmt"
	"sync"
)

// Job interface defines a single method that any job should implement
type Job interface {
	Process()
}

// JobQueue manages the queue of jobs and the pool of workers
type JobQueue struct {
	JobQueue   chan Job         // Channel for incoming jobs
	WorkerPool chan chan Job    // Channel for available workers
	Workers    []*Worker        // Slice to keep track of workers
	wg         *sync.WaitGroup  // WaitGroup to wait for all workers to finish
}

// Worker represents a single worker that executes jobs
type Worker struct {
	ID         int            // ID of the worker
	JobQueue   chan Job       // Channel for jobs assigned to this worker
	WorkerPool chan chan Job  // Channel for worker pool to get jobs
	quit       chan bool      // Channel to signal the worker to stop
}

// NewWorker creates a new worker with a given ID and worker pool
func NewWorker(id int, workerPool chan chan Job) *Worker {
	return &Worker{
		ID:         id,
		JobQueue:   make(chan Job),
		WorkerPool: workerPool,
		quit:       make(chan bool),
	}
}

// Start runs the worker in a separate goroutine and listens for jobs or quit signals
func (w *Worker) Start(wg *sync.WaitGroup) {
	go func() {
		defer wg.Done()
		for {
			// Add this worker's job queue to the worker pool
			w.WorkerPool <- w.JobQueue

			select {
			case job := <-w.JobQueue:
				// Process the job
				job.Process()
			case <-w.quit:
				// Stop the worker
				fmt.Printf("Worker %d stopping\n", w.ID)
				return
			}
		}
	}()
}

// Stop sends a signal to stop the worker
func (w *Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}

// NewJobQueue creates a new job queue with a specified number of workers
func NewJobQueue(numWorkers int) *JobQueue {
	workerPool := make(chan chan Job, numWorkers)
	jobQueue := make(chan Job)
	wg := &sync.WaitGroup{}

	queue := &JobQueue{
		JobQueue:   jobQueue,
		WorkerPool: workerPool,
		Workers:    make([]*Worker, numWorkers),
		wg:         wg,
	}

	// Create and start workers
	for i := 0; i < numWorkers; i++ {
		worker := NewWorker(i, workerPool)
		queue.Workers[i] = worker
		wg.Add(1)
		worker.Start(wg)
	}

	// Start the job dispatcher
	go queue.dispatch()
	return queue
}

// Start all workers
func (q *JobQueue) Start() {
	for _, worker := range q.Workers {
		q.wg.Add(1)
		worker.Start(q.wg)
	}
}

// dispatch listens for incoming jobs and assigns them to available workers
func (q *JobQueue) dispatch() {
	for {
		select {
		case job := <-q.JobQueue:
			go func(job Job) {
				// Wait for an available worker
				jobChannel := <-q.WorkerPool
				// Assign the job to the worker
				jobChannel <- job
			}(job)
		}
	}
}

// Submit adds a job to the job queue
func (q *JobQueue) Submit(job Job) {
	q.JobQueue <- job
}

// Stop stops all workers and waits for them to finish
func (q *JobQueue) Stop() {
	for _, worker := range q.Workers {
		worker.Stop()
	}
	q.wg.Wait()
}
