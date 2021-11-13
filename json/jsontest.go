package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
)

type jsonobject struct {
	A  int
	B  string
	C  bool
	AA []int
	// 字段首字母小写不会被序列化
	privateField string
	MM           map[string]interface{}
	nestobject
}

type nestobject struct {
	Name string
}

func main() {

	o := jsonobject{
		A:            1,
		B:            "ss",
		C:            true,
		AA:           []int{2, 4, 1},
		privateField: "nothing",
		MM: map[string]interface{}{
			"k": 1,
		},
		nestobject: nestobject{
			Name: "kk",
		},
	}

	// object to json
	if content, err := json.Marshal(o); err == nil {
		s := string(content)
		fmt.Println(s)
		_, _ = os.Stdout.Write(content)
		fmt.Println()

		// json to object
		var ob jsonobject
		_ = json.Unmarshal(content, &ob)
		fmt.Println(ob)

		var ob2 = new(jsonobject)
		_ = json.Unmarshal(content, ob2)
		fmt.Println(*ob2)
	}

	fmt.Println(os.Hostname())
	fmt.Println(os.UserHomeDir())
	fmt.Println(os.Environ())
	fmt.Println(os.Getgid())
	fmt.Println(os.Getgroups())
	fmt.Println(os.Getwd())
	fmt.Println(os.Geteuid())
	fmt.Println(os.Getpagesize())
	fmt.Println(os.Getpid())
	fmt.Println(os.Getuid())
	fmt.Println(os.Getenv("GO111MODULE")) // go module是否开启
	fmt.Println(user.Current())
	os.Exit(0) // 正常执行，非0异常

}
