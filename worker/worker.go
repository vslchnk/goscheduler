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

type worker struct {
	sync.Mutex
	jobs map[string]*job
}

type job struct {
	sync.Mutex
	task      t.Task
	ctx       context.Context
	cancelCtx context.CancelFunc
	status    jobstatus
}

// creates new worker
func NewWorker() *worker {
	w := worker{}
	w.jobs = make(map[string]*job)
	return &w
}

// checks if there is a job with name n in job pool
func (w *worker) check(n string) bool {
	defer w.Unlock()
	w.Lock()

	_, ok := w.jobs[n]

	return ok
}

// change task in job pool by its name
func (w *worker) ChangeTask(n string, task t.Task) error {
	if !w.check(n) {
		return fmt.Errorf("No job with name %v in job pool", n)
	}

	defer w.jobs[n].Unlock()
	w.jobs[n].Lock()

	if w.jobs[n].status == wtnf || w.jobs[n].status == wtf {
		return fmt.Errorf("Job is working, must be stopped before being changed")
	}

	w.jobs[n].task = task

	return nil
}

// adds task to job pool, if name n of job is unique, if ok return number of job in job pool, if not return number of job with the same name and error
func (w *worker) Add(task t.Task, n string) error {

	if w.check(n) {
		return fmt.Errorf("Function with name %v already exist", n)
	}

	j := &job{task: task, status: cnw}

	w.Lock()
	w.jobs[n] = j
	w.Unlock()

	return nil
}

// prints jobs in job pool
func (w worker) PrintAll() error {
	for k, _ := range w.jobs {
		if err := w.Print(k); err != nil {
			return fmt.Errorf("Error in PrintAll(): %v", err)
		}
	}

	return nil
}

// prints job from job pool by its name
func (w worker) Print(n string) error {
	if !w.check(n) {
		return fmt.Errorf("No job with name %v in job pool", n)
	}

	fmt.Printf("Name: %v; Status: %v; ", n, w.jobs[n].status)
	w.jobs[n].task.Print()

	return nil
}

// starts job by its name
func (w *worker) Start(n string) error {
	if !w.check(n) {
		return fmt.Errorf("No job with name %v in job pool", n)
	}

	w.jobs[n].Lock()

	if w.jobs[n].status == wtnf || w.jobs[n].status == wtf {
		w.jobs[n].Unlock()
		return fmt.Errorf("Job name %v is already working, its status: %v", n, w.jobs[n].status)
	}

	ctx, cancel := context.WithCancel(context.Background())

	w.jobs[n].status = wtnf
	w.jobs[n].ctx = ctx
	w.jobs[n].cancelCtx = cancel
	w.jobs[n].Unlock()

	go w.startJob(n)

	return nil
}

// starts all jobs if they are stopped or not started
func (w *worker) StartAll() error {
	for k, _ := range w.jobs {
		if err := w.Start(k); err != nil {
			return fmt.Errorf("Error in StartAll(): %v", err)
		}
	}

	return nil
}

// stops job by its name
func (w *worker) Stop(n string) error {
	if !w.check(n) {
		return fmt.Errorf("No job with name %v in job pool", n)
	}

	w.jobs[n].Lock()
	if w.jobs[n].status == sss || w.jobs[n].status == ste || w.jobs[n].status == k || w.jobs[n].status == cnw {
		w.jobs[n].Unlock()
		return fmt.Errorf("Job name %v is not working, its status: %v", n, w.jobs[n].status)
	}
	w.jobs[n].status = sss
	w.jobs[n].cancelCtx()
	w.jobs[n].Unlock()

	return nil
}

// stops all jobs if they are not stopped or not started
func (w *worker) StopAll() error {
	for k, _ := range w.jobs {
		if err := w.Stop(k); err != nil {
			return fmt.Errorf("Error in StopAll(): %v", err)
		}
	}

	return nil
}

// kills job and removes it from pool
func (w *worker) Kill(n string) error {
	if !w.check(n) {
		return fmt.Errorf("No job with name %v in job pool", n)
	}

	w.jobs[n].Lock()

	if w.jobs[n].status == sss || w.jobs[n].status == ste || w.jobs[n].status == k || w.jobs[n].status == cnw {
		w.jobs[n].Unlock()
		return fmt.Errorf("Job name %v is not working, its status: %v", n, w.jobs[n].status)
	}

	w.jobs[n].status = k
	w.jobs[n].cancelCtx()
	w.jobs[n].Unlock()

	return nil
}

// kills all jobs and removes them from pool
func (w *worker) KillAll() error {
	for k, _ := range w.jobs {
		if err := w.Kill(k); err != nil {
			return fmt.Errorf("Error in KillAll(): %v", err)
		}
	}

	return nil
}

// delets job from pool by its number
func (w *worker) Delete(n string) error {
	if !w.check(n) {
		return fmt.Errorf("No job with name %v in job pool", n)
	}

	w.Lock()
	if w.jobs[n].status == wtf || w.jobs[n].status == wtnf {
		w.Unlock()
		return fmt.Errorf("Job name %v is working, its status: %v", n, w.jobs[n].status)
	}

	delete(w.jobs, n)
	w.Unlock()

	return nil
}

// controls work of job
func (w *worker) startJob(n string) {
	defer w.deleteKilled(n)
	defer w.jobs[n].cancelCtx()

	delayChan := time.NewTimer(time.Duration(w.jobs[n].task.GetDelay())).C
	for {
		select {
		case <-w.jobs[n].ctx.Done():

			return
		case <-delayChan:

			tickChan := time.NewTimer(0).C
			for {
				select {
				case <-w.jobs[n].ctx.Done():

					return
				case <-tickChan:

					c2, cancel := context.WithCancel(context.Background())
					c1 := context.WithValue(c2, "func", cancel)

					go w.jobs[n].task.GetDoFunc()(c1)

					tickChan = time.NewTimer(time.Duration(w.jobs[n].task.GetPeriod())).C
					expiredChan := time.NewTimer(time.Duration(w.jobs[n].task.GetTaskTime())).C
				Looptick:
					for {
						select {
						case <-w.jobs[n].ctx.Done():

							w.jobs[n].Lock()
							if w.jobs[n].status == k {
								cancel()
							} else {
							Loopfinished:
								for {
									select {
									case <-c2.Done():
										break Loopfinished
									}
								}
							}
							w.jobs[n].Unlock()

							return
						case <-expiredChan:
							cancel()
							w.jobs[n].Lock()
							w.jobs[n].status = ste
							w.jobs[n].Unlock()

							return
						case <-c2.Done():
							w.jobs[n].Lock()
							w.jobs[n].status = wtf
							w.jobs[n].Unlock()

							break Looptick
						}
					}
					fmt.Println("break looptick")
				}
			}
		}
	}
}

// deletes job if it's killed from pool
func (w *worker) deleteKilled(n string) error {
	if !w.check(n) {
		return fmt.Errorf("No job with name %v in job pool", n)
	}

	if w.jobs[n].status == k {
		if err := w.Delete(n); err != nil {
			return fmt.Errorf("Error in deleteKilled(): %v", err)
		}
	}

	return nil
}
