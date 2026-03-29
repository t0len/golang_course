package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func solutionMutex() {
	fmt.Println("=== Solution 1: sync.Mutex ===")

	var counter int
	var mu sync.Mutex
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			counter++
			mu.Unlock()
		}()
	}

	wg.Wait()
	fmt.Printf("Counter: %d\n\n", counter)
}

func solutionAtomic() {
	fmt.Println("=== Solution 2: sync/atomic ===")

	var counter int64
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			atomic.AddInt64(&counter, 1)
		}()
	}

	wg.Wait()
	fmt.Printf("Counter: %d\n\n", counter)
}

func main() {
	solutionMutex()
	solutionAtomic()
}
