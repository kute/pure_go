package main

import (
	"fmt"
	"reflect"
)

func main() {

	numbers := make(map[string]int)
	fmt.Println(reflect.TypeOf(numbers)) // map[string]int

}
