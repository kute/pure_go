package main

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

/**
进程/线程
进程是程序在操作系统中的一次执行过程，系统进行资源分配和调度的一个独立单位。

线程是进程的一个执行实体，是 CPU 调度和分派的基本单位，它是比进程更小的能独立运行的基本单位。

一个进程可以创建和撤销多个线程，同一个进程中的多个线程之间可以并发执行。
并发/并行
多线程程序在单核心的 cpu 上运行，称为并发；多线程程序在多核心的 cpu 上运行，称为并行。

并发与并行并不相同，并发主要由切换时间片来实现“同时”运行，并行则是直接利用多核实现多线程的运行，Go程序可以设置使用核心数，以发挥多核计算机的能力。
协程/线程
协程：独立的栈空间，共享堆空间，调度由用户自己控制，本质上有点类似于用户级线程，这些用户级线程的调度也是自己实现的。

线程：一个线程上可以跑多个协程，协程是轻量级的线程。
*/

func main() {

	fmt.Println("NumCPU", runtime.NumCPU(), runtime.NumGoroutine())
	runtime.GOMAXPROCS(runtime.NumCPU())

	//funcDefine()

	//go timePrint()
	//var content string
	//fmt.Scanln(&content)
	//fmt.Println("receive content", content)

	//test2()

	//test3()

	//test3_1()

	//test3_2()

	//testChannel()

	//testChannel_single()

}

// 单向通道：通道只能读 或者 写：将一个 channel 变量传递到一个函数时，可以通过将其指定为单向 channel 变量，从而限制该函数中可以对此 channel 的操作
func testChannel_single() {

	// 创建通道，每个通道只允许指定类型的数据放入
	var intChan = make(chan int32)
	//var intReadChan = make(<-chan int32)  // 创建只读通道
	//var intWriteChan = make(chan<- int32) // 创建只写通道
	defer close(intChan) // 关闭通道
	var count int32 = 0
	var wg sync.WaitGroup

	wg.Add(2)

	// 声明该通道只能 读
	go func(intChan <-chan int32) {
		defer wg.Done()
		// 循环阻塞接收通道数据
		for data := range intChan {
			fmt.Println("consume value", data)
			if data >= 10 {
				break
			}
		}
		fmt.Println("consume done")

		// intChan <- 88 // Panic，无法发送，通道被声明为只能读

	}(intChan)

	// 声明该通道只能 写
	go func(intChan chan<- int32) {
		defer wg.Done()
		var ticket = time.Tick(time.Second)
		for range ticket {
			// 原子性递增
			atomic.AddInt32(&count, 1)
			fmt.Println("produce value:", count)
			// 发送到通过内，阻塞，必须等待上一个发送的数据被接收才能发送下一个
			intChan <- count
			if count >= 10 {
				break
			}
		}
		fmt.Println("produce done")

		//data := <-intChan  // Panic，无法接收，通道被声明为只能写
	}(intChan)
	wg.Wait()
	fmt.Println("main end ")

}

/**
1. 在任何时候，同时只能有一个 goroutine 访问通道进行发送和获取数据
2. 遵循先入先出（First In First Out）的规则，类似一个队列，保证收发数据的顺序
① 通道的收发操作在不同的两个 goroutine 间进行。
由于通道的数据在没有接收方处理时，数据发送方会持续阻塞，因此通道的接收必定在另外一个 goroutine 中进行。
② 接收将持续阻塞直到发送方发送数据。
如果接收方接收时，通道中没有发送方发送数据，接收方也会发生阻塞，直到发送方发送数据为止。
data := <- ch   // 此为阻塞方式
data, ok := <- ch  // ok表示通道是否处于打开，true 打开状态，false 关闭
<-ch  // 阻塞接收数据后，忽略从通道返回的数据
for data := range ch {}   // 循环接收通道数据
③ 每次接收一个元素。
通道一次只能接收一个数据元素
*/
func testChannel() {

	// 创建通道，每个通道只允许指定类型的数据放入
	var intChan = make(chan int32)
	defer close(intChan)
	//var bodyChan = make(chan *Body)
	var count int32 = 0
	var wg sync.WaitGroup

	wg.Add(2)

	go func(intChan chan int32) {
		defer wg.Done()
		// 循环阻塞接收通道数据
		for data := range intChan {
			fmt.Println("consume value", data)
			if data >= 10 {
				break
			}
		}
		fmt.Println("consume done")
		intChan <- 88
		fmt.Println("consume produce new value done")
	}(intChan)

	go func(intChan chan int32) {
		defer wg.Done()
		var ticket = time.Tick(time.Second)
		for range ticket {
			atomic.AddInt32(&count, 1)
			fmt.Println("produce value:", count)
			// 发送到通过内，阻塞，必须等待上一个发送的数据被接收才能发送下一个
			intChan <- count
			if count >= 10 {
				break
			}
		}
		fmt.Println("produce done")
		data := <-intChan
		fmt.Println("produce consume new value done", data)
	}(intChan)
	wg.Wait()
	fmt.Println("main end ")

}

