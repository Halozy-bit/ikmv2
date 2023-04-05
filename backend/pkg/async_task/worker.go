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
}

func (w worker) AddTask(newTask Task) error {
	return w.r.InsertTask(newTask)
}

func (w *worker) checker(refreshDur time.Duration, stopChecker <-chan struct{}) {
	defer close(w.stopChecker)
	log.Print("[async pool] running real time task check")
	ticker := time.NewTicker(refreshDur)
free:
	for {
		taskNumber := w.r.Check()
		if taskNumber > maxTask {
			continue
		}

		go w.do(taskNumber)

		select {
		case <-stopChecker:
			break free
		default:
			<-ticker.C
		}
	}
	log.Println("[async pool] stopped")
}

func (w *worker) do(taskNumber int) {
	log.Print("[task] Running: ", w.r.task[taskNumber].GetName())
	w.r.Run(taskNumber)
	log.Print("[task] exiting: ", w.r.task[taskNumber].GetName())
}

func (w *worker) Start(refreshDur time.Duration) error {
	if w.isRunning {
		return fmt.Errorf("worker already running")
	}

	w.stopChecker = make(chan struct{})

	go w.checker(refreshDur, w.stopChecker)
	return nil
}

func (w *worker) Stop() {
	// signal to stop checker
	w.stopChecker <- struct{}{}
	// change state running to false
	w.isRunning = false
}
