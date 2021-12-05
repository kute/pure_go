package main

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"sort"
)

type Kuter interface {
	Done(x interface{}) (y int, err error)
}

type Whyer interface {
	省略变量名(int) (int, error)
}

type Param struct {
	data interface{}
}

type Meter struct {
	data interface{}
}

// 实现接口 Kuter，函数签名一致
func (s *Meter) Done(x interface{}) (y int, err error) {
	fmt.Println(reflect.TypeOf(s))
	return 0, nil
}

// 未实现接口
func (s *Meter) DoneX(x interface{}) (y int, err error) {
	fmt.Println(reflect.TypeOf(s))
	return 0, nil
}

// 实现接口 Whyer
func (s *Meter) 省略变量名(x int) (y int, err error) {
	fmt.Println(reflect.TypeOf(s))
	return 0, nil
}

// 实现接口
func (s *Param) Done(x interface{}) (y int, err error) {
	fmt.Println(reflect.TypeOf(s))
	if x == 11 {
		s.data = x
	}
	return 0, nil
}

// 接口作为参数，并不关心实现
func usingKuter(kuter Kuter, x interface{}) {
	_, _ = kuter.Done(x)
}

func usingWhyer(whyer Whyer, x interface{}) {
	// 注意，因为 Whyer 接口的方法 『省略变量名』的入参是 int，而这里的x的类型是 interface{}，所以需要这里使用断言，若断言失败，则触发panic
	whyer.省略变量名(x.(int))

	// 断言判断，若 这里不接收ok的值失败时也会panic
	if intValue, ok := x.(int); ok {
		whyer.省略变量名(intValue)
	}

	// 断言 配合switch
	switch x.(type) {
	case int, string:
		fmt.Println("int, string")
	default:
		fmt.Println("default")
	}
}

//实现排序接口 sort.Interface，需要实现接口里的所有方法
type SortObject struct {
	seq      []int
	reversed bool
}

func (o SortObject) Len() int {
	return len(o.seq)
}

// 默认 i < j，当 Less 为true时调用Swap
func (o SortObject) Less(i, j int) bool {
	var seq = o.seq
	if o.reversed {
		return seq[i] > seq[j]
	}
	return seq[i] < seq[j]
}

func (o SortObject) Swap(i, j int) {
	var seq = o.seq
	seq[i], seq[j] = seq[j], seq[i]
}

func main() {

	fmt.Println("begin")

	var ins = &Param{}
	ins.Done(11)

	var pns = &Meter{}
	pns.Done(11)
	pns.DoneX(11)
	pns.省略变量名(11)

	var kns Kuter = ins
	kns.Done(11)

	var mns Whyer = pns
	mns.省略变量名(11)

	// 排序测试
	var o = &SortObject{
		seq: []int{4, 3, 5, 2, 1, 7},
		// 自定义升降序
		reversed: true,
	}

	// 使用自定义的排序
	fmt.Println("IsSorted=", sort.IsSorted(o), "IntsAreSorted=", sort.IntsAreSorted(o.seq))
	sort.Sort(o)
	fmt.Println("IsSorted=", sort.IsSorted(o), "IntsAreSorted=", sort.IntsAreSorted(o.seq))
	fmt.Println(o.seq)

	// 反序
	sort.Sort(sort.Reverse(o))
	fmt.Println(o.seq)

	// 内建的排序方法
	sort.Ints(o.seq)
	fmt.Println(o.seq)
	// 自定义less函数
	sort.Slice(o.seq, func(i, j int) bool {
		return o.seq[i] > o.seq[j]
	})
	fmt.Println(o.seq)

	var w io.Writer
	w = os.Stdout
	if rw, ok := w.(io.ReadWriter); ok {
		fmt.Println(reflect.TypeOf(rw))
	} else {
		fmt.Println("dd")
	}

}
