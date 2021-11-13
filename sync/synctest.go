package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {

	// 一个互斥锁只能同时被一个 goroutine 锁定，其它 goroutine 将阻塞直到互斥锁被解锁（重新争抢对互斥锁的锁定）
	testMutex()

	//同时只能有一个 goroutine 能够获得写锁定；
	//同时可以有任意多个 gorouinte 获得读锁定；
	//同时只能存在写锁定或读锁定（读和写互斥）
	testRWMutex()

}

func testRWMutex() {
}

func testMutex() {
	var a int
	var lock sync.Mutex
	for i := 0; i < 1000; i++ {
		go func() {
			lock.Lock()
			defer lock.Unlock()
			a++
			fmt.Println(a)
		}()
	}
	time.Sleep(time.Second)
}
