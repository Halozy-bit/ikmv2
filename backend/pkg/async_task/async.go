package asynctask

import (
	"fmt"
	"reflect"
	"time"
)

var async worker

func init() {
	async = worker{
		r:         newRunner(),
		isRunning: false,
	}
}

func NewTask(name string, interval time.Duration, functionName interface{}, params ...interface{}) (Task, error) {
	fn := reflect.ValueOf(functionName)
	if fn.Kind() != reflect.Func {
		return TaskImp{}, fmt.Errorf("not kind of function")
	}

	methodParam := make([]reflect.Value, fn.Type().NumIn())

	for i := 0; i < fn.Type().NumIn(); i++ {
		methodParam[i] = reflect.ValueOf(params[i])
	}

	return TaskImp{
		TaskIdentifier: TaskIdentifier{
			Name:     name,
			Interval: interval,
		},
		fn:    fn,
		param: methodParam,
	}, nil
}

func AddTask(newTask Task) error {
	return async.AddTask(newTask)
}

func Start(refreshDur time.Duration) error {
	return async.Start(refreshDur)
}

func Stop() {
	async.Stop()
}
