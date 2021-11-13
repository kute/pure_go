package main

import (
	"bufio"
	"bytes"
	"container/list"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"reflect"
	"strconv"
	"sync"
	"time"
	"unsafe"
)

// 单个包的初始化过程是，先初始化常量，然后是全局变量，最后执行包的 init 函数
// 初始化函数
func init() {
	fmt.Println("ini before main")
}

func IoTest() {

	content := []byte("接口技术的")
	// 创建一个reader
	br := bytes.NewReader(content) // 不可重复使用
	ir := bufio.NewReader(br)

	fmt.Println("ir.Read")
	var data [15]byte // 声明一个数组，固定长度
	// 将 ir中的字节读到 data中
	if size, err := ir.Read(data[:]); err == nil {
		fmt.Println(string(data[:]), size, err)
	}

	fmt.Println("ir.ReadByte")
	ir2 := bufio.NewReader(bytes.NewReader(content))
	var data1 []byte // 声明一个切片，非固定
	for {
		if c, err := ir2.ReadByte(); err == nil {
			data1 = append(data1, c)
		} else {
			break
		}
	}
	fmt.Println(string(data1), len(data1))

	fmt.Println("ir.ReadBytes")
	ir3 := bufio.NewReader(bytes.NewReader(content))
	var delim byte = '\n'
	// 读取数据直到遇到第一个分隔符“delim”，并返回读取的字节序列（包括“delim”）。
	// 如果 ReadBytes 在读到第一个“delim”之前出错，它返回已读取的数据和那个错误（通常是 io.EOF）。
	// 只有当返回的数据不以“delim”结尾时，返回的 err 才不为空值
	if databytes, err := ir3.ReadBytes(delim); err == nil || err == io.EOF {
		fmt.Println(string(databytes), len(databytes), err)
	} else {
		fmt.Println(err)
	}

	fmt.Println("ir.ReadString")
	ir4 := bufio.NewReader(bytes.NewReader(content))
	// 读取数据直到分隔符“delim”第一次出现，并返回一个包含“delim”的字符串。
	//如果 ReadString 在读取到“delim”前遇到错误，它返回已读字符串和那个错误（通常是 io.EOF）。
	//只有当返回的字符串不以“delim”结尾时，ReadString 才返回非空 err
	if databytes, err := ir4.ReadString(delim); err == nil || err == io.EOF {
		fmt.Println(string(databytes), len(databytes), err)
	} else {
		fmt.Println(err)
	}

	fmt.Println("ir.Peek")
	ir5 := bufio.NewReader(bytes.NewReader(content))
	// 读取数据直到分隔符“delim”第一次出现，并返回一个包含“delim”的字符串。
	//如果 ReadString 在读取到“delim”前遇到错误，它返回已读字符串和那个错误（通常是 io.EOF）。
	//只有当返回的字符串不以“delim”结尾时，ReadString 才返回非空 err
	if databytes, err := ir5.Peek(200); err == nil || err == io.EOF {
		fmt.Println(string(databytes), len(databytes), err)
	} else {
		fmt.Println(err)
	}

	timeTest()

}

func timeTest() {
	time.Sleep(time.Second)
	fmt.Println(time.April.String())
	now := time.Now() //获取当前时间
	fmt.Printf("current time:%v\n", now)
	year := now.Year()            //年
	month := now.Month()          //月
	day := now.Day()              //日
	hour := now.Hour()            //小时
	minute := now.Minute()        //分钟
	second := now.Second()        //秒
	secondtimestamp := now.Unix() // 时间戳,秒
	timestamp := now.UnixNano()   // 纳秒时间戳
	fmt.Printf("%d-%02d-%02d %02d:%02d:%02d\n", year, month, day, hour, minute, second)
	fmt.Println(secondtimestamp, timestamp)
	fmt.Println(time.Millisecond)
	fmt.Println(time.Unix(secondtimestamp, 0))
	fmt.Println(time.UnixDate)

	ticker := time.Tick(time.Second) //定义一个1秒间隔的定时器
	for i := range ticker {
		fmt.Println(i.Format("2006-01-03 15:04:05")) //每秒都会执行的任务
	}

}

var x, y int

// 批量声明变量
var (
	xx bool
	yy float32
)

// rune 等价于 int32，用来表示unicode
// byte 和 uint8 等价
var a rune

// 类型定义
type newIntType int

// 类型别名
type bytesAlia = float64

// 闭包
func getSequence() func() int {
	i := 0
	return func() int {
		i += 1
		//i++
		return i
	}
}

// params ...interface{} : 表示多个不同类型的参数
// ...type：表示类型为 type的可变参数，即不定数量的参数
func Params(a, b string, id int, params ...interface{}) {
	fmt.Println(id)
	if len(params) > 0 {
		fmt.Println(params[0])
	}
	fmt.Println(reflect.TypeOf(a))
	fmt.Println(reflect.TypeOf(b))
}

