package worker

import (
	"context"
	"fmt"
	"sync"
	"time"

	t "github.com/vslchnk/goscheduler/task"
)

type jobstatus string

const (
	sss  jobstatus = "stopped by stop signal"
	cnw  jobstatus = "created, not working"
	wtnf jobstatus = "working, task isn't finished"
	wtf  jobstatus = "working, task is finished"
	ste  jobstatus = "stopped, time has expired"
	k    jobstatus = "killed"
)

type Worker struct {
	sync.Mutex
	jobs []*job
}

type job struct {
	sync.Mutex
	task      t.Task
	ctx       context.Context
	cancelCtx context.CancelFunc
	status    jobstatus
}

// checks if there is a job with number n in job pool
func (w *Worker) check(n int) bool {
	defer w.Unlock()
	w.Lock()

	if n < 0 || n > len(w.jobs)-1 {
		return false
	}

	return true
}

// change task in job pool by its number
func (w *Worker) ChangeTask(n int, task t.Task) error {
	if !w.check(n) {
		return fmt.Errorf("No job with number %v in job pool", n)
	}

	defer w.jobs[n].Unlock()
	w.jobs[n].Lock()

	if w.jobs[n].status == wtnf || w.jobs[n].status == wtf {
		return fmt.Errorf("Job is working, must be stopped before being changed")
	}

	w.jobs[n].task.SetName(task.GetName())
	if err := w.jobs[n].task.SetPeriod(task.GetPeriod()); err != nil {
		return fmt.Errorf("Error in ChangeTask(): %v", err)
	}
	if err := w.jobs[n].task.SetTaskTime(task.GetTaskTime()); err != nil {
		return fmt.Errorf("Error in ChangeTask(): %v", err)
	}
	if err := w.jobs[n].task.SetDelay(task.GetDelay()); err != nil {
		return fmt.Errorf("Error in ChangeTask(): %v", err)
	}
	w.jobs[n].task.SetDoFunc(task.GetDoFunc())

	return nil
}

// adds task to job pool, if name of task is unique, if ok return number of task in job pool, if not return number of task with the same name and error
func (w *Worker) Add(task t.Task) (int, error) {
	for k, v := range w.jobs {
		if v.task.GetName() == task.GetName() {
			return k, fmt.Errorf("Function with name %v already exist", task.GetName())
		}
	}

	j := &job{task: task, status: cnw}

	w.Lock()
	w.jobs = append(w.jobs, j)
	w.Unlock()

	return len(w.jobs) - 1, nil
}

// prints jobs in job pool
func (w Worker) PrintAll() {
	for k, _ := range w.jobs {
		if err := w.Print(k); err != nil {
			fmt.Errorf("Error in PrintAll(): %v", err)
		}
	}
}

// prints job from job pool by its number
func (w Worker) Print(n int) error {
	if !w.check(n) {
		return fmt.Errorf("No job with number %v in job pool", n)
	}

	fmt.Printf("Number: %v; Status: %v; ", n, w.jobs[n].status)
	w.jobs[n].task.Print()

	return nil
}

// returns job number in job pool by task's name
func (w Worker) GetNumberByName(name string) (int, error) {
	for k, v := range w.jobs {
		if v.task.GetName() == name {
			return k, nil
		}
	}

	return -100, fmt.Errorf("No task with such name")
}

// starts job by its number
func (w *Worker) Start(n int) error {
	if !w.check(n) {
		return fmt.Errorf("No job with number %v in job pool", n)
	}

	w.jobs[n].Lock()

	if w.jobs[n].status == wtnf || w.jobs[n].status == wtf {
		w.jobs[n].Unlock()
		return fmt.Errorf("Job number %v is already working, its status: %v", n, w.jobs[n].status)
	}
	fmt.Println("worker started")
	ctx, cancel := context.WithCancel(context.Background())

	w.jobs[n].status = wtnf
	w.jobs[n].ctx = ctx
	w.jobs[n].cancelCtx = cancel
	w.jobs[n].Unlock()

	go w.startJob(n)

	return nil
}

// starts all jobs if they are stopped or not started
func (w *Worker) StartAll() error {
	for k, _ := range w.jobs {
		if err := w.Start(k); err != nil {
			return fmt.Errorf("Error in StartAll(): %v", err)
		}
	}

	return nil
}

// stops job by its number
func (w *Worker) Stop(n int) error {
	if !w.check(n) {
		return fmt.Errorf("No job with number %v in job pool", n)
	}

	w.jobs[n].Lock()
	if w.jobs[n].status == sss || w.jobs[n].status == ste || w.jobs[n].status == k || w.jobs[n].status == cnw {
		w.jobs[n].Unlock()
		return fmt.Errorf("Job number %v is not working, its status: %v", n, w.jobs[n].status)
	}
	fmt.Println("Stop", w.jobs[n].status)
	w.jobs[n].status = sss
	/*if w.jobs[n].cancelCtx == nil {
		time.Sleep(time.Nanosecond * 10000) // in case of when startJob haven't started
	}*/
	w.jobs[n].cancelCtx()
	w.jobs[n].Unlock()

	return nil
}

