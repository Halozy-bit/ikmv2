package asynctask

import (
	"fmt"
	"time"
)

type runner struct {
	task    []Task
	nextRun []int64
	c       chan int
}

func (r *runner) InsertTask(newTask Task) error {
	if newTask.GetInterval() == emptyInterval {
		return ErrTaskDurationEmpty
	}

	if newTask.GetName() == "" {
		return ErrTaskNameEmpty
	}

	r.task = append(r.task, newTask)
	r.nextRun = append(r.nextRun, time.Now().Unix())
	return nil
}

func (r *runner) Receive() <-chan int {
	return r.c
}

func (r *runner) incrementInterval(i int) {
	r.nextRun[i] = time.Now().Add(
		r.task[i].GetInterval(),
	).Unix()
}

func (r runner) Run(taskNumber int) error {
	if taskNumber > len(r.task)-1 {
		return fmt.Errorf("no task")
	}

	r.task[taskNumber].Run()
	r.incrementInterval(taskNumber)
	return nil
}

func (r *runner) Check() {
	now := time.Now().Unix()
	for i := 0; i < len(r.nextRun); i++ {
		if r.nextRun[i] <= now {
			r.c <- i
		}
	}
}

func newRunner() *runner {
	return &runner{}
}