func Params2(params ...interface{}) {
	Params("a", "b", 1, params...)
}

func testMath() {
	maxFloat32 := math.MaxFloat32
	fmt.Println(maxFloat32)
}

func testFmt() (string, string, newIntType) {
	Params("b", "a", 3)
	fmt.Print("a", "b")
	fmt.Println("a" + "b")
	var a string = "asdf"
	var b, c string = "asdf", "fsdf"

	d := 3
	f, g := 1, 2

	nf := float32(d) * float32(g)
	fmt.Println(nf)

	const constant = "sf"
	fmt.Println(a, b, c, d, f, g, constant)

	strAry := []string{"a", "b"}
	fmt.Println(strAry)

	var complexVar complex128 = complex(1., .23)
	fmt.Println(complexVar)
	fmt.Println(real(complexVar))
	fmt.Println(imag(complexVar))
	if d > 4 {
		if g < 1 && f > 1 {
			fmt.Println("ffffffffff")
		} else {
			fmt.Println("xxxxx")
		}
	} else {
		fmt.Println("gggggggg")
	}

	// iota特殊常量，iota 在 const关键字出现时将被重置为 0(const 内部的第一行之前)，const 中每新增一行常量声明将使 iota 计数一次(iota 可理解为 const 语句块中的行索引)。
	const (
		a1 = iota // 0
		a2 = iota // 1
		a3 = iota // 2
	)
	fmt.Println(a1, a2, a3)
	const (
		b1 = iota // 0
		b2        // 1
		b3        // 2
		_
		_
		b6 // 5
	)
	fmt.Println(b1, b2, b3, b6)

	const (
		IgEggs         newIntType = 1 << iota // 1 << 0 which is 00000001
		IgChocolate                           // 1 << 1 which is 00000010
		IgNuts                                // 1 << 2 which is 00000100
		IgStrawberries                        // 1 << 3 which is 00001000
		IgShellfish                           // 1 << 4 which is 00010000
	)

	type ByteSize float64
	const (
		_           = iota             // ignore first value by assigning to blank identifier
		KB ByteSize = 1 << (10 * iota) // 1 << (10*1)
		MB                             // 1 << (10*2)
		GB                             // 1 << (10*3)
		TB                             // 1 << (10*4)
		PB                             // 1 << (10*5)
		EB                             // 1 << (10*6)
		ZB                             // 1 << (10*7)
		YB                             // 1 << (10*8)
	)

	const (
		Apple, Banana = iota + 1, iota + 2
		Cherimoya, Durian
		Elderberry, Fig
	)
	// 1 2 2 3 3 4
	fmt.Println(Apple, Banana, Cherimoya, Durian, Elderberry, Fig)

	var x interface{}

	switch p := x.(type) {
	case nil:
		fmt.Printf(" x 的类型 :%v", p)
	case int:
		fmt.Printf("x 是 int 型")
	case float64:
		fmt.Printf("x 是 float64 型")
	case func(int) float64:
		fmt.Printf("x 是 func(int) 型")
	case bool, string:
		fmt.Printf("x 是 bool 或 string 型")
	default:
		fmt.Printf("未知型")
	}

	// 等价于 for(;;)
	for {
		break
	}

	// 等价于 while(condition){}
	for d < 9 {
		fmt.Println("fsdf = " + strconv.Itoa(d))
		d++
	}

	sum := 0
	for i := 1; i < 10; i += 2 {
		sum += i
	}
	fmt.Println(sum)

	for sum <= 10 {
		sum += sum
	}
	fmt.Println(sum)

	for i, v := range strAry {
		fmt.Println(i, v)
	}

	// 初始化大小，默认为0
	numbers := [6]int{1, 2, 3, 5}
	for i, x := range numbers {
		fmt.Printf("第 %d 位 x 的值 = %d\n", i, x)
	}

	var y1 [10]int
	y1[0] = 1
	var y2 = []int{1, 2} // 初始化
	y3 := [4]int{1, 2}   // 初始化大小，其他默认为 类型的零值
	y4 := [...]int{2, 4} // ... 自行推断大小
	var y5 = [...]int{2, 4}
	fmt.Println(len(y2))
	y6 := [...][2]int{{1, 2}, {3, 4}}
	y6[0][1], y6[1][0] = y6[1][0], y6[0][1]
	fmt.Println(y6)
	fmt.Println(y1, y2, y3, y4, y5, y6)
	for _, t := range y4 {
		fmt.Println(t)
	}

	//var y6 [2][3]int // 多维数组
	var y8 [][]int // 多维数组
	//var y10 [...][...]int

	y8 = append(y8, y2)

	//y10 = append(y10, y2)

	//y9 := [][]int{{1, 2}}
	//fmt.Println(y6, y8, y9, y10)

	// 切片，不会发生内存分配
	fmt.Println(y1[2:3])

	// 动态创建切片，一定发生了内存分配
	var l1 = make([]int, 2, 10)
	l1 = append(l1, 2, 2)              // 添加元素
	l1 = append([]int{4}, l1...)       // 头部添加元素
	l1 = append([]int{4, 5, 6}, l1...) // 头部添加切片
	var index = 4
	var ele = 7
	l1 = append(l1[:index], append([]int{ele}, l1[index:]...)...) // 第 5 个位置 插入 7
	fmt.Println(l1)

	// 删除 开头元素
	l1 = l1[1:] // 移动指针删除
	fmt.Println(l1)

	l1 = append([]int{}, l1[1:]...) // append 原地完成（所谓原地完成是指在原有的切片数据对应的内存区间内完成，不会导致内存空间结构的变化）
	fmt.Println(l1)

	l1 = l1[:copy(l1, l1[1:])]
	fmt.Println(l1)

	l1 = l1[:len(l1)-1] // 删除尾部元素

	l1 = append(l1[:2], l1[2+1:]...) // 删除第三个元素
	fmt.Println(l1)

	MapTest()

	return a, b, 3
}

