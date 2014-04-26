package channel

import (
	"fmt"
	"sync"
)

func Run() {
	c := gen(2, 3, 4, 5, 6)
	out1 := sq(c)
	out2 := sq(c)
	out3 := sq(c)

	done := make(chan struct{}, 2)

	out := merge(done, out1, out2, out3)
	fmt.Println(<-out)

	done <- struct{}{}
	done <- struct{}{}
	done <- struct{}{}

}

func gen(nums ...int) <-chan int {
	out := make(chan int, len(nums))
	for _, n := range nums {
		out <- n
	}
	close(out)
	return out
}

func sq(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

func merge(done <-chan struct{}, cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	output := func(in <-chan int) {
		for v := range in {
			out <- v
			select {
			case out <- v:
			case <-done:
			}
		}
		wg.Done()
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
