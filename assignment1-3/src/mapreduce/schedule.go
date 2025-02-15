package mapreduce

import (
	"fmt"
	"sync"
)

// schedule assigns tasks to available workers and ensures all tasks are completed.
func (mr *Master) schedule(phase jobPhase) {
	var ntasks int
	var nios int // Number of inputs (for reduce) or outputs (for map)

	switch phase {
	case mapPhase:
		ntasks = len(mr.files)
		nios = mr.nReduce
	case reducePhase:
		ntasks = mr.nReduce
		nios = len(mr.files)
	}

	debug("Schedule: %v %v tasks (%d I/Os)\n", ntasks, phase, nios)

	var wg sync.WaitGroup
	taskChannel := make(chan int, ntasks) // Queue for task numbers

	// Add tasks to the queue
	for i := 0; i < ntasks; i++ {
		taskChannel <- i
	}
	close(taskChannel) // No more new tasks will be added

	for i := 0; i < ntasks; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for taskNum := range taskChannel {
				for {
					worker := <-mr.registerChannel // Get an available worker

					// Prepare RPC arguments
					task := DoTaskArgs{
						JobName:       mr.jobName,
						File:          mr.files[taskNum],
						Phase:         phase,
						TaskNumber:    taskNum,
						NumOtherPhase: nios,
					}

					// Try sending task to worker
					success := call(worker, "Worker.DoTask", task, nil)
					if success {
						// Re-register worker and break out of retry loop
						go func() { mr.registerChannel <- worker }()
						break
					} else {
						// Worker failed, try another worker
						fmt.Printf("Worker %s failed, retrying task %d\n", worker, taskNum)
					}
				}
			}
		}()
	}

	wg.Wait() // Wait for all tasks to finish
	debug("Schedule: %v phase done\n", phase)
}

// package mapreduce

// // schedule starts and waits for all tasks in the given phase (Map or Reduce).
// func (mr *Master) schedule(phase jobPhase) {
// 	var ntasks int
// 	var nios int // number of inputs (for reduce) or outputs (for map)
// 	switch phase {
// 	case mapPhase:
// 		ntasks = len(mr.files)
// 		nios = mr.nReduce
// 	case reducePhase:
// 		ntasks = mr.nReduce
// 		nios = len(mr.files)
// 	}

// 	debug("Schedule: %v %v tasks (%d I/Os)\n", ntasks, phase, nios)

// 	// All ntasks tasks have to be scheduled on workers, and only once all of
// 	// them have been completed successfully should the function return.
// 	// Remember that workers may fail, and that any given worker may finish
// 	// multiple tasks.
// 	//
// 	// TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO
// 	//
// 	debug("Schedule: %v phase done\n", phase)
// }