func MapTest() {

	var m1 map[string]int // 非线程安全
	m1 = map[string]int{"a": 1, "b": 2}
	m2 := make(map[string]int)
	m2["b"] = 3

	b := m1[""]
	fmt.Println(b)

	delete(m1, "a")
	fmt.Println(m1, m2)

	var m map[int]bool
	a, ok := m[1]
	fmt.Println(a, ok)

	// 线程安全map
	var m3 sync.Map
	m3.Store("a", 98)
	m3.Store("b", 98)

	// 遍历匿名函数
	m3.Range(func(k, v interface{}) bool {
		fmt.Println("k=", k, ", v=", v)
		return true
	})

	v, ifExists := m3.Load("a")
	if ifExists {
		fmt.Println("found value=", v)
	}
	v, ifExists = m3.LoadAndDelete("a")
	if ifExists {
		fmt.Println("load value=", v)
		_, ifExists = m3.Load("a")
		if !ifExists {
			fmt.Println("LoadAndDelete ok")
		}
	}
	m3.Delete("a")

	ListTest()

}

func ListTest() {

	ll := list.New()
	ll.PushFront(99)
	fmt.Println(ll)
	var l list.List

	l.PushBack(1)
	ele2 := l.PushBack(2)
	l.PushFront(3)
	l.PushBackList(ll)
	l.PushFrontList(ll)
	l.InsertAfter(88, ele2)
	l.InsertBefore(77, ele2)
	l.Remove(ele2)
	for i := l.Front(); i != nil; i = i.Next() {
		fmt.Println(i.Value)
	}

	for i := l.Back(); i != nil; i = i.Prev() {
		fmt.Println(i.Value)
	}
	fmt.Println(unsafe.Sizeof(l)) // 内存大小

	LoopTest()

}

func LoopTest() {

	var i int

	for i < 10 {
		if i < 10 {
			break
		}
	}

	switch i {
	case 2:
		fmt.Println("2")
	default:
		fmt.Println("default")
	}

	switch {
	case i < 3:
		fmt.Println("3")
		fallthrough // 继续判定
	case i < 2:
		fmt.Println("2")
	default:
		fmt.Println("default")
	}

	// 匿名函数定义调用
	func(x int) {
		fmt.Println("xxx", x)
	}(100)

	var f = func(x int) bool {
		return true
	}
	f(1)

	// 函数回调
	CallBackFunc([]int{1, 2, 3}, func(v int) {
		fmt.Println(v)
	})

	// function map
	var skill = map[string]func(){
		"fire": func() {
			fmt.Println("chicken fire")
		},
		"run": func() {
			fmt.Println("soldier run")
		},
		"fly": func() {
			fmt.Println("angel fly")
		},
	}

	skill["up"] = func() {
		fmt.Println("any up")
	}

	if f, ok := skill["fire"]; ok {
		f()
	} else {
		fmt.Println("skill not found")
	}

	getSequence := getSequence()
	fmt.Println(&getSequence)
	fmt.Println(getSequence())
	fmt.Println(getSequence())
	fmt.Println(getSequence())

	player1 := player("kute")
	fmt.Println(&player1)
	fmt.Println(player1())
	fmt.Println(player1())
	fmt.Println(player1())

	fmt.Println(reflect.TypeOf(player1))

	DeferTest()

}

func player(name string) func() (string, int) {
	blood := 100
	return func() (string, int) {
		blood--
		return name, blood
	}
}

