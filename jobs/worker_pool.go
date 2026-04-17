package jobs

import (
	"fmt"
	"sync"
	"time"
)

func StartWorkerPool(workerCount int, queue *Queue, wg *sync.WaitGroup) {
	for i := 1; i <= workerCount; i++ {
		workerID := i

		go func() {
			for job := range queue.Jobs {
				fmt.Printf("Worker %d picked up job %d (%s) for item %d\n", workerID, job.ID, job.Type, job.ItemID)

				// Simulate background processing
				time.Sleep(2 * time.Second)

				fmt.Printf("Worker %d finished job %d (%s) for item %d\n", workerID, job.ID, job.Type, job.ItemID)

				wg.Done()
			}
		}()
	}
}
