package main

import (
	"fmt"
)

var printChan = make(chan string)

func makeSource(nums ...int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for _, n := range nums {
			printChan <- fmt.Sprintf("Source: %d", n)
			out <- n
		}
	}()

	return out
}

func stage0(in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for n := range in {
			printChan <- fmt.Sprintf("stage0: %d", n)
			out <- n * n
		}
	}()

	return out
}

func main() {
	quitChan := make(chan int)
	go func() {
		defer close(quitChan)

		for s := range printChan {
			fmt.Println(s)
		}
		quitChan <- 1
	}()

	source := makeSource(3, 4, 5, 6, 7, 8)
	out := stage0(source)

	for x := range out {
		printChan <- fmt.Sprintf("0: %d", x)
	}
	close(printChan)

	// time.Sleep(5 * time.Second)
	// fmt.Println("Woke up")

	<-quitChan
}
