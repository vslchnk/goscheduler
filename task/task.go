package task

import (
	"context"
	"fmt"
	"reflect"
	"runtime"
	"time"
)

type Task struct {
	Name     string
	Period   time.Duration
	TaskTime time.Duration
	Delay    time.Duration
	Do       func(ctx context.Context) error
}

func Create(name string, period time.Duration, taskTime time.Duration, delay time.Duration, do func(ctx context.Context) error) (task Task, err error) {
	if period < 0 || taskTime < 0 || delay < 0 {
		return task, fmt.Errorf("Period, TaskTime or Delay less than 0")
	}

	if do == nil {
		return task, fmt.Errorf("No function provided")
	}

	task.Name = name
	task.Period = period
	task.TaskTime = taskTime
	task.Delay = delay
	task.Do = do

	return task, nil
}

func (t Task) Print() {
	fmt.Printf("Name: %v; Period: %v; TaskTime: %v; Delay: %v; Do: %v\n", t.Name, t.Period, t.TaskTime, t.Delay, runtime.FuncForPC(reflect.ValueOf(t.Do).Pointer()).Name())
}