// stops all jobs if they are not stopped or not started
func (w *Worker) StopAll() error {
	for k, _ := range w.jobs {
		if err := w.Stop(k); err != nil {
			return fmt.Errorf("Error in StopAll(): %v", err)
		}
	}

	return nil
}

// kills job and removes it from pool
func (w *Worker) Kill(n int) error {
	if !w.check(n) {
		return fmt.Errorf("No job with number %v in job pool", n)
	}

	w.jobs[n].Lock()

	if w.jobs[n].status == sss || w.jobs[n].status == ste || w.jobs[n].status == k || w.jobs[n].status == cnw {
		w.jobs[n].Unlock()
		return fmt.Errorf("Job number %v is not working, its status: %v", n, w.jobs[n].status)
	}
	//fmt.Println("Stop", w.jobs[n].status)
	w.jobs[n].status = k
	/*if w.jobs[n].cancelCtx == nil {
		time.Sleep(time.Nanosecond * 10000) // in case of when startJob haven't started
	}*/
	w.jobs[n].cancelCtx()
	w.jobs[n].Unlock()

	return nil
}

// kills all jobs and removes them from pool
func (w *Worker) KillAll() error {
	for k, _ := range w.jobs {
		if err := w.Kill(k); err != nil {
			return fmt.Errorf("Error in KillAll(): %v", err)
		}
	}

	return nil
}

// delets job from pool by its number
func (w *Worker) Delete(n int) error {
	if !w.check(n) {
		return fmt.Errorf("No job with number %v in job pool", n)
	}

	w.Lock()
	if w.jobs[n].status == wtf || w.jobs[n].status == wtnf {
		w.jobs[n].Unlock()
		return fmt.Errorf("Job number %v is working, its status: %v", n, w.jobs[n].status)
	}

	copy(w.jobs[n:], w.jobs[n+1:])
	w.jobs[len(w.jobs)-1] = nil
	w.jobs = w.jobs[:len(w.jobs)-1]
	w.Unlock()

	return nil
}

// controls work of job
func (w *Worker) startJob(n int) {
	defer w.deleteKilled(n)
	defer w.jobs[n].cancelCtx()

	delayChan := time.NewTimer(time.Duration(w.jobs[n].task.GetDelay())).C
	for {
		select {
		case <-w.jobs[n].ctx.Done():
			fmt.Println("Stopped with context, exiting in 500 milliseconds")
			time.Sleep(500 * time.Millisecond)

			return
		case <-delayChan:
			fmt.Println("Delay is done")

			tickChan := time.NewTimer(0).C
			for {
				select {
				case <-w.jobs[n].ctx.Done():
					fmt.Println("Stopped with context, exiting in 500 milliseconds")
					time.Sleep(500 * time.Millisecond)

					return
				case <-tickChan:
					fmt.Println("Ticker ticked")
					//c2, cancel := context.WithCancel(w.jobs[n].ctx)
					c2, cancel := context.WithCancel(context.Background())
					c1 := context.WithValue(c2, "func", cancel)
					fmt.Println(w.jobs[n].status)
					//if w.jobs[n].status != sss {
					go w.jobs[n].task.GetDoFunc()(c1)
					//}

					tickChan = time.NewTimer(time.Duration(w.jobs[n].task.GetPeriod())).C
					expiredChan := time.NewTimer(time.Duration(w.jobs[n].task.GetTaskTime())).C
				Looptick:
					for {
						select {
						case <-w.jobs[n].ctx.Done():

							w.jobs[n].Lock()
							if w.jobs[n].status == k {
								fmt.Println("killed")
								cancel()
							} else {
							Loopfinished:
								for {
									select {
									case <-c2.Done():
										fmt.Println("stop signal and finished")
										break Loopfinished
									}
								}
							}
							w.jobs[n].Unlock()

							fmt.Println("Stopped with context, exiting in 500 milliseconds")
							time.Sleep(500 * time.Millisecond)

							return
						case <-expiredChan:
							fmt.Println("Stopped and expired")
							cancel()
							w.jobs[n].Lock()
							w.jobs[n].status = ste
							w.jobs[n].Unlock()

							return
						case <-c2.Done():
							w.jobs[n].Lock()
							w.jobs[n].status = wtf
							w.jobs[n].Unlock()
							fmt.Println(w.jobs[n].status)

							break Looptick
						}
					}
					fmt.Println("break looptick")
				}
			}
		}
	}
}

func (w *Worker) deleteKilled(n int) error {
	if w.jobs[n].status == k {
		if err := w.Delete(n); err != nil {
			return fmt.Errorf("Error in deleteKilled(): %v", err)
		}
	}

	return nil
}
