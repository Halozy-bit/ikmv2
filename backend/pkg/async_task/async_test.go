package asynctask

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/ikmv2/backend/pkg/repository"
	"github.com/stretchr/testify/assert"
)

func logger(str string) {
	log.Print("hello ", str)
}

func TestNewTask(t *testing.T) {
	task1, err := NewTask(repository.RandName(true), time.Second*1, logger, repository.RandName(true))
	assert.NoError(t, err)
	task1.Run()
	defer func(t *testing.T) {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
			return
		}
		t.Fail()
	}(t)
	NewTask(repository.RandName(true), time.Second*1, logger)
}

func TestRunner(t *testing.T) {
	rn := &runner{C: make(chan int)}
	amountTask := 5
	for i := 0; i < amountTask; i++ {
		randDur := time.Duration(int(repository.RandInt(1, 5)))
		task, err := NewTask(repository.RandName(true), randDur*time.Second, logger, repository.RandName(true))
		assert.NoError(t, err)
		err = rn.InsertTask(task)
		assert.NoError(t, err)
	}

	go func(r *runner) {
		for {
			run := <-r.Receive()
			r.Run(run)
		}
	}(rn)

	time.Sleep(time.Second * 6)
	rn.Check()
}

func TestAsync(t *testing.T) {
	amountTask := 1
	for i := 0; i < amountTask; i++ {
		randDur := time.Duration(int(repository.RandInt(1, 10)))
		name := repository.RandName(true)
		task, err := NewTask(name, randDur*time.Second, logger, name)
		t.Log(task.GetName(), " interval ", task.GetInterval())
		assert.NoError(t, err)
		err = AddTask(task)
		assert.NoError(t, err)
	}

	Start(time.Second * 1)
	time.Sleep(time.Second * 20)
	log.Println("sending stop signal")
	Stop()
}
