package main

import (
	"context"
	"fmt"
	t "scheduler/task"
	"time"
)

func outer(name string) func() error {
	text := "Modified " + name

	foo := func() error {
		fmt.Println(text)

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
	foo := outer("hello")
	foo2 := outer2("hello", "piss")

	a, _ := t.Create(10.0, 10.0, 10.0, foo)
	a.Print()
	a.Delay = 20.0
	a.Print()

	a.Do()

	foo()
	foo2()

	c1, cancel := context.WithCancel(context.Background())

	exitCh := make(chan struct{})

	go startJob(c1, a, exitCh)

	go func() {
		time.Sleep(5000 * time.Millisecond)
		cancel()
		return
	}()

	<-exitCh
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
					expiredChan := time.NewTimer(time.Nanosecond).C
					task.Do()
					select {
					case <-expiredChan:
						fmt.Println("Stopped and expired")
						//status = "Stopped and expired"
						exitCh <- struct{}{}
						return
					default:
						fmt.Println("Stopped and work is done")
						//status = "Stopped and work is done"
						break
					}
				}
			}
		}
	}
}
