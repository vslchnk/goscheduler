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
		//time.Sleep(time.Second * 2)
		//fmt.Println("start")
		//time.Sleep(time.Second * 6)
		c2, cancel := context.WithCancel(ctx)
		defer cancel()
		for i := 0; i < 3; i++ {
			time.Sleep(time.Second * 2)
			select {
			case <-c2.Done():
				fmt.Println("Done")
				return nil
			default:
				fmt.Println(text)
			}
		}
		fmt.Println(ctx.Err() == context.Canceled)
		//fmt.Println(text /*+ ctx.Value("func").(string)*/)
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

	a, _ := t.Create("printing", time.Second*3, time.Second*3, time.Second*1, foo)
	b, _ := t.Create("printing2", time.Second*3, time.Second*3, time.Second*1, foo)
	//a.Print()
	//a.Delay = 20.0
	//a.Print()
	worker.Add(a)
	worker.Add(b)
	//worker.Print(0)
	worker.PrintAll()
	//worker.Delete(1)
	worker.PrintAll()
	fmt.Println(worker.GetNumberByName("printing"))
	fmt.Println(worker.Start(0))
	time.Sleep(4 * time.Second)
	//worker.Stop(0)
	worker.Kill(0)
	time.Sleep(9 * time.Second)
	/*worker.Start(0)
	time.Sleep(4 * time.Second)
	//worker.Stop(0)
	worker.Kill(0)
	time.Sleep(9 * time.Second)*/
	worker.PrintAll()
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
