package asynctask

import (
	"fmt"
	"log"
	"time"
)

type runner struct {
	task    []Task
	nextRun []int64
	C       chan int
}

func (r *runner) InsertTask(newTask Task) error {
	if len(r.task) >= 16 {
		return fmt.Errorf("task full")
	}

	if newTask.GetInterval() == emptyInterval {
		return ErrTaskDurationEmpty
	}

	if newTask.GetName() == "" {
		return ErrTaskNameEmpty
	}

	r.task = append(r.task, newTask)
	r.nextRun = append(r.nextRun, time.Now().Unix())
	log.Print("task ", newTask.GetName(), " registered")
	return nil
}

func (r *runner) Receive() <-chan int {
	return r.C
}

func (r *runner) incrementInterval(i int) {
	r.nextRun[i] = time.Now().Add(
		r.task[i].GetInterval(),
	).Unix()
}

func (r runner) Run(taskNumber int) error {
	tNum := int(taskNumber)
	if tNum > len(r.task)-1 {
		return fmt.Errorf("no task")
	}

	r.task[tNum].Run()
	r.incrementInterval(tNum)
	return nil
}

func (r *runner) Check() {
	now := time.Now().Unix()
	for i := 0; i < len(r.nextRun); i++ {
		next := r.nextRun[i]
		if next <= now {
			r.C <- i
		}
	}
}

func newRunner() *runner {
	return &runner{}
}
