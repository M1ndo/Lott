package utils

import "sync"

const numWorkers = 1 // Adjust based on your system's capabilities

type Job struct {
	Year        int
	Result      map[string][]uint64
	AllResults  *map[string]map[string][]uint64
	mutex       sync.Mutex
}

// jobQueue is a buffered channel that we can send work requests on.
var jobQueue = make(chan Job, numWorkers)
// done is a channel to signal when a worker is finished and can receive another job.
var done = make(chan bool, numWorkers)

// startWorkers creates numWorkers goroutines that listen for jobs on jobQueue.
func (m *DB) startWorkers(numWorkers int) {
	for i := 0; i < numWorkers; i++ {
		m.L.Logger.Log.Info("Starting Many Times")
		go m.worker(i, jobQueue, done)
	}
}

// worker is a function that gets jobs from jobQueue and processes them.
func (m *DB) worker(id int, jobQueue chan Job, done chan bool) {
	for j := range jobQueue {
		Years := YearSuccessive[j.Year]
		count := 0
		for date, row := range j.Result {
			dataa := m.Improvedcentry(row, date, Years)
			j.mutex.Lock() // Acquire the lock before modifying the AllResults map
			(*j.AllResults)[date] = dataa
			j.mutex.Unlock() // Release the lock after modifying the AllResults map
			count++
			if count > 9 {
				break
			}
		}
		done <- true
	}
}

// AddJob adds a job to the queue.
func addJob(j Job) {
	jobQueue <- j
}

// CloseAndWait closes the jobQueue and waits for all jobs to finish.
func closeAndWait() {
	close(jobQueue)
	for i := 0; i < cap(done); i++ {
		<-done
	}
}
