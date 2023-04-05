package asynctask

import (
	"fmt"
	"log"
	"time"
)

// max task
const maxTask = int(16)

type runner struct {
	task    []Task
	nextRun []int64
}

func (r *runner) InsertTask(newTask Task) error {
	if len(r.task) >= maxTask {
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
	log.Print("[task] ", newTask.GetName(), " registered")
	return nil
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

func (r *runner) Check() int {
	now := time.Now().Unix()
	for i := 0; i < len(r.nextRun); i++ {
		next := r.nextRun[i]
		if next <= now {
			return i
		}
	}
	return 17
}

func newRunner() *runner {
	return &runner{}
}
