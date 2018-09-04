package task

import (
	"context"
	"fmt"
	"reflect"
	"runtime"
	"time"
)

type Task struct {
	name     string
	period   time.Duration
	taskTime time.Duration
	delay    time.Duration
	do       func(ctx context.Context) error
}

// creates task
func Create(name string, period time.Duration, taskTime time.Duration, delay time.Duration, do func(ctx context.Context) error) (task Task, err error) {
	if period < 0 || taskTime < 0 || delay < 0 {
		return task, fmt.Errorf("Period, TaskTime or Delay less than 0")
	}

	if do == nil {
		return task, fmt.Errorf("No function provided")
	}

	task.name = name
	task.period = period
	task.taskTime = taskTime
	task.delay = delay
	task.do = do

	return task, nil
}

// prints task's parametres
func (t Task) Print() {
	fmt.Printf("Name: %v; Period: %v; TaskTime: %v; Delay: %v; Do: %v\n", t.name, t.period, t.taskTime, t.delay, runtime.FuncForPC(reflect.ValueOf(t.do).Pointer()).Name())
}

// sets name of task
func (t *Task) SetName(name string) {
	t.name = name
}

// returns task's name
func (t Task) GetName() string {
	return t.name
}

// sets task's period
func (t *Task) SetPeriod(period time.Duration) error {
	if period < 0 {
		return fmt.Errorf("Period is less than 0")
	}

	t.period = period

	return nil
}

// returns task's period
func (t *Task) GetPeriod() time.Duration {
	return t.period
}

// sets task's time
func (t *Task) SetTaskTime(taskTime time.Duration) error {
	if taskTime < 0 {
		return fmt.Errorf("Task time is less than 0")
	}

	t.taskTime = taskTime

	return nil
}

// returns taks's time
func (t *Task) GetTaskTime() time.Duration {
	return t.taskTime
}

// sets task's delay
func (t *Task) SetDelay(delay time.Duration) error {
	if delay < 0 {
		return fmt.Errorf("Delay is less than 0")
	}

	t.delay = delay

	return nil
}

// returns task's delay
func (t *Task) GetDelay() time.Duration {
	return t.delay
}

// sets function of task
func (t *Task) SetDoFunc(do func(ctx context.Context) error) error {
	if do == nil {
		return fmt.Errorf("No function provided")
	}

	t.do = do

	return nil
}

// returns task func
func (t *Task) GetDoFunc() func(ctx context.Context) error {
	return t.do
}
