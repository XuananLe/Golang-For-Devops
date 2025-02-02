package main

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Task struct {
	id int
}

var l = sync.RWMutex{};

type Queue interface {
	Enqueue(t Task) error
	Dequeue() (Task, error)
}

type Worker struct {
	id int
}

var WorkerPool = make([]Worker, 100);

func (w *Worker) Process(q Queue, wg *sync.WaitGroup ) {
	defer wg.Done();
	for {
		task, err := q.Dequeue()
		if err != nil {
			return
		}
		fmt.Printf("Worker %d is processing Task %d\n", w.id, task.id)
		time.Sleep(time.Duration(rand.Intn(5)+1) * time.Millisecond * 100) // Simulate processing time
		fmt.Printf("Worker %d done processing Task %d\n", w.id, task.id)
	}
}

func InitializeWorkers(numWorkers int) {
	WorkerPool = make([]Worker, numWorkers)
	for i := 0; i < numWorkers; i++ {
		WorkerPool[i] = Worker{id: i + 1}
	}
}

type NativeChannelQueue struct {
	queue chan Task
}

func doubleSize(c chan Task) chan Task {
	d := make(chan Task, cap(c)*2)
	for len(c) > 0 {
		d <- <-c
	}
	return d
}

func (q *NativeChannelQueue) Enqueue(t Task) error {
	l.Lock();
	defer l.Unlock();
	if len(q.queue)+1 >= cap(q.queue) {
		q.queue = doubleSize(q.queue)
	}
	select {
	case q.queue <- t:
		return nil
	default:
		return errors.New("queue is full")
	}
}

func (q *NativeChannelQueue) Dequeue() (Task, error) {
	l.Lock();
	defer l.Unlock();
	select {
	case task := <-q.queue:
		return task, nil
	default:
		return Task{}, errors.New("Queue is empty")
	}
}

func main() {
	numWorkers := 1000
	numTasks := 100

	q := &NativeChannelQueue{
		queue: make(chan Task, numTasks),
	}

	InitializeWorkers(numWorkers)

	for i := 0; i < numTasks; i++ {
		q.Enqueue(Task{id: i})
	}

	var wg sync.WaitGroup
	start := time.Now()
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go WorkerPool[i].Process(q, &wg)
	}

	wg.Wait()
	fmt.Printf("Time taken with %d workers: %v\n", numWorkers, time.Since(start))

	// === SECOND TEST WITH ONE WORKER ===

	numWorkers = 1
	InitializeWorkers(numWorkers)

	// Reset queue
	q = &NativeChannelQueue{
		queue: make(chan Task, numTasks),
	}

	for i := 0; i < numTasks; i++ {
		q.Enqueue(Task{id: i})
	}

	start = time.Now()
	wg = sync.WaitGroup{}

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go WorkerPool[i].Process(q, &wg)
	}

	wg.Wait() 
	fmt.Printf("Time taken with %d worker: %v\n", numWorkers, time.Since(start))
}