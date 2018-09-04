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
	worker := Worker{}
	foo := outer("hello")

	a, err := tk.Create("printing", time.Second*2, time.Second*3, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	pos, err := worker.Add(a)

	if err != nil {
		t.Error("Failed to add task to worker: ", err)
	}

	if err := worker.Start(pos); err != nil {
		t.Error("Failed to start worker: ", err)
	}

	time.Sleep(4 * time.Second)

	if err := worker.Stop(pos); err != nil {
		t.Error("Failed to start worker: ", err)
	}

	time.Sleep(2 * time.Second)
}

func Test_Worker_StartExpired(t *testing.T) {
	fmt.Println("//Test_Worker_StartExpired//")
	worker := Worker{}
	foo := outer("hello")

	a, err := tk.Create("printing", time.Second*2, time.Millisecond*1, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	pos, err := worker.Add(a)

	if err != nil {
		t.Error("Failed to add task to worker: ", err)
	}

	if err := worker.Start(pos); err != nil {
		t.Error("Failed to start worker: ", err)
	}

	time.Sleep(4 * time.Second)

	if err := worker.Stop(pos); err == nil {
		t.Error("Failed to detect expired:")
	}

	time.Sleep(2 * time.Second)
}

func Test_Worker_StartJobFinished(t *testing.T) {
	fmt.Println("//Test_Worker_StartJobFinished//")
	worker := Worker{}
	foo := outer("hello")

	a, err := tk.Create("printing", time.Second*6, time.Second*6, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	pos, err := worker.Add(a)

	if err != nil {
		t.Error("Failed to add task to worker: ", err)
	}

	if err := worker.Start(pos); err != nil {
		t.Error("Failed to start worker: ", err)
	}

	time.Sleep(8 * time.Second)

	if err := worker.Stop(pos); err != nil {
		t.Error("Failed to detect expired:")
	}

	time.Sleep(2 * time.Second)
}

func Test_Worker_StartDouble(t *testing.T) {
	fmt.Println("//Test_Worker_StartDouble//")
	worker := Worker{}
	foo := outer("hello")

	a, err := tk.Create("printing", time.Second*2, time.Second*3, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	pos, err := worker.Add(a)

	if err != nil {
		t.Error("Failed to add task to worker: ", err)
	}

	if err := worker.Start(pos); err != nil {
		t.Error("Failed to start worker: ", err)
	}

	if err := worker.Start(pos); err == nil {
		t.Error("Failed to detect error while starting work")
	}

	if err := worker.Stop(pos); err != nil {
		t.Error("Failed to start worker: ", err)
	}

	time.Sleep(2 * time.Second)
}

func Test_Worker_StartError(t *testing.T) {
	fmt.Println("//Test_Worker_StartError//")
	worker := Worker{}
	foo := outer("hello")

	a, err := tk.Create("printing", time.Second*2, time.Second*3, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	pos, err := worker.Add(a)

	if err != nil {
		t.Error("Failed to add task to worker: ", err)
	}

	if err := worker.Start(pos + 1); err == nil {
		t.Error("Failed to detect error while starting worker: ", err)
	}

	time.Sleep(2 * time.Second)
}

func Test_Worker_StopDouble(t *testing.T) {
	fmt.Println("//Test_Worker_StopDouble//")
	worker := Worker{}
	foo := outer("hello")

	a, err := tk.Create("printing", time.Second*2, time.Second*3, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	pos, err := worker.Add(a)

	if err != nil {
		t.Error("Failed to add task to worker: ", err)
	}

	if err := worker.Start(pos); err != nil {
		t.Error("Failed to start worker: ", err)
	}

	if err := worker.Stop(pos); err != nil {
		t.Error("Failed to start worker: ", err)
	}

	if err := worker.Stop(pos); err == nil {
		t.Error("Failed to detect error while stoping worker: ", err)
	}

	time.Sleep(2 * time.Second)
}

func Test_Worker_StopError(t *testing.T) {
	fmt.Println("//Test_Worker_StopError//")
	worker := Worker{}
	foo := outer("hello")

	a, err := tk.Create("printing", time.Second*2, time.Second*3, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	pos, err := worker.Add(a)

	if err != nil {
		t.Error("Failed to add task to worker: ", err)
	}

	if err := worker.Start(pos); err != nil {
		t.Error("Failed to start worker: ", err)
	}

	if err := worker.Stop(pos); err != nil {
		t.Error("Failed to start worker: ", err)
	}

	if err := worker.Stop(pos + 1); err == nil {
		t.Error("Failed to detect error while stoping worker: ", err)
	}

	time.Sleep(2 * time.Second)
}

func Test_Worker_StartAndStopWithoutDelay(t *testing.T) {
	fmt.Println("//Test_Worker_StartAndStopWithoutDelay//")
	worker := Worker{}
	foo := outer("hello")

	a, err := tk.Create("printing", time.Second*2, time.Second*3, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	pos, err := worker.Add(a)

	if err != nil {
		t.Error("Failed to add task to worker: ", err)
	}

	if err := worker.Start(pos); err != nil {
		t.Error("Failed to start worker: ", err)
	}

	if err := worker.Stop(pos); err != nil {
		t.Error("Failed to start worker: ", err)
	}

	time.Sleep(2 * time.Second)
}

func Test_Worker_StartAndKillWithDelay(t *testing.T) {
	fmt.Println("//Test_Worker_StartAndKillWithDelay//")
	worker := Worker{}
	foo := outer("hello")

	a, err := tk.Create("printing", time.Second*2, time.Second*3, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	pos, err := worker.Add(a)
	if err != nil {
		t.Error("Failed to add task to worker: ", err)
	}

	if err := worker.Start(pos); err != nil {
		t.Error("Failed to start worker: ", err)
	}

	time.Sleep(2 * time.Second)

	if err := worker.Kill(pos); err != nil {
		t.Error("Failed to start worker: ", err)
	}

	time.Sleep(1 * time.Second)
}

func Test_Worker_StartAndKillWithoutDelay(t *testing.T) {
	fmt.Println("//Test_Worker_StartAndKillWithoutDelay//")
	worker := Worker{}
	foo := outer("hello")

	a, err := tk.Create("printing", time.Second*2, time.Second*3, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	pos, err := worker.Add(a)
	if err != nil {
		t.Error("Failed to add task to worker: ", err)
	}

	if err := worker.Start(pos); err != nil {
		t.Error("Failed to start worker: ", err)
	}

	if err := worker.Kill(pos); err != nil {
		t.Error("Failed to start worker: ", err)
	}

	time.Sleep(1 * time.Second)
}

func Test_Worker_KillDouble(t *testing.T) {
	fmt.Println("//Test_Worker_KillDouble//")
	worker := Worker{}
	foo := outer("hello")

	a, err := tk.Create("printing", time.Second*2, time.Second*3, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	pos, err := worker.Add(a)
	if err != nil {
		t.Error("Failed to add task to worker: ", err)
	}

	if err := worker.Start(pos); err != nil {
		t.Error("Failed to start worker: ", err)
	}

	if err := worker.Kill(pos); err != nil {
		t.Error("Failed to start worker: ", err)
	}

	if err := worker.Kill(pos); err == nil {
		t.Error("Failed to detect error while killing worker: ", err)
	}

	time.Sleep(1 * time.Second)
}

func Test_Worker_KillError(t *testing.T) {
	fmt.Println("//Test_Worker_KillError//")
	worker := Worker{}
	foo := outer("hello")

	a, err := tk.Create("printing", time.Second*2, time.Second*3, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	pos, err := worker.Add(a)
	if err != nil {
		t.Error("Failed to add task to worker: ", err)
	}

	if err := worker.Start(pos); err != nil {
		t.Error("Failed to start worker: ", err)
	}

	if err := worker.Kill(pos); err != nil {
		t.Error("Failed to start worker: ", err)
	}

	if err := worker.Kill(pos + 1); err == nil {
		t.Error("Failed to detect error while killing worker: ", err)
	}

	time.Sleep(1 * time.Second)
}

func Test_Worker_StartAllAndStopAllWithDelay(t *testing.T) {
	fmt.Println("//Test_Worker_StartAllAndStopAllWithDelay//")
	worker := Worker{}
	foo := outer("hello")

	a, err := tk.Create("printing", time.Second*2, time.Second*3, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	b, err := tk.Create("printing2", time.Second*2, time.Second*3, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	_, err = worker.Add(a)
	if err != nil {
		t.Error("Failed to add task to worker: ", err)
	}

	_, err = worker.Add(b)
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
	worker := Worker{}
	foo := outer("hello")

	a, err := tk.Create("printing", time.Second*2, time.Second*3, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	b, err := tk.Create("printing2", time.Second*2, time.Second*3, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	_, err = worker.Add(a)
	if err != nil {
		t.Error("Failed to add task to worker: ", err)
	}

	_, err = worker.Add(b)
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
	worker := Worker{}
	foo := outer("hello")

	a, err := tk.Create("printing", time.Second*2, time.Second*3, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	b, err := tk.Create("printing2", time.Second*2, time.Second*3, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	_, err = worker.Add(a)
	if err != nil {
		t.Error("Failed to add task to worker: ", err)
	}

	_, err = worker.Add(b)
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
	worker := Worker{}
	foo := outer("hello")

	a, err := tk.Create("printing", time.Second*2, time.Second*3, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	b, err := tk.Create("printing2", time.Second*2, time.Second*3, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	_, err = worker.Add(a)
	if err != nil {
		t.Error("Failed to add task to worker: ", err)
	}

	_, err = worker.Add(b)
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
	worker := Worker{}
	foo := outer("hello")

	a, err := tk.Create("printing", time.Second*2, time.Second*3, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	_, err = worker.Add(a)

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
	worker := Worker{}
	foo := outer("hello")

	a, err := tk.Create("printing", time.Second*2, time.Second*3, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	_, err = worker.Add(a)

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
	worker := Worker{}
	foo := outer("hello")

	a, err := tk.Create("printing", time.Second*3, time.Second*3, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	pos, err := worker.Add(a)

	if pos != 0 {
		t.Error("Failed to add task to worker: position of added task is not OK")
	}

	if len(worker.jobs) != 1 {
		t.Error("Failed to add task to worker: length is not the same")
	}

	if err != nil {
		t.Error("Failed to add task to worker: ", err)
	}
}

func Test_Worker_AddWithError(t *testing.T) {
	worker := Worker{}
	foo := outer("hello")

	a, err := tk.Create("printing", time.Second*3, time.Second*3, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	_, err = worker.Add(a)
	_, err = worker.Add(a)

	if err == nil {
		t.Error("Failed to detect error while adding task to worker: ", err)
	}
}

func Test_Worker_Check(t *testing.T) {
	worker := Worker{}
	foo := outer("hello")

	a, err := tk.Create("printing", time.Second*3, time.Second*3, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	n, err := worker.Add(a)

	if err != nil {
		t.Error("Failed to add task to worker: ", err)
	}

	if !worker.check(n) {
		t.Error("Can't find added job")
	}

	if worker.check(n + 1) {
		t.Error("Find not added job")
	}
}

func Test_Worker_ChangeTask(t *testing.T) {
	worker := Worker{}
	foo := outer("hello")

	a, err := tk.Create("printing", time.Second*3, time.Second*3, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	n, err := worker.Add(a)

	if err != nil {
		t.Error("Failed to add task to worker: ", err)
	}

	a.SetName("printingChanged")
	a.SetPeriod(time.Second * 2)
	a.SetTaskTime(time.Second * 2)
	a.SetDelay(time.Second * 2)

	err = worker.ChangeTask(n, a)

	if err != nil {
		t.Error("Failed to change task: ", err)
	}

	if worker.jobs[n].task.GetName() != "printingChanged" || worker.jobs[n].task.GetPeriod() != time.Second*2 || worker.jobs[n].task.GetTaskTime() != time.Second*2 ||
		worker.jobs[n].task.GetDelay() != time.Second*2 {
		t.Error("Failed to change task: diferent values")
	}
}

func Test_Worker_ChangeTaskError(t *testing.T) {
	worker := Worker{}
	foo := outer("hello")

	a, err := tk.Create("printing", time.Second*3, time.Second*3, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	n, err := worker.Add(a)

	if err != nil {
		t.Error("Failed to add task to worker: ", err)
	}

	a.SetName("printingChanged")
	a.SetPeriod(time.Second * 2)
	a.SetTaskTime(time.Second * 2)
	a.SetDelay(time.Second * 2)

	err = worker.ChangeTask(n, a)

	err = worker.ChangeTask(n+1, a)

	if err == nil {
		t.Error("Failed to change task: ", err)
	}

	a.SetName("printingChanged")
	err = worker.ChangeTask(n, a)

	if err == nil {
		t.Error("Failed to change task: ", err)
	}

	worker.jobs[n].status = wtf
	a.SetName("printingChanged2")
	err = worker.ChangeTask(n, a)

	if err == nil {
		t.Error("Failed to change task: ", err)
	}
}

func Test_Worker_Print(t *testing.T) {
	worker := Worker{}
	foo := outer("hello")

	a, err := tk.Create("printing", time.Second*3, time.Second*3, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	n, err := worker.Add(a)

	if err != nil {
		t.Error("Failed to add task to worker: ", err)
	}
	if err := worker.Print(n); err != nil {
		t.Error("Failed to print: ", err)
	}
}

func Test_Worker_PrintError(t *testing.T) {
	worker := Worker{}
	foo := outer("hello")

	a, err := tk.Create("printing", time.Second*3, time.Second*3, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	n, err := worker.Add(a)

	if err != nil {
		t.Error("Failed to add task to worker: ", err)
	}
	if err := worker.Print(n + 1); err == nil {
		t.Error("Failed to detect error while printing")
	}
}

func Test_Worker_PrintAll(t *testing.T) {
	worker := Worker{}
	foo := outer("hello")

	a, err := tk.Create("printing", time.Second*3, time.Second*3, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	_, err = worker.Add(a)

	if err != nil {
		t.Error("Failed to add task to worker: ", err)
	}
	if err := worker.PrintAll(); err != nil {
		t.Error("Failed to print all: ", err)
	}
}

func Test_Worker_GetNumberByName(t *testing.T) {
	worker := Worker{}
	foo := outer("hello")

	a, err := tk.Create("printing", time.Second*3, time.Second*3, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	n, err := worker.Add(a)

	if err != nil {
		t.Error("Failed to add task to worker: ", err)
	}

	num, err := worker.GetNumberByName("printing")
	if num != n {
		t.Error("Different numbers")
	}
	if err != nil {
		t.Error("Failed to get number by namme: ", err)
	}
}

func Test_Worker_GetNumberByNameError(t *testing.T) {
	worker := Worker{}
	foo := outer("hello")

	a, err := tk.Create("printing", time.Second*3, time.Second*3, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	_, err = worker.Add(a)

	if err != nil {
		t.Error("Failed to add task to worker: ", err)
	}

	num, err := worker.GetNumberByName("printin")

	if err == nil && num != -100 {
		t.Error("Failed to detect error while getting number by name")
	}
}

func Test_Worker_Delete(t *testing.T) {
	worker := Worker{}
	foo := outer("hello")

	a, err := tk.Create("printing", time.Second*3, time.Second*3, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	n, err := worker.Add(a)

	if err != nil {
		t.Error("Failed to add task to worker: ", err)
	}

	prevLen := len(worker.jobs)
	err = worker.Delete(n)
	if err != nil {
		t.Error("Failed to delete: ", err)
	}
	if len(worker.jobs) != prevLen-1 {
		t.Error("Failed to delete: size not the same")
	}
}

func Test_Worker_DeleteError(t *testing.T) {
	worker := Worker{}
	foo := outer("hello")

	a, err := tk.Create("printing", time.Second*3, time.Second*3, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	n, err := worker.Add(a)

	if err != nil {
		t.Error("Failed to add task to worker: ", err)
	}

	err = worker.Delete(n + 1)
	if err == nil {
		t.Error("Failed to detect error while deleting")
	}

	worker.jobs[n].status = wtf
	err = worker.Delete(n)
	if err == nil {
		t.Error("Failed to detect error while deleting")
	}
}

func Test_Worker_DeleteKilled(t *testing.T) {
	worker := Worker{}
	foo := outer("hello")

	a, err := tk.Create("printing", time.Second*3, time.Second*3, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	n, err := worker.Add(a)

	if err != nil {
		t.Error("Failed to add task to worker: ", err)
	}

	prevLen := len(worker.jobs)

	err = worker.deleteKilled(n)
	if err != nil {
		t.Error("Failed to delete killed: ", err)
	}
	if len(worker.jobs) != prevLen {
		t.Error("Deleted not killd")
	}

	worker.jobs[n].status = k
	err = worker.deleteKilled(n)
	if err != nil {
		t.Error("Failed to delete killed: ", err)
	}
	if len(worker.jobs) != prevLen-1 {
		t.Error("Failed to delete: size not the same")
	}
}

func Test_Worker_DeleteKilledError(t *testing.T) {
	worker := Worker{}
	foo := outer("hello")

	a, err := tk.Create("printing", time.Second*3, time.Second*3, time.Second*1, foo)

	if err != nil {
		t.Error("Failed to create task: ", err)
	}

	n, err := worker.Add(a)

	if err != nil {
		t.Error("Failed to add task to worker: ", err)
	}

	err = worker.deleteKilled(n + 1)
	if err == nil {
		t.Error("Failed to detect error while deleting killed: ", err)
	}
}
