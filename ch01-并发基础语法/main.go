package main

import (
	"fmt"
	"time"
)

func test() {
	fmt.Println("hello world")
}

func testChannel(ch <-chan int) {
	for v := range ch {
		fmt.Println(v)
	}
}

func testChannel1(ch chan<- int) {
	var i int
	for {
		ch <- i
		i++
	}
}

func main() {
	// 启动一个goroutine
	go test()
	time.Sleep(time.Second)

	// channel示例
	ch := make(chan int, 2)
	go testChannel1(ch)
	go testChannel(ch)
	go testChannel(ch)
	time.Sleep(time.Second)

}
