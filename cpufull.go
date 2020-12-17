package main

import (
	"sync"
)

func main() {

	quit := make(chan bool)
	wg := sync.WaitGroup{}

	for i := 0; i < 100000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-quit:
					break
				default:

				}
			}
		}()
	}

	wg.Wait()
}
