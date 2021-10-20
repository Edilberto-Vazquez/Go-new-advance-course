package main

import (
	"fmt"
	"sync"
	"time"
)

func ExpensiveFibonacci(n int) int {
	fmt.Printf("Calculate Expensive Fibonnaci fo %d\n", n)
	time.Sleep(5 * time.Second)
	return n
}

type Service struct {
	InProgress map[int]bool
	IsPending  map[int][]chan int
	lock       sync.RWMutex
}

func (s *Service) Work(job int) {
	s.lock.RLock()
	exists := s.InProgress[job]
	if exists {
		s.lock.RUnlock()
		response := make(chan int)
		defer close(response)

		s.lock.Lock()
		s.IsPending[job] = append(s.IsPending[job], response)
		s.lock.Unlock()
		fmt.Printf("Waiting for response: %d\n", job)
		resp := <-response
		fmt.Printf("Response done, received %d\n", resp)
		return
	}
	s.lock.RUnlock()

	s.lock.Lock()
	s.InProgress[job] = true
	s.lock.Unlock()

	fmt.Printf("Calculate Fibonacci for %d", job)
	result := ExpensiveFibonacci(job)
	s.lock.RLock()
	pendingWorkers, exists := s.IsPending[job]
	s.lock.RUnlock()

	if exists {
		for _, pendingWorker := range pendingWorkers {
			pendingWorker <- result
		}
		fmt.Printf("Result sent - all pending workers ready job: %d\n", job)
	}
	s.lock.Lock()
	s.InProgress[job] = false
	s.IsPending[job] = make([]chan int, 0)
	s.lock.Unlock()
}

func NewService() *Service {
	return &Service{
		InProgress: make(map[int]bool),
		IsPending:  make(map[int][]chan int),
	}
}

func main() {
	service := NewService()
	jobs := []int{3, 4, 5, 5, 4, 8, 8, 8, 3}
	var wg sync.WaitGroup
	wg.Add(len(jobs))
	for _, v := range jobs {
		go func(job int) {
			defer wg.Done()
			service.Work(job)
		}(v)
	}
	wg.Wait()
}
