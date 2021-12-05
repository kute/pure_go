package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {

	// 一个互斥锁只能同时被一个 goroutine 锁定，其它 goroutine 将阻塞直到互斥锁被解锁（重新争抢对互斥锁的锁定）
	//testMutex()

	//同时只能有一个 goroutine 能够获得写锁定；
	//同时可以有任意多个 gorouinte 获得读锁定；
	//同时只能存在写锁定或读锁定（读和写互斥）
	//testRWMutex()

	ss := []int{1, 2, 3}
	fmt.Printf("len=%d, cap=%d\n", len(ss), cap(ss))
	//fmt.Println(ss[0:4]) // panic,切片越界

	dd := make([]int, 3, 5)
	fmt.Println(dd)
	//fmt.Println(dd[0:5]) // 不会越界，容量<=5就不会越界

	dd = append(dd, 1)
	fmt.Printf("dd=%v, len=%d, cap=%d\n", dd, len(dd), cap(dd))

	dd = append(dd, 2)
	fmt.Printf("dd=%v, len=%d, cap=%d\n", dd, len(dd), cap(dd))

	// 当添加3之后，长度 大于 5 了，所以会进行扩容，直接翻倍
	dd = append(dd, 3)
	fmt.Printf("dd=%v, len=%d, cap=%d\n", dd, len(dd), cap(dd))

	hh := dd[3:4]
	fmt.Printf("hh=%v, len=%d, cap=%d\n", hh, len(hh), cap(hh))

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
