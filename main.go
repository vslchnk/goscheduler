package main

import (
	"context"
	"fmt"
	"time"

	t "github.com/vslchnk/goscheduler/task"
	w "github.com/vslchnk/goscheduler/worker"
)

func outer(name string) func(ctx context.Context) error {
	text := "Modified " + name

	foo := func(ctx context.Context) error {
		defer ctx.Value("func").(context.CancelFunc)()
		//select {
		//case <-ctx.Done():
		time.Sleep(time.Second * 2)
		fmt.Println("start")
		//c2, cancel := context.WithCancel(ctx)
		//defer cancel()
		/*for {
			select {
			case <-c2.Done():
				fmt.Println("Done")
			default:
				fmt.Println(text)
			}
		}*/
		fmt.Println(ctx.Err() == context.Canceled)
		fmt.Println(text /*+ ctx.Value("func").(string)*/)
		//}

		return nil
	}

	return foo
}

func outer2(name string, name2 string) func() error {

	text := "Modified " + name + name2

	foo := func() error {
		fmt.Println(text)

		return nil
	}

	return foo
}

func main() {
	worker := w.Worker{}
	foo := outer("hello")
	//foo2 := outer2("hello", "piss")

	a, _ := t.Create("printing", time.Second*1, time.Second*1, time.Nanosecond, foo)
	a.Print()
	//a.Delay = 20.0
	a.Print()
	worker.Add(a)
	worker.Print(0)
	worker.PrintAll()
	fmt.Println(worker.GetNumberByName("printing"))
	worker.Start(0)
	//time.Sleep(9 * time.Second)
	worker.Stop(0)
	time.Sleep(5 * time.Second)
	/*c2, cancel := context.WithCancel(context.Background())
	c1 := context.WithValue(c2, "func", cancel)
	a.Do(c1)
	fmt.Println(c2.Err() == context.Canceled)

	foo(c1)
	foo2()*/

	/*c1, cancel := context.WithCancel(context.Background())

	exitCh := make(chan struct{})

	go startJob(c1, a, exitCh)

	go func() {
		time.Sleep(5000 * time.Millisecond)
		cancel()
		return
	}()

	go startJob(c1, a, exitCh)

	go func() {
		time.Sleep(5000 * time.Millisecond)
		cancel()
		return
	}()

	<-exitCh*/
}

func startJob(ctx context.Context, task t.Task, exitCh chan struct{}) {
	delayChan := time.NewTimer(time.Second).C
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Stopped with context, exiting in 500 milliseconds")
			time.Sleep(500 * time.Millisecond)
			exitCh <- struct{}{}
			return
		case <-delayChan:
			fmt.Println("Delay is done")
			tickChan := time.NewTicker(time.Millisecond * 400).C
			for {
				select {
				case <-ctx.Done():
					fmt.Println("Stopped with context, exiting in 500 milliseconds")
					time.Sleep(500 * time.Millisecond)
					exitCh <- struct{}{}
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
						/*case <-ctx.Done():
						fmt.Println("Stopped with context, exiting in 500 milliseconds")
						time.Sleep(500 * time.Millisecond)
						exitCh <- struct{}{}
						return*/
						case <-expiredChan:
							fmt.Println("Stopped and expired")

							//status = "Stopped and expired"
							exitCh <- struct{}{}
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
}
