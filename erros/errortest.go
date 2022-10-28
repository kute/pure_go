package main

import (
	"errors"
	"fmt"
)

func main() {
	fmt.Println("hello")

	var error = errors.New("this is a error message")
	fmt.Println(error)

	const name, id = "bimmler", 17
	err := fmt.Errorf("user %q (id %d) not found", name, id)
	if err != nil {
		fmt.Print(err)
	}

}
