// Key Concepts Fan in and Fan out, Worker Pool Pattern
package main

import (
	"fmt"
	"sync"
	"time"
)

var images = []string{"test1", "test2", "test3", "test4", "test5"}

type Result struct {
	Value string
	Error error
}

func worker(wg *sync.WaitGroup, jobsChan chan string, resChan chan Result) {
	defer wg.Done()
	for url := range jobsChan {
		time.Sleep(time.Millisecond * 50)
		fmt.Printf("Image Processed : %s\n", url)
		resChan <- Result{
			Value: url,
			Error: nil,
		}
	}
}

func main() {
	var wg sync.WaitGroup
	totalWorkers := 2
	resultChan := make(chan Result, 5)
	jobsChan := make(chan string, len(images))
	startTime := time.Now()

	for range totalWorkers {
		wg.Add(1)
		go worker(&wg, jobsChan, resultChan)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// send jobs
	for i := 1; i <= len(images); i++ {
		jobsChan <- fmt.Sprintf("test%d", i)
	}
	close(jobsChan)

	for i := range resultChan {
		fmt.Printf("Result : %s\n", i.Value)
	}
	fmt.Printf("it took %s to process\n", time.Since(startTime))
}
