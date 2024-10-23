package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
)

// Task represents a unit of work.
type Task struct {
  id      int
  data    string
  retries int // To track the number of retries
}

var completedTasks = 0

const taskLength = 5

// Worker function that processes tasks. If a worker fails, the task will be sent to failChan.
func worker(id int, taskChan <-chan Task, wg *sync.WaitGroup, failChan chan<- Task, retryLimit int) {
  defer wg.Done()

  for task := range taskChan {
    fmt.Printf(Blue+"[INFO] Worker %d assigned : %s"+Reset+"\n", id, task.data)

    // Simulate random failure (30% chance of failure)
    if rand.Float32() < 0.3 {
      fmt.Printf(Red + "[ERROR] Worker %d failed on task %d"+Reset+"\n", id, task.id)
      task.retries++

      if task.retries > retryLimit {
        fmt.Printf(Magenta + "[Warning] Task %d reached max retries (%d), logging as failed"+Reset+"\n", task.id, retryLimit)
		completedTasks += 1
	if completedTasks == taskLength{
		close(failChan)
	}
      } else {
        failChan <- task // Send the failed task for reassignment
      }
      return
    }

    // Simulate task processing time
    time.Sleep(time.Duration(rand.Intn(3)+1) * time.Second)
    fmt.Printf(Green + "[SUCCESS] Worker %d completed task %d"+Reset+"\n", id, task.id)

	completedTasks += 1
	if completedTasks == taskLength{
		close(failChan)
	}
  }
}

func main() {
  rand.Seed(time.Now().UnixNano())

  // Define a set of tasks to be executed
  tasks := []Task{
    {id: 1, data: "Task 1"},
    {id: 2, data: "Task 2"},
    {id: 3, data: "Task 3"},
    {id: 4, data: "Task 4"},
    {id: 5, data: "Task 5"},
  }

  // Channels for task distribution and failure handling
  taskChan := make(chan Task, len(tasks))
  failChan := make(chan Task, len(tasks))

  // WaitGroup to ensure all workers finish their tasks
  var wg sync.WaitGroup

  // Number of workers (simulating processors)
  numWorkers := 3
  retryLimit := 1 // Maximum retries for a failed task

  // Start worker goroutines
  for i := 1; i <= numWorkers; i++ {
    wg.Add(1)
    go worker(i, taskChan, &wg, failChan, retryLimit)
  }

  // Distribute tasks to workers using a round-robin approach
  for _,task := range tasks{
	taskChan <- task
  }

  // Handle failed tasks by redistributing them with a retry mechanism
  go func() {
	currentWorkerId := 0
    for failedTask := range failChan {
      fmt.Printf(Yellow + "[INFO] Reassigning failed task %d"+Reset+"\n", failedTask.id)
      
      fmt.Printf(Yellow + "[INFO] Reassigning task %d to worker %d after failure"+Reset+"\n", failedTask.id, currentWorkerId + 1)
      wg.Add(1)
	  taskChan <- failedTask
      go worker(currentWorkerId + 1, taskChan, &wg, failChan, retryLimit)
	  currentWorkerId = (currentWorkerId + 1) % numWorkers
    }
	close(taskChan)
  }()

  // Wait for all workers to finish
  wg.Wait()

  fmt.Println("All tasks completed.")
}
