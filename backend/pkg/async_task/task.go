package asynctask

import (
	"reflect"
	"time"
)

type Task interface {
	Run()
	GetName() string
	GetInterval() time.Duration
}

type TaskIdentifier struct {
	Name     string
	Interval time.Duration
}

func (ti TaskIdentifier) GetInterval() time.Duration {
	return ti.Interval
}

func (ti TaskIdentifier) GetName() string {
	return ti.Name
}

type TaskImp struct {
	TaskIdentifier
	fn    reflect.Value
	param []reflect.Value
}

func (ti TaskImp) Run() {
	ti.fn.Call(ti.param)
}

func (ti TaskImp) GetName() string {
	return ti.Name
}
