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
	r           *runner
	isRunning   bool
	stopChecker chan struct{}
	stopWorker  chan struct{}
}

func (w worker) AddTask(newTask Task) error {
	return w.r.InsertTask(newTask)
}

func (w *worker) checker(refreshDur time.Duration, r *runner, stopChecker <-chan struct{}) {
	defer close(w.r.C)
	defer close(w.stopChecker)
	log.Print("running goroutine checker")
	ticker := time.NewTicker(refreshDur)
free:
	for {
		r.Check()
		select {
		case <-stopChecker:
			break free
		default:
			<-ticker.C
		}
	}
}

func (w *worker) do(r *runner, stopWorker <-chan struct{}) {
	log.Print("running goroutine worker")
free:
	for {
		select {
		case <-stopWorker:
			break free
		case call := <-r.C:
			log.Print("Running: ", w.r.task[call].GetName())
			r.Run(call)
			log.Print("exiting: ", w.r.task[call].GetName())
		}
	}
	log.Println("worker stopped")
}

func (w *worker) Start(refreshDur time.Duration) error {
	if w.isRunning {
		return fmt.Errorf("worker already running")
	}

	w.r.C = make(chan int)
	w.stopWorker = make(chan struct{})
	w.stopChecker = make(chan struct{})

	go w.checker(refreshDur, w.r, w.stopChecker)
	go w.do(w.r, w.stopWorker)
	return nil
}

func (w *worker) Stop() {
	w.stopWorker <- struct{}{}
	w.stopChecker <- struct{}{}
	w.isRunning = false
}
