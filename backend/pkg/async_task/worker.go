package asynctask

import (
	"fmt"
	"log"
	"time"
)

var ErrTaskDurationEmpty = fmt.Errorf("cannot insert task with empty duration")
var ErrTaskNameEmpty = fmt.Errorf("cannot insert task with empty name")
var emptyInterval time.Duration

type worker struct {
	r         *runner
	isRunning bool
	done      chan struct{}
}

func (w worker) AddTask(newTask Task) error {
	return w.r.InsertTask(newTask)
}

func (w *worker) checker(refreshDur time.Duration, r *runner, done <-chan struct{}) {
	defer close(w.r.c)
	defer close(w.done)
	ticker := time.NewTicker(refreshDur)
free:
	for {
		r.Check()
		select {
		case <-done:
			break free
		default:
			<-ticker.C
		}
	}
}

func (w *worker) do(r *runner, done <-chan struct{}) {
free:
	for {
		select {
		case <-done:
			break free
		case call := <-r.Receive():
			r.Run(call)
		}
	}
	log.Println("worker stopped")
}

func (w *worker) Start(refreshDur time.Duration) error {
	if w.isRunning {
		return fmt.Errorf("worker already running")
	}

	w.r.c = make(chan int)
	w.done = make(chan struct{})

	go w.checker(refreshDur, w.r, w.done)
	go w.do(w.r, w.done)
	return nil
}

func (w *worker) Stop() {
	w.done <- struct{}{}
}
