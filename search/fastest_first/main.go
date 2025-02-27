package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Result string
type Search func(query string) Result

func First(query string, replicas ...Search) Result {
	c := make(chan Result)
	for _, replica := range replicas {
		go func() {
			c <- replica(query)
		}()
	}
	return <-c
}

func fakeSearch(kind string) Search {
	return func(query string) Result {
		time.Sleep(time.Duration(rand.Intn(100) * int(time.Millisecond)))
		return Result(fmt.Sprintf("%s result for %q\n", kind, query))
	}

}

func firstMain() {
	start := time.Now()
	searchers := []Search{
		fakeSearch("r1"),
		fakeSearch("r2"),
		fakeSearch("r3"),
	}

	result := First("maw", searchers...)
	elapsed := time.Since(start)
	fmt.Println(result)
	fmt.Printf("%dms\n", elapsed.Milliseconds())
}

func main() {
	r := make(chan Result)
	query := "golang"

	go func() {
		r <- First(
			query,
			fakeSearch("web-r1"),
			fakeSearch("web-r2"),
			fakeSearch("web-r3"),
		)
	}()

	go func() {
		r <- First(
			query,
			fakeSearch("image-r1"),
			fakeSearch("image-r2"),
			fakeSearch("image-r3"),
		)
	}()

	go func() {
		r <- First(
			query,
			fakeSearch("video-r1"),
			fakeSearch("video-r2"),
			fakeSearch("video-r3"),
		)
	}()

	timout := time.After(80 * time.Millisecond)
	start := time.Now()
	timedout := false
	for range 3 {
		select {
		case result := <-r:
			fmt.Print("Result: ", result)
		case <-timout:
			timedout = true
			fmt.Println("⏰ timeout!")
		}
	}

	if !timedout {
		elapsed := time.Since(start)
		fmt.Printf("Elapsed: %v\n", elapsed)
	}
}
