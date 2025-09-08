package main

import (
	"fmt"
	"sync"
	"time"
)

// ==================== PART 1: BASIC GOROUTINES ====================

// Example 1: Basic goroutine
func basicGoroutine() {
	fmt.Println("\n=== Basic Goroutine Example ===")
	
	// This runs in the main goroutine
	fmt.Println("Main goroutine: Starting")
	
	// Launch a new goroutine
	go func() {
		fmt.Println("Goroutine: Hello from goroutine!")
		time.Sleep(1 * time.Second)
		fmt.Println("Goroutine: Finished sleeping")
	}()
	
	// Main goroutine continues
	fmt.Println("Main goroutine: Launched goroutine")
	
	// Wait for goroutine to finish (don't do this in real code)
	time.Sleep(2 * time.Second)
	fmt.Println("Main goroutine: Done")
}

// ==================== PART 2: WAITGROUPS ====================

// Example 2: Using WaitGroup to wait for goroutines
func waitGroupExample() {
	fmt.Println("\n=== WaitGroup Example ===")
	
	var wg sync.WaitGroup
	
	// Launch multiple goroutines
	for i := 1; i <= 3; i++ {
		wg.Add(1) // Increment counter
		
		go func(id int) {
			defer wg.Done() // Decrement counter when done
			
			fmt.Printf("Worker %d: Starting\n", id)
			time.Sleep(time.Duration(id) * time.Second)
			fmt.Printf("Worker %d: Finished\n", id)
		}(i) // Pass i as parameter to avoid closure issues
	}
	
	wg.Wait() // Wait for all goroutines to finish
	fmt.Println("All workers completed!")
}

// ==================== PART 3: CHANNELS ====================

// Example 3: Basic channels
func basicChannels() {
	fmt.Println("\n=== Basic Channels Example ===")
	
	// Create a channel
	ch := make(chan string)
	
	// Send data in a goroutine
	go func() {
		time.Sleep(1 * time.Second)
		ch <- "Hello from channel!" // Send
	}()
	
	// Receive data (blocks until data is available)
	message := <-ch // Receive
	fmt.Println("Received:", message)
}

// Example 4: Buffered channels
func bufferedChannels() {
	fmt.Println("\n=== Buffered Channels Example ===")
	
	// Create buffered channel (can hold 3 values)
	ch := make(chan int, 3)
	
	// Send values (won't block because buffer has space)
	ch <- 1
	ch <- 2
	ch <- 3
	
	// Receive values
	fmt.Println("Received:", <-ch)
	fmt.Println("Received:", <-ch)
	fmt.Println("Received:", <-ch)
}

// ==================== PART 4: CHANNEL PATTERNS ====================

// Example 5: Producer-Consumer pattern
func producerConsumer() {
	fmt.Println("\n=== Producer-Consumer Example ===")
	
	jobs := make(chan int, 5)
	results := make(chan int, 5)
	
	// Producer
	go func() {
		for i := 1; i <= 5; i++ {
			jobs <- i
			fmt.Printf("Produced job: %d\n", i)
		}
		close(jobs) // Close channel when done
	}()
	
	// Consumer
	go func() {
		for job := range jobs { // Range over channel until closed
			result := job * job
			results <- result
			fmt.Printf("Processed job %d, result: %d\n", job, result)
		}
		close(results)
	}()
	
	// Collect results
	for result := range results {
		fmt.Printf("Final result: %d\n", result)
	}
}

// ==================== PART 5: SELECT STATEMENT ====================

// Example 6: Select statement for non-blocking operations
func selectExample() {
	fmt.Println("\n=== Select Statement Example ===")
	
	ch1 := make(chan string)
	ch2 := make(chan string)
	
	// Send data after different delays
	go func() {
		time.Sleep(1 * time.Second)
		ch1 <- "Message from channel 1"
	}()
	
	go func() {
		time.Sleep(2 * time.Second)
		ch2 <- "Message from channel 2"
	}()
	
	// Select waits for first available channel
	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-ch1:
			fmt.Println("Received:", msg1)
		case msg2 := <-ch2:
			fmt.Println("Received:", msg2)
		case <-time.After(3 * time.Second):
			fmt.Println("Timeout!")
		}
	}
}

// ==================== PART 6: MUTEX FOR SHARED STATE ====================

// Example 7: Race condition problem and solution
type Counter struct {
	value int
	mutex sync.Mutex
}

func (c *Counter) Increment() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.value++
}

