package worker

import (
	"context"
	"fmt"
	"time"

	t "github.com/vslchnk/goscheduler/task"
)

type Worker struct {
	Jobs []*Job
}

type Job struct {
	task      t.Task
	ctx       context.Context
	cancelCtx context.CancelFunc
	status    string
}

// adds task to job pool, if name of task is unique, if ok return number of task in job pool, if not return number of task with the same name and error
func (w *Worker) Add(task t.Task) (int, error) {
	for k, v := range w.Jobs {
		if v.task.Name == task.Name {
			return k, fmt.Errorf("Function with name %v already exist", task.Name)
		}
	}

	//c, cancel := context.WithCancel(context.Background())
	j := &Job{task: task, status: "created, not working"}

	w.Jobs = append(w.Jobs, j)
	return len(w.Jobs) - 1, nil
}

// prints tasks in job pool
func (w Worker) PrintAll() {
	for k, _ := range w.Jobs {
		w.Print(k)
	}
}

func (w Worker) Print(n int) error {
	fmt.Printf("Number: %v; Status: %v; ", n, w.Jobs[n].status)
	w.Jobs[n].task.Print()

	return nil
}

// returns task number in job pool by its' name
func (w Worker) GetNumberByName(name string) (int, error) {
	for k, v := range w.Jobs {
		if v.task.Name == name {
			return k, nil
		}
	}

	return -100, fmt.Errorf("No task with such name")
}

func (w *Worker) Start(n int) error {
	fmt.Println("worker started")

	go w.startJob(n)

	return nil
}

func StartAll() error {
	return nil
}

func (w *Worker) Stop(n int) error {
	fmt.Println("Stop", w.Jobs[n].status)
	w.Jobs[n].status = "stopped"
	w.Jobs[n].cancelCtx()

	return nil
}

func StopAll() error {
	return nil
}

func Kill(n int) error {
	return nil
}

func KillAll() error {
	return nil
}

func (w *Worker) startJob(n int) {
	ctx, cancel := context.WithCancel(context.Background())

	w.Jobs[n].status = "working"
	w.Jobs[n].ctx = ctx
	w.Jobs[n].cancelCtx = cancel

	delayChan := time.NewTimer(time.Second * time.Duration(w.Jobs[n].task.Delay)).C
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Stopped with context, exiting in 500 milliseconds")
			time.Sleep(500 * time.Millisecond)

			return
		case <-delayChan:
			fmt.Println("Delay is done")
			var tickChan *time.Ticker
			//tickChan := time.NewTicker(time.Nanosecond).C

			for {
				select {
				case <-ctx.Done():
					fmt.Println("Stopped with context, exiting in 500 milliseconds")
					time.Sleep(500 * time.Millisecond)

					return
				case <-tickChan:
					fmt.Println("Ticker ticked")
					c2, cancel := context.WithCancel(ctx)
					c1 := context.WithValue(c2, "func", cancel)
					go w.Jobs[n].task.Do(c1)
					tickChan = time.NewTicker(time.Second * time.Duration(w.Jobs[n].task.Period)).C
					expiredChan := time.NewTimer(time.Second * time.Duration(w.Jobs[n].task.TaskTime)).C
				Looptick:
					for {
						select {
						case <-ctx.Done():
							fmt.Println("Stopped with context, exiting in 500 milliseconds")
							time.Sleep(500 * time.Millisecond)

							return
						case <-expiredChan:
							fmt.Println("Stopped and expired")
							w.Jobs[n].status = "stopped, expired"
							//status = "Stopped and expired"

							return
						case <-c2.Done():
							w.Jobs[n].status = "stopped, work is done"
							fmt.Println(w.Jobs[n].status)

							//status = "Stopped and work is done"
							break Looptick
						}
					}
					fmt.Println("break looptick")
				}
			}
		}
	}
}

/*func startJob(job *Job) {
	delayChan := time.NewTimer(time.Second).C
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Stopped with context, exiting in 500 milliseconds")
			time.Sleep(500 * time.Millisecond)

			return
		case <-delayChan:
			fmt.Println("Delay is done")
			tickChan := time.NewTicker(time.Millisecond * 400).C
			for {
				select {
				case <-ctx.Done():
					fmt.Println("Stopped with context, exiting in 500 milliseconds")
					time.Sleep(500 * time.Millisecond)

					return
				case <-tickChan:
					fmt.Println("Ticker ticked")
					expiredChan := time.NewTimer(time.Millisecond * 50000).C
					c2, cancel := context.WithCancel(ctx)
					c1 := context.WithValue(c2, "func", cancel)
					go task.Do(c1)
				Looptick:
					for {
						select {
						case <-ctx.Done():
							fmt.Println("Stopped with context, exiting in 500 milliseconds")
							time.Sleep(500 * time.Millisecond)

							return
						case <-expiredChan:
							fmt.Println("Stopped and expired")

							//status = "Stopped and expired"

							return
						case <-c2.Done():
							fmt.Println("Stopped and work is done")
							//status = "Stopped and work is done"
							break Looptick
						}
					}
					fmt.Println("break looptick")
				}
			}
		}
	}
}*/