type Doer interface {
	Do(p []byte) (n int, err error)
}

func CallBackFunc(list []int, f func(x int)) {
	for _, v := range list {
		f(v)
	}
}

/**
代码的延迟顺序与最终的执行顺序是反向的
延迟调用是在 defer 所在函数结束时进行，函数结束可以是正常返回时，也可以是发生宕机时
*/
func DeferTest() {
	fmt.Println("defer begin")

	defer fmt.Println("1")
	defer fmt.Println("2")
	defer fmt.Println("3")

	fmt.Println("defer end")

	var valueLock sync.Mutex
	valueLock.Lock()

	fmt.Println("ssss")

	// 最后释放锁
	defer valueLock.Unlock()
	fmt.Println("ssxx")

	StructTest()
}

type PlayerManager struct {
	x int
	y string
	z *int
}

func newPlayerManager(x int, y string, z *int) *PlayerManager {
	return &PlayerManager{
		x: x,
		y: y,
		z: z,
	}
}

func StructTest() {
	// 结构体实例化
	// 基本实例化, 普通实例
	var z int = 4
	var p1 PlayerManager
	p1.x = 1
	p1.z = &z
	fmt.Println(reflect.TypeOf(p1), p1.x)

	// 指针实例化
	p2 := new(PlayerManager)
	p2.x = 1
	fmt.Println(reflect.TypeOf(p2), p2.x)

	// & 取地址操作也会当做一次new的操作，取结构体的地址实例化
	p3 := &PlayerManager{}
	p3.x = 1
	fmt.Println(reflect.TypeOf(p3), (*p3).x) // p3.x 是语法糖，会被转化为 (*p3).x

	p4 := &PlayerManager{1, "a", &z}
	fmt.Println(reflect.TypeOf(p4), p4.x)

	// 函数包装实例化
	p5 := newPlayerManager(1, "a", &z)
	fmt.Println(reflect.TypeOf(p5), p5.x, *p5.z)
	*p5.z = 6
	fmt.Println(reflect.TypeOf(p5), p5.x, *p5.z)

	// 成员初始化，选择性初始化
	ins := &People{
		name: "x",
		// 结构体成员中只能包含结构体的指针类型，包含非指针类型会引起编译错误
		child: &People{
			name: "y",
			age:  uint8(18),
		},
	}
	fmt.Println(reflect.TypeOf(ins), ins.name, ins.child.name, (*(*ins).child).name)
	if c := ins.child; c != nil {
		fmt.Println(c.name)
	}

	// 使用闭包构造不同的结构体
	pn := NewPeopleByName("s")
	pa := NewPeopleByAge(18)
	fmt.Printf("%T\n%T\n", pn, pa)

	// 定义和初始化匿名结构体
	spu := &struct {
		name string
		age  int
	}{
		name: "k",
		age:  18,
	}
	fmt.Println(reflect.TypeOf(spu), spu.age)
	fmt.Printf("%T\n", spu)

	// 内嵌结构体
	pps := People{"sd", uint8(18), nil, 19, 22, PlayerManager{
		x: 0,
		y: "ff",
		z: &z,
	}, struct {
		Power int    // 功率
		Type  string // 类型
	}{
		Power: 44,
		Type:  "fssf",
	}}
	fmt.Printf("%T\n", pps)
	fmt.Println(pps.PlayerManager)
	fmt.Println(pps.x, pps.PlayerManager.x)

	IoTest()

}

func NewPeopleByName(name string) *People {
	return &People{
		name: name,
	}
}

func NewPeopleByAge(age int) *People {
	return &People{
		age: uint8(age),
	}
}

type People struct {
	name  string
	age   uint8
	child *People
	x     int
	int   // 匿名字段，在一个结构体中对于每一种数据类型只能有一个匿名字段
	// 内嵌的结构体可以直接访问其成员变量,一个结构体只能嵌入一个同类型的成员
	PlayerManager // 内嵌结构体，类比于 继承，拥有了 PlayerManager 的字段，同名字段还是需要区分的

	// 匿名结构体
	Engine struct {
		Power int    // 功率
		Type  string // 类型
	}
}

var zeroError = errors.New("zero error")

func fileSize(filename string) int64 {

	f, err := os.Open(filename)

	if err != nil {
		return 0
	}

	// 延迟调用Close, 此时Close不会被调用
	defer f.Close()

	info, err := f.Stat()

	if err != nil {
		// defer机制触发, 调用Close关闭文件
		return 0
	}

	size := info.Size()

	// defer机制触发, 调用Close关闭文件
	return size
}

func main() {
	//hello
	/*sdf*/
	fmt.Println("Hello World!")

	// _ 占位符
	_, b, d := testFmt()
	fmt.Println(b, d)
}
