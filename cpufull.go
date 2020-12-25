package main

import (
	"container/list"
	"etcdexample/commons"
	"sync"
)

func main() {

	quit := make(chan bool)
	wg := sync.WaitGroup{}

	str := commons.RandString(1024)
	for i := 0; i < 1000; i++ {
		//for {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				select {
				case <-quit:
					break
				default:
					compute(str)

				}
			}
		}()
	}

	wg.Wait()
}

func compute(str string) {
	l := list.New()

	for i := 0; i < 100000; i++ {
		l.PushBack(str)
	}

	for i := 0; i < 100000; i++ {
		for e := l.Front(); e != nil; e = e.Next() {
		}

	}
}

func fibonacci(num int) int {
	if num < 2 {
		return 1
	}

	return fibonacci(num-1) + fibonacci(num-2)
}
