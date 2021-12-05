package main

import "fmt"

func main() {

	var arr = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	var sum = 9

	/**
	数组中，两数之和等于目标，返回首次出现的元素组合
	*/
	var twoSumFunc = func(arr []int, sum int) (int, int) {
		var bucket = make(map[int]int, len(arr))
		for i := 0; i < len(arr); i++ {
			// 过滤掉过大的数
			if arr[i] > sum {
				continue
			}
			var sub = sum - arr[i]
			if _, ifExists := bucket[sub]; ifExists {
				return sub, arr[i]
			}
			bucket[arr[i]] = i
		}
		return -1, -1
	}
	fmt.Println(twoSumFunc(arr, sum))

	/**
	三数之和
	*/
	func(arr []int, sum int) {
		//TODO
	}(arr, sum)

	/**
	数组中所有的组合等于sum，返回组合列表
	*/
	func(arr []int, sum int) {
		//TODO
	}(arr, sum)

}
