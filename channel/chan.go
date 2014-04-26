package channel

import (
	"fmt"
	"sync"
)

func Run() {
	done := make(chan struct{}, 2)
	defer close(done)

	c := gen(done, 2, 3, 4, 5, 6, 7, 8)
	out1 := sq(done, c)
	out2 := sq(done, c)
	out3 := sq(done, c)

	out := merge(done, out1, out2, out3)
	fmt.Println(<-out)
	fmt.Println(<-out)
	fmt.Println(<-out)
}

func gen(done <-chan struct{}, nums ...int) <-chan int {
	out := make(chan int, len(nums))
	go func() {
		defer close(out)
		for _, n := range nums {
			select {
			case out <- n:
			case <-done:
				return
			}
		}
	}()
	return out
}

func sq(done chan struct{}, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			select {
			case out <- n * n:
			case <-done:
				return
			}
		}
	}()
	return out
}

func merge(done <-chan struct{}, cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	output := func(in <-chan int) {
		defer wg.Done()
		for v := range in {
			select {
			case out <- v:
			case <-done:
				return
			}
		}
	}
	wg.Add(len(cs))
	for _, v := range cs {
		go output(v)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