func (c *Counter) Get() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.value
}

func mutexExample() {
	fmt.Println("\n=== Mutex Example (Avoiding Race Conditions) ===")
	
	counter := &Counter{}
	var wg sync.WaitGroup
	
	// Launch 100 goroutines that increment counter
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Increment()
		}()
	}
	
	wg.Wait()
	fmt.Printf("Final counter value: %d\n", counter.Get())
}

// ==================== PART 7: WORKER POOL PATTERN ====================

// Example 8: Worker pool for parallel processing
func workerPool() {
	fmt.Println("\n=== Worker Pool Example ===")
	
	const numWorkers = 3
	const numJobs = 10
	
	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)
	
	// Start workers
	var wg sync.WaitGroup
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for job := range jobs {
				fmt.Printf("Worker %d processing job %d\n", workerID, job)
				time.Sleep(500 * time.Millisecond) // Simulate work
				results <- job * 2
			}
		}(w)
	}
	
	// Send jobs
	go func() {
		for j := 1; j <= numJobs; j++ {
			jobs <- j
		}
		close(jobs)
	}()
	
	// Wait for workers to finish
	go func() {
		wg.Wait()
		close(results)
	}()
	
	// Collect results
	for result := range results {
		fmt.Printf("Result: %d\n", result)
	}
}

// ==================== PART 8: CONTEXT FOR CANCELLATION ====================

import "context"

// Example 9: Using context for cancellation
func contextExample() {
	fmt.Println("\n=== Context Example ===")
	
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Goroutine: Context cancelled, stopping work")
				return
			default:
				fmt.Println("Goroutine: Working...")
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()
	
	time.Sleep(5 * time.Second) // Context will timeout after 3 seconds
}

// ==================== PART 9: PIPELINE PATTERN ====================

// Example 10: Pipeline for data processing
func pipelineExample() {
	fmt.Println("\n=== Pipeline Example ===")
	
	// Stage 1: Generate numbers
	numbers := make(chan int)
	go func() {
		defer close(numbers)
		for i := 1; i <= 5; i++ {
			numbers <- i
		}
	}()
	
	// Stage 2: Square numbers
	squares := make(chan int)
	go func() {
		defer close(squares)
		for num := range numbers {
			squares <- num * num
		}
	}()
	
	// Stage 3: Double squares
	doubled := make(chan int)
	go func() {
		defer close(doubled)
		for square := range squares {
			doubled <- square * 2
		}
	}()
	
	// Final stage: Print results
	for result := range doubled {
		fmt.Printf("Pipeline result: %d\n", result)
	}
}

// ==================== MAIN FUNCTION ====================

func main() {
	fmt.Println("Go Concurrency Examples")
	fmt.Println("=======================")
	
	// Run all examples
	basicGoroutine()
	waitGroupExample()
	basicChannels()
	bufferedChannels()
	producerConsumer()
	selectExample()
	mutexExample()
	workerPool()
	contextExample()
	pipelineExample()
	
	fmt.Println("\n=== All examples completed! ===")
}

/*
CONCURRENCY CONCEPTS EXPLAINED:

1. GOROUTINES:
   - Lightweight threads managed by Go runtime
   - Use 'go' keyword to launch
   - Much cheaper than OS threads

2. CHANNELS:
   - Communication mechanism between goroutines
   - "Don't communicate by sharing memory; share memory by communicating"
   - Types: unbuffered (synchronous) and buffered (asynchronous)

3. WAITGROUP:
   - Synchronization primitive to wait for goroutines
   - Add() to increment counter, Done() to decrement, Wait() to block

4. MUTEX:
   - Protects shared memory from race conditions
   - Lock() before accessing, Unlock() after
   - Use defer for automatic unlocking

5. SELECT:
   - Like switch for channels
   - Non-blocking channel operations
   - Can include timeout cases

6. CONTEXT:
   - Carries deadlines, cancellation signals
   - Graceful shutdown and timeout handling

COMMON PATTERNS:
- Worker Pool: Multiple workers processing jobs from a channel
- Pipeline: Chain of processing stages connected by channels  
- Producer-Consumer: One produces data, another consumes
- Fan-out/Fan-in: Distribute work to multiple goroutines, collect results

BEST PRACTICES:
- Always handle channel closure
- Use context for cancellation
- Avoid sharing memory when possible, use channels
- Close channels when done producing
- Use buffered channels for async communication
- Prefer sync.WaitGroup over time.Sleep for waiting
*/