package task

import (
	"fmt"
	"reflect"
	"runtime"
)

type Task struct {
	Period   float64
	TaskTime float64
	Delay    float64
	Do       func() error
}

func Create(period float64, taskTime float64, delay float64, do func() error) (task Task, err error) {
	if period < 0 || taskTime < 0 || delay < 0 {
		return task, fmt.Errorf("Period, TaskTime or Delay less than 0")
	}

	if do == nil {
		return task, fmt.Errorf("No function provided")
	}

	task.Period = period
	task.TaskTime = taskTime
	task.Delay = delay
	task.Do = do

	return task, nil
}

func (t Task) Print() {
	fmt.Printf("Period: %v TaskTime: %v Delay: %v Do: %v\n", t.Period, t.TaskTime, t.Delay, runtime.FuncForPC(reflect.ValueOf(t.Do).Pointer()).Name())
}
