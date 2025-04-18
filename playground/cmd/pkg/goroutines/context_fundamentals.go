/*
# What is Fan-Out and Fan-In

1. Fan-Out : On goroutine spawns multiple worker goroutines to perform tasks concurrently

2. Fan-In : Combines results from multiple channels into a single output channel

# What is sync.WaitGroup

1. A WaitGroup is used to wait for a collection of goroutines to finish. It provides a simple and safe way to synchronize your program without using complex channel coordination

2. Add(n) increments the counter

3. Done() is shorthand for Add(-1)

4. Wait() blocks until the counter hits 0, and then unblocks using a runtime semaphore (via runtime_Semacquire and runtime_Semrelease)
*/
package goroutines

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

func cancel_dowork(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			// context cancelled
			fmt.Println("cancel_dowork cancelled:", ctx.Err())
			return
		default:
			fmt.Println("doing cancel_dowork...")
			time.Sleep(500 * time.Millisecond)
		}
	}
}

// context that cancels after 2 seconds
func ContextCancellationExample() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	go cancel_dowork(ctx)

	time.Sleep(3 * time.Second)
	fmt.Println("::: main finished :::")
}

func fetch_data(ctx context.Context) (string, error) {
	select {
	case <-time.After(5 * time.Second):
		return "Here is your Data!", nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

// if client closes the browser tab, ctx.Done() will unblock and prevent wasting resources
func HTTPHandlerSample(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	data, err := fetch_data(ctx)
	if err != nil {
		http.Error(w, "Request canceled or failed", http.StatusRequestTimeout)
		return
	}

	fmt.Fprintln(w, data)
}

func fan_in_out_fetch(ctx context.Context, url string) (string, error) {
	delay := time.Duration(rand.Intn(5)) * time.Second

	select {
	case <-time.After(delay):
		return "Data from " + url, nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func fanout_fetch(ctx context.Context, urls []string) <-chan string {
	out := make(chan string)

	go func() {
		defer close(out)
		resultCh := make(chan string)

		// FAN-OUT
		for _, url := range urls {
			go func(u string) {
				data, err := fan_in_out_fetch(ctx, u)
				if err != nil {
					fmt.Println("error:", err)
					return
				}
				resultCh <- data
			}(url)
		}

		// FAN-IN
		for i := 0; i < len(urls); i++ {
			select {
			case <-ctx.Done():
				fmt.Println("Cancelled:", ctx.Err())
				return
			case result := <-resultCh:
				out <- result
			}
		}
	}()

	return out
}

// Fan-Out + Fan-In with context that simulate downloading data from multiple sources and aggregating the results.
//
// - fan-out: Spawns a goroutine for each url to fetch concurrently
//
// - fan-in : Collects from resultCh and sends into out channel
//
// - ctx.Done() : Aborts everything if timeout hits
//
// - defer close(out) : Prevent channel leak
func FanOutFanInWithContext() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	urls := []string{"a.com", "b.com", "c.com", "d.com"}

	out := fanout_fetch(ctx, urls)
	for data := range out {
		fmt.Println("Result:", data)
	}
}

// WaitGroup Sample One
func WaitGroupSampleOne() {
	var wg sync.WaitGroup

	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			fmt.Printf("Worker %d done\n", id)
			time.Sleep(time.Second)
		}(i)
	}

	wg.Wait()
	fmt.Println("All workers done!")
}

func wg_worker(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Worker %d starting\n", id)
	time.Sleep(time.Second)
	fmt.Printf("Worker %d done\n", id)
}

// Basic Example – Waiting for multiple goroutines to finish
func WaitingForMultipleGoRoutinesToFinish() {
	var wg sync.WaitGroup

	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go wg_worker(i, &wg)
	}

	wg.Wait()
	fmt.Println("::: All workers Finished :::")
}

func wg_process_file(file string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Processing %s\n", file)
	time.Sleep(500 * time.Millisecond) // Simulate I/O
	fmt.Printf("Done with %s\n", file)
}

// WG - File Processing Example
func WgProcessMultipleFileExample() {
	files := []string{"a.txt", "b.txt", "c.txt"}

	var wg sync.WaitGroup
	for _, file := range files {
		wg.Add(1)
		go wg_process_file(file, &wg)
	}

	wg.Wait()
	fmt.Println("Finished processing all files.")
}

func wg_worker_fan_out(id int, jobs <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		fmt.Printf("Worker %d started job %d\n", id, j)
		time.Sleep(time.Second)
		fmt.Printf("Worker %d finished job %d\n", id, j)
	}
}

// Goroutines with Channels + WaitGroup (Fan-out)
func GoroutineswithChannelsAndWaitGroup() {
	const numWorkers = 3
	jobs := make(chan int, 5)
	var wg sync.WaitGroup

	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go wg_worker_fan_out(w, jobs, &wg)
	}

	for j := 1; j <= 5; j++ {
		jobs <- j
	}
	close(jobs)

	wg.Wait()
	fmt.Println("All jobs done")
}

func wg_do_work(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Goroutine canceled.")
			return
		default:
			fmt.Println("Working...")
		}
	}
}

func CancelableWaitGroupPattern() {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	wg.Add(1)
	go wg_do_work(ctx, &wg)

	cancel()
	wg.Wait()
	fmt.Println("::: Done :::")
}