// 原子读写
// main 函数使用 StoreInt64 函数来安全地修改 count 变量的值。如果哪个 goroutine 试图在 main 函数调用 StoreInt64 的同时调用 LoadInt64 函数，那么原子函数会将这些调用互相同步，保证这些操作都是安全的，不会进入竞争状态
func test3_2() {
	var count int32 = 0
	var size int = 2
	var wg sync.WaitGroup

	wg.Add(size)
	for i := 0; i < size; i++ {
		go func(i int) {
			defer wg.Done()
			for {
				fmt.Println("goroutine -", i, "sleep")
				time.Sleep(250 * time.Millisecond)
				// 当前协程发现main有调用StoreInt32的情况下，又再调用 LoadInt32 会自动同步
				if atomic.LoadInt32(&count) == 8 {
					fmt.Println("goroutine -", i, "shut down")
					break
				}
			}
		}(i)
	}
	fmt.Println("total NumGoroutine", runtime.NumGoroutine())
	time.Sleep(5 * time.Second)
	// main函数调用 StoreInt32 修改count值
	atomic.StoreInt32(&count, 8)
	wg.Wait() //等待goroutine结束
}

// 原子读写
func test3_1() {
	var count int32 = 0
	var size int = 100
	var wg sync.WaitGroup

	wg.Add(size)
	for i := 0; i < size; i++ {
		go func() {
			defer wg.Done()
			atomic.AddInt32(&count, 1) // 原子增加
		}()
	}
	wg.Wait() //等待goroutine结束
	fmt.Println("after for count", count)
}

func test3() {
	var count, size = 0, 100
	var wg sync.WaitGroup

	wg.Add(size)
	for i := 0; i < size; i++ {
		go func() {
			defer wg.Done()
			// 读写未枷锁，结果有问题,通过 go build -race concurrenttest.go ，然后执行可执行文件看到控制台日志关于共享资源的竞争问题
			count++
		}()
	}
	wg.Wait() //等待goroutine结束
	fmt.Println("after for count", count)
}

func test2() {
	var count = 0
	var lock sync.Mutex
	for i := 0; i < 10; i++ {
		go func() {
			lock.Lock()
			defer lock.Unlock()
			{
				count++
			}
			time.Sleep(time.Second)
		}()
	}
	// check all goroutine if done
	for {
		lock.Lock()
		c := count
		lock.Unlock()
		// 让出cpu时间片
		runtime.Gosched()
		if c >= 10 {
			fmt.Println("all goroutine done for count", count)
			break
		}
	}
}

func timePrint() {
	var ticket = time.Tick(time.Second)
	for i := range ticket {
		fmt.Println(i.Format("2006-01-03 15:04:05"))
	}
}

func sleep() {
	fmt.Println("begin execute depth and sleep")
	time.Sleep(5 * time.Second)
}

func funcDefine() {
	fmt.Println("main start")
	// 函数调用前添加go关键字
	go sleep()

	// 匿名方法执行
	go func() {
		time.Sleep(5 * time.Second)
	}()

	fmt.Println("main end")
}
