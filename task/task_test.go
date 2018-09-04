package task

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func outer(name string) func(ctx context.Context) error {
	text := "Modified " + name

	foo := func(ctx context.Context) error {
		fmt.Println(text)

		return nil
	}

	return foo
}

func Test_Task_Create(t *testing.T) {
	foo := outer("hello")

	_, err := Create("printing", time.Second*3, time.Second*3, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}
}

func Test_Task_CreateWithNilFunc(t *testing.T) {

	_, err := Create("printing", time.Second*3, time.Second*3, time.Second*1, nil)

	if err == nil {
		t.Error("Failed to detect error while creating task: ", err)
	}
}

func Test_Task_CreateWithNegativePeriod(t *testing.T) {
	foo := outer("hello")

	_, err := Create("printing", time.Second*(-3), time.Second*3, time.Second*1, foo)

	if err == nil {
		t.Error("Failed to detect error while creating task: ", err)
	}
}

func Test_Task_CreateWithNegativeTaskTime(t *testing.T) {
	foo := outer("hello")

	_, err := Create("printing", time.Second*3, time.Second*(-3), time.Second*1, foo)

	if err == nil {
		t.Error("Failed to detect error while creating task: ", err)
	}
}

func Test_Task_CreateWithNegativeDelayTime(t *testing.T) {
	foo := outer("hello")

	_, err := Create("printing", time.Second*3, time.Second*3, time.Second*(-1), foo)

	if err == nil {
		t.Error("Failed to detect error while creating task: ", err)
	}
}

func Test_Task_Print(t *testing.T) {
	foo := outer("hello")

	task, err := Create("printing", time.Second*3, time.Second*3, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	task.Print()
}

func Test_Task_SetGetName(t *testing.T) {
	foo := outer("hello")

	task, err := Create("printing", time.Second*3, time.Second*3, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	task.SetName("printingChanged")

	if task.GetName() != "printingChanged" {
		t.Error("Failed to set new name: names not the same")
	}
}

func Test_Task_SetGetPeriod(t *testing.T) {
	foo := outer("hello")

	task, err := Create("printing", time.Second*3, time.Second*3, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	err = task.SetPeriod(time.Second * 2)

	if err != nil {
		t.Error("Failed to set new period: ", err)
	}

	if task.GetPeriod() != time.Second*2 {
		t.Error("Failed to set new period: periods not the same", err)
	}
}

func Test_Task_SetPeriodFailure(t *testing.T) {
	foo := outer("hello")

	task, err := Create("printing", time.Second*3, time.Second*3, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	err = task.SetPeriod(time.Second * (-2))

	if err == nil {
		t.Error("Failed to detect error while setting new period: ", err)
	}
}

func Test_Task_SetGetTaskTime(t *testing.T) {
	foo := outer("hello")

	task, err := Create("printing", time.Second*3, time.Second*3, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	err = task.SetTaskTime(time.Second * 2)

	if err != nil {
		t.Error("Failed to set new task time: ", err)
	}

	if task.GetTaskTime() != time.Second*2 {
		t.Error("Failed to set new task time: task time not the same", err)
	}
}

func Test_Task_SetTaskTimeFailure(t *testing.T) {
	foo := outer("hello")

	task, err := Create("printing", time.Second*3, time.Second*3, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	err = task.SetTaskTime(time.Second * (-2))

	if err == nil {
		t.Error("Failed to detect error while setting new task time: ", err)
	}
}

func Test_Task_SetGetDelay(t *testing.T) {
	foo := outer("hello")

	task, err := Create("printing", time.Second*3, time.Second*3, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	err = task.SetDelay(time.Second * 2)

	if err != nil {
		t.Error("Failed to set new delay: ", err)
	}

	if task.GetDelay() != time.Second*2 {
		t.Error("Failed to set new delay: delay not the same", err)
	}
}

func Test_Task_SetDelayFailure(t *testing.T) {
	foo := outer("hello")

	task, err := Create("printing", time.Second*3, time.Second*3, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	err = task.SetDelay(time.Second * (-2))

	if err == nil {
		t.Error("Failed to detect error while setting new delay: ", err)
	}
}

func Test_Task_SetGetDoFunc(t *testing.T) {
	foo := outer("hello")

	task, err := Create("printing", time.Second*3, time.Second*3, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	foo2 := outer("hello2")

	err = task.SetDoFunc(foo2)

	if err != nil {
		t.Error("Failed to set new func: ", err)
	}

	if task.GetDoFunc()(context.Background()) != foo2(context.Background()) {
		t.Error("Failed to set new func: func not the same", err)
	}
}

func Test_Task_SetDoFuncFailure(t *testing.T) {
	foo := outer("hello")

	task, err := Create("printing", time.Second*3, time.Second*3, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	err = task.SetDoFunc(nil)

	if err == nil {
		t.Error("Failed to detect error while setting new delay: ", err)
	}
}
