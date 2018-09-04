package worker

import (
	"context"
	"fmt"
	"testing"
	"time"

	tk "github.com/vslchnk/goscheduler/task"
)

func outer(name string) func(ctx context.Context) error {
	text := "Modified " + name

	foo := func(ctx context.Context) error {
		defer ctx.Value("func").(context.CancelFunc)()

		c2, cancel := context.WithCancel(ctx)
		defer cancel()

		for i := 0; i < 3; i++ {
			time.Sleep(time.Second * 1)
			select {
			case <-c2.Done():
				fmt.Println("Done")
				return nil
			default:
				fmt.Println(text)
			}
		}
		fmt.Println(ctx.Err() == context.Canceled)

		return nil
	}

	return foo
}

func Test_Worker_StartAndStopWithDelay(t *testing.T) {
	fmt.Println("//Test_Worker_StartAndStopWithDelay//")
	worker := NewWorker()
	foo := outer("hello")

	a, err := tk.Create(time.Second*5, time.Second*5, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	name := "printing"

	err = worker.Add(a, name)

	if err != nil {
		t.Error("Failed to add task to worker: ", err)
	}

	if err := worker.Start(name); err != nil {
		t.Error("Failed to start worker: ", err)
	}

	time.Sleep(3 * time.Second)

	if err := worker.Stop(name); err != nil {
		t.Error("Failed to stop worker: ", err)
	}

	time.Sleep(2 * time.Second)
}

func Test_Worker_StartAndStopAfterDelayBeforeTick(t *testing.T) {
	fmt.Println("//Test_Worker_StartAndStopAfterDelayBeforeTick//")
	worker := NewWorker()
	foo := outer("hello")

	a, err := tk.Create(time.Second*5, time.Second*5, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	name := "printing"

	err = worker.Add(a, name)

	if err != nil {
		t.Error("Failed to add task to worker: ", err)
	}

	if err := worker.Start(name); err != nil {
		t.Error("Failed to start worker: ", err)
	}

	time.Sleep(time.Millisecond * 1000)

	if err := worker.Stop(name); err != nil {
		t.Error("Failed to stop worker: ", err)
	}

	time.Sleep(2 * time.Second)
}

func Test_Worker_StartExpired(t *testing.T) {
	fmt.Println("//Test_Worker_StartExpired//")
	worker := NewWorker()
	foo := outer("hello")

	a, err := tk.Create(time.Second*2, time.Millisecond*1, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	name := "printing"

	err = worker.Add(a, name)

	if err != nil {
		t.Error("Failed to add task to worker: ", err)
	}

	if err := worker.Start(name); err != nil {
		t.Error("Failed to start worker: ", err)
	}

	time.Sleep(4 * time.Second)

	if err := worker.Stop(name); err == nil {
		t.Error("Failed to detect expired:")
	}

	time.Sleep(2 * time.Second)
}

func Test_Worker_StartJobFinished(t *testing.T) {
	fmt.Println("//Test_Worker_StartJobFinished//")
	worker := NewWorker()
	name := "printing"
	foo := outer("hello")

	a, err := tk.Create(time.Second*6, time.Second*6, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	err = worker.Add(a, "printing")

	if err != nil {
		t.Error("Failed to add task to worker: ", err)
	}

	if err := worker.Start(name); err != nil {
		t.Error("Failed to start worker: ", err)
	}

	time.Sleep(8 * time.Second)

	if err := worker.Stop(name); err != nil {
		t.Error("Failed to detect expired:")
	}

	time.Sleep(2 * time.Second)
}

func Test_Worker_StartDouble(t *testing.T) {
	fmt.Println("//Test_Worker_StartDouble//")
	worker := NewWorker()
	name := "printing"
	foo := outer("hello")

	a, err := tk.Create(time.Second*5, time.Second*5, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	err = worker.Add(a, "printing")

	if err != nil {
		t.Error("Failed to add task to worker: ", err)
	}

	if err := worker.Start(name); err != nil {
		t.Error("Failed to start worker: ", err)
	}

	if err := worker.Start(name); err == nil {
		t.Error("Failed to detect error while starting work")
	}

	if err := worker.Stop(name); err != nil {
		t.Error("Failed to start worker: ", err)
	}

	time.Sleep(2 * time.Second)
}

func Test_Worker_StartError(t *testing.T) {
	fmt.Println("//Test_Worker_StartError//")
	worker := NewWorker()
	name := "printing"
	foo := outer("hello")

	a, err := tk.Create(time.Second*5, time.Second*5, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	err = worker.Add(a, "printing")

	if err != nil {
		t.Error("Failed to add task to worker: ", err)
	}

	if err := worker.Start(name + "2"); err == nil {
		t.Error("Failed to detect error while starting worker: ", err)
	}

	time.Sleep(2 * time.Second)
}

func Test_Worker_StopDouble(t *testing.T) {
	fmt.Println("//Test_Worker_StopDouble//")
	worker := NewWorker()
	name := "printing"
	foo := outer("hello")

	a, err := tk.Create(time.Second*5, time.Second*5, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	err = worker.Add(a, "printing")

	if err != nil {
		t.Error("Failed to add task to worker: ", err)
	}

	if err := worker.Start(name); err != nil {
		t.Error("Failed to start worker: ", err)
	}

	if err := worker.Stop(name); err != nil {
		t.Error("Failed to start worker: ", err)
	}

	if err := worker.Stop(name); err == nil {
		t.Error("Failed to detect error while stoping worker: ", err)
	}

	time.Sleep(2 * time.Second)
}

func Test_Worker_StopError(t *testing.T) {
	fmt.Println("//Test_Worker_StopError//")
	worker := NewWorker()
	name := "printing"
	foo := outer("hello")

	a, err := tk.Create(time.Second*5, time.Second*5, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	err = worker.Add(a, "printing")

	if err != nil {
		t.Error("Failed to add task to worker: ", err)
	}

	if err := worker.Start(name); err != nil {
		t.Error("Failed to start worker: ", err)
	}

	if err := worker.Stop(name); err != nil {
		t.Error("Failed to start worker: ", err)
	}

	if err := worker.Stop(name + "2"); err == nil {
		t.Error("Failed to detect error while stoping worker: ", err)
	}

	time.Sleep(2 * time.Second)
}

func Test_Worker_StartAndStopWithoutDelay(t *testing.T) {
	fmt.Println("//Test_Worker_StartAndStopWithoutDelay//")
	worker := NewWorker()
	name := "printing"
	foo := outer("hello")

	a, err := tk.Create(time.Second*5, time.Second*5, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	err = worker.Add(a, "printing")

	if err != nil {
		t.Error("Failed to add task to worker: ", err)
	}

	if err := worker.Start(name); err != nil {
		t.Error("Failed to start worker: ", err)
	}

	if err := worker.Stop(name); err != nil {
		t.Error("Failed to start worker: ", err)
	}

	time.Sleep(2 * time.Second)
}

func Test_Worker_StartAndKillWithDelay(t *testing.T) {
	fmt.Println("//Test_Worker_StartAndKillWithDelay//")
	worker := NewWorker()
	name := "printing"
	foo := outer("hello")

	a, err := tk.Create(time.Second*5, time.Second*5, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	err = worker.Add(a, "printing")
	if err != nil {
		t.Error("Failed to add task to worker: ", err)
	}

	if err := worker.Start(name); err != nil {
		t.Error("Failed to start worker: ", err)
	}

	time.Sleep(2 * time.Second)

	if err := worker.Kill(name); err != nil {
		t.Error("Failed to start worker: ", err)
	}

	time.Sleep(1 * time.Second)
}

func Test_Worker_StartAndKillWithoutDelay(t *testing.T) {
	fmt.Println("//Test_Worker_StartAndKillWithoutDelay//")
	worker := NewWorker()
	name := "printing"
	foo := outer("hello")

	a, err := tk.Create(time.Second*5, time.Second*5, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	err = worker.Add(a, "printing")
	if err != nil {
		t.Error("Failed to add task to worker: ", err)
	}

	if err := worker.Start(name); err != nil {
		t.Error("Failed to start worker: ", err)
	}

	if err := worker.Kill(name); err != nil {
		t.Error("Failed to start worker: ", err)
	}

	time.Sleep(1 * time.Second)
}

func Test_Worker_KillDouble(t *testing.T) {
	fmt.Println("//Test_Worker_KillDouble//")
	worker := NewWorker()
	name := "printing"
	foo := outer("hello")

	a, err := tk.Create(time.Second*5, time.Second*5, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	err = worker.Add(a, "printing")
	if err != nil {
		t.Error("Failed to add task to worker: ", err)
	}

	if err := worker.Start(name); err != nil {
		t.Error("Failed to start worker: ", err)
	}

	if err := worker.Kill(name); err != nil {
		t.Error("Failed to start worker: ", err)
	}

	if err := worker.Kill(name); err == nil {
		t.Error("Failed to detect error while killing worker: ", err)
	}

	time.Sleep(1 * time.Second)
}

func Test_Worker_KillError(t *testing.T) {
	fmt.Println("//Test_Worker_KillError//")
	worker := NewWorker()
	name := "printing"
	foo := outer("hello")

	a, err := tk.Create(time.Second*5, time.Second*5, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	err = worker.Add(a, "printing")
	if err != nil {
		t.Error("Failed to add task to worker: ", err)
	}

	if err := worker.Start(name); err != nil {
		t.Error("Failed to start worker: ", err)
	}

	if err := worker.Kill(name); err != nil {
		t.Error("Failed to start worker: ", err)
	}

	if err := worker.Kill(name + "1"); err == nil {
		t.Error("Failed to detect error while killing worker: ", err)
	}

	time.Sleep(1 * time.Second)
}

func Test_Worker_StartAllAndStopAllWithDelay(t *testing.T) {
	fmt.Println("//Test_Worker_StartAllAndStopAllWithDelay//")
	worker := NewWorker()
	name := "printing"
	name2 := "printing2"
	foo := outer("hello")

	a, err := tk.Create(time.Second*5, time.Second*5, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	b, err := tk.Create(time.Second*5, time.Second*5, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	err = worker.Add(a, name)
	if err != nil {
		t.Error("Failed to add task to worker: ", err)
	}

	err = worker.Add(b, name2)
	if err != nil {
		t.Error("Failed to add task to worker: ", err)
	}

	if err := worker.StartAll(); err != nil {
		t.Error("Failed to start worker: ", err)
	}

	time.Sleep(4 * time.Second)

	if err := worker.StopAll(); err != nil {
		t.Error("Failed to start worker: ", err)
	}

	time.Sleep(2 * time.Second)
}

func Test_Worker_StartAllAndKillAllWithDelay(t *testing.T) {
	fmt.Println("//Test_Worker_StartAllAndKillAllWithDelay//")
	worker := NewWorker()
	name := "printing"
	name2 := "printing2"
	foo := outer("hello")

	a, err := tk.Create(time.Second*5, time.Second*5, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	b, err := tk.Create(time.Second*5, time.Second*5, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	err = worker.Add(a, name)
	if err != nil {
		t.Error("Failed to add task to worker: ", err)
	}

	err = worker.Add(b, name2)
	if err != nil {
		t.Error("Failed to add task to worker: ", err)
	}

	if err := worker.StartAll(); err != nil {
		t.Error("Failed to start worker: ", err)
	}

	time.Sleep(2 * time.Second)

	if err := worker.KillAll(); err != nil {
		t.Error("Failed to start worker: ", err)
	}

	time.Sleep(1 * time.Second)
}

func Test_Worker_StartAllAndKillAllWithoutDelay(t *testing.T) {
	fmt.Println("//Test_Worker_StartAllAndKillAllWithoutDelay//")
	worker := NewWorker()
	name := "printing"
	name2 := "printing2"
	foo := outer("hello")

	a, err := tk.Create(time.Second*5, time.Second*5, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	b, err := tk.Create(time.Second*5, time.Second*5, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	err = worker.Add(a, name)
	if err != nil {
		t.Error("Failed to add task to worker: ", err)
	}

	err = worker.Add(b, name2)
	if err != nil {
		t.Error("Failed to add task to worker: ", err)
	}

	if err := worker.StartAll(); err != nil {
		t.Error("Failed to start worker: ", err)
	}

	if err := worker.KillAll(); err != nil {
		t.Error("Failed to start worker: ", err)
	}

	time.Sleep(1 * time.Second)
}

func Test_Worker_KillAllDouble(t *testing.T) {
	fmt.Println("//Test_Worker_KillAllDouble//")
	worker := NewWorker()
	name := "printing"
	name2 := "printing2"
	foo := outer("hello")

	a, err := tk.Create(time.Second*5, time.Second*5, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	b, err := tk.Create(time.Second*5, time.Second*5, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	err = worker.Add(a, name)
	if err != nil {
		t.Error("Failed to add task to worker: ", err)
	}

	err = worker.Add(b, name2)
	if err != nil {
		t.Error("Failed to add task to worker: ", err)
	}

	if err := worker.StartAll(); err != nil {
		t.Error("Failed to start worker: ", err)
	}

	if err := worker.KillAll(); err != nil {
		t.Error("Failed to start worker: ", err)
	}

	if err := worker.KillAll(); err == nil {
		t.Error("Failed to detect error while killing worker: ", err)
	}

	time.Sleep(1 * time.Second)
}

func Test_Worker_StartAllDouble(t *testing.T) {
	fmt.Println("//Test_Worker_StarAllDouble//")
	worker := NewWorker()
	name := "printing"
	foo := outer("hello")

	a, err := tk.Create(time.Second*5, time.Second*5, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	err = worker.Add(a, name)

	if err != nil {
		t.Error("Failed to add task to worker: ", err)
	}

	if err := worker.StartAll(); err != nil {
		t.Error("Failed to start worker: ", err)
	}

	if err := worker.StartAll(); err == nil {
		t.Error("Failed to detect error while starting worker: ", err)
	}

	if err := worker.StopAll(); err != nil {
		t.Error("Failed to stop worker: ", err)
	}

	time.Sleep(2 * time.Second)
}

func Test_Worker_StopAllDouble(t *testing.T) {
	fmt.Println("//Test_Worker_StopAllDouble//")
	worker := NewWorker()
	name := "printing"
	foo := outer("hello")

	a, err := tk.Create(time.Second*5, time.Second*5, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	err = worker.Add(a, name)

	if err != nil {
		t.Error("Failed to add task to worker: ", err)
	}

	if err := worker.StartAll(); err != nil {
		t.Error("Failed to start worker: ", err)
	}

	if err := worker.StopAll(); err != nil {
		t.Error("Failed to start worker: ", err)
	}

	if err := worker.StopAll(); err == nil {
		t.Error("Failed to detect error while stoping worker: ", err)
	}

	time.Sleep(2 * time.Second)
}

func Test_Worker_Add(t *testing.T) {
	worker := NewWorker()
	name := "printing"
	foo := outer("hello")

	a, err := tk.Create(time.Second*3, time.Second*3, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	err = worker.Add(a, name)

	if len(worker.jobs) != 1 {
		t.Error("Failed to add task to worker: length is not the same")
	}

	if err != nil {
		t.Error("Failed to add task to worker: ", err)
	}
}

func Test_Worker_AddWithError(t *testing.T) {
	worker := NewWorker()
	name := "printing"
	foo := outer("hello")

	a, err := tk.Create(time.Second*3, time.Second*3, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	err = worker.Add(a, name)
	err = worker.Add(a, name)

	if err == nil {
		t.Error("Failed to detect error while adding task to worker: ", err)
	}
}

func Test_Worker_Check(t *testing.T) {
	worker := NewWorker()
	name := "printing"
	foo := outer("hello")

	a, err := tk.Create(time.Second*3, time.Second*3, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	err = worker.Add(a, name)

	if err != nil {
		t.Error("Failed to add task to worker: ", err)
	}

	if !worker.check(name) {
		t.Error("Can't find added job")
	}

	if worker.check(name + "1") {
		t.Error("Find not added job")
	}
}

func Test_Worker_ChangeTask(t *testing.T) {
	worker := NewWorker()
	name := "printing"
	foo := outer("hello")

	a, err := tk.Create(time.Second*3, time.Second*3, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	err = worker.Add(a, name)

	if err != nil {
		t.Error("Failed to add task to worker: ", err)
	}

	a.SetPeriod(time.Second * 2)
	a.SetTaskTime(time.Second * 2)
	a.SetDelay(time.Second * 2)

	err = worker.ChangeTask(name, a)

	if err != nil {
		t.Error("Failed to change task: ", err)
	}

	if worker.jobs[name].task.GetPeriod() != time.Second*2 || worker.jobs[name].task.GetTaskTime() != time.Second*2 || worker.jobs[name].task.GetDelay() != time.Second*2 {
		t.Error("Failed to change task: diferent values")
	}
}

func Test_Worker_ChangeTaskError(t *testing.T) {
	worker := NewWorker()
	name := "printing"
	foo := outer("hello")

	a, err := tk.Create(time.Second*3, time.Second*3, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	err = worker.Add(a, name)

	if err != nil {
		t.Error("Failed to add task to worker: ", err)
	}

	a.SetPeriod(time.Second * 2)
	a.SetTaskTime(time.Second * 2)
	a.SetDelay(time.Second * 2)

	err = worker.ChangeTask(name, a)

	err = worker.ChangeTask(name+"1", a)

	if err == nil {
		t.Error("Failed to change task: ", err)
	}

	worker.jobs[name].status = wtf

	err = worker.ChangeTask(name, a)

	if err == nil {
		t.Error("Failed to change task: ", err)
	}
}

func Test_Worker_Print(t *testing.T) {
	worker := NewWorker()
	name := "printing"
	foo := outer("hello")

	a, err := tk.Create(time.Second*3, time.Second*3, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	err = worker.Add(a, name)

	if err != nil {
		t.Error("Failed to add task to worker: ", err)
	}
	if err := worker.Print(name); err != nil {
		t.Error("Failed to print: ", err)
	}
}

func Test_Worker_PrintError(t *testing.T) {
	worker := NewWorker()
	name := "printing"
	foo := outer("hello")

	a, err := tk.Create(time.Second*3, time.Second*3, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	err = worker.Add(a, name)

	if err != nil {
		t.Error("Failed to add task to worker: ", err)
	}
	if err := worker.Print(name + "1"); err == nil {
		t.Error("Failed to detect error while printing")
	}
}

func Test_Worker_PrintAll(t *testing.T) {
	worker := NewWorker()
	name := "printing"
	foo := outer("hello")

	a, err := tk.Create(time.Second*3, time.Second*3, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	err = worker.Add(a, name)

	if err != nil {
		t.Error("Failed to add task to worker: ", err)
	}
	if err := worker.PrintAll(); err != nil {
		t.Error("Failed to print all: ", err)
	}
}

func Test_Worker_Delete(t *testing.T) {
	worker := NewWorker()
	name := "printing"
	foo := outer("hello")

	a, err := tk.Create(time.Second*3, time.Second*3, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	err = worker.Add(a, name)

	if err != nil {
		t.Error("Failed to add task to worker: ", err)
	}

	prevLen := len(worker.jobs)
	err = worker.Delete(name)
	if err != nil {
		t.Error("Failed to delete: ", err)
	}
	if len(worker.jobs) != prevLen-1 {
		t.Error("Failed to delete: size not the same")
	}
}

func Test_Worker_DeleteError(t *testing.T) {
	worker := NewWorker()
	name := "printing"
	foo := outer("hello")

	a, err := tk.Create(time.Second*3, time.Second*3, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	err = worker.Add(a, name)

	if err != nil {
		t.Error("Failed to add task to worker: ", err)
	}

	err = worker.Delete(name + "1")
	if err == nil {
		t.Error("Failed to detect error while deleting")
	}

	worker.jobs[name].status = wtf
	err = worker.Delete(name)
	if err == nil {
		t.Error("Failed to detect error while deleting")
	}
}

func Test_Worker_DeleteKilled(t *testing.T) {
	worker := NewWorker()
	name := "printing"
	foo := outer("hello")

	a, err := tk.Create(time.Second*3, time.Second*3, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	err = worker.Add(a, name)

	if err != nil {
		t.Error("Failed to add task to worker: ", err)
	}

	prevLen := len(worker.jobs)

	err = worker.deleteKilled(name)
	if err != nil {
		t.Error("Failed to delete killed: ", err)
	}
	if len(worker.jobs) != prevLen {
		t.Error("Deleted not killd")
	}

	worker.jobs[name].status = k
	err = worker.deleteKilled(name)
	if err != nil {
		t.Error("Failed to delete killed: ", err)
	}
	if len(worker.jobs) != prevLen-1 {
		t.Error("Failed to delete: size not the same")
	}
}

func Test_Worker_DeleteKilledError(t *testing.T) {
	worker := NewWorker()
	name := "printing"
	foo := outer("hello")

	a, err := tk.Create(time.Second*3, time.Second*3, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	err = worker.Add(a, name)

	if err != nil {
		t.Error("Failed to add task to worker: ", err)
	}

	err = worker.deleteKilled(name + "1")
	if err == nil {
		t.Error("Failed to detect error while deleting killed: ", err)
	}
}
