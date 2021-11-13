package main

import (
	"io/ioutil"
	"log"
	"net/http"
	_ "os" // 匿名引用，如果包有init函数会执行，没有也不会报错
)

// 包别名
import F "fmt"

//这种格式相当于把 fmt 包直接合并到当前程序中，在使用 fmt 包内的方法是可以不用加前缀fmt.，直接引用
//import . "fmt"

// 初始化函数
//在运行时，被最后导入的包会最先初始化并调用 init() 函数
func init() {
	F.Println("ini before main")
}

func main() {

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(200)
		_, _ = F.Fprintln(writer, "root page")
	})

	http.HandleFunc("/index", func(writer http.ResponseWriter, request *http.Request) {
		if content, err := ioutil.ReadFile("./http/index.html"); err == nil {
			_, _ = writer.Write(content)
		}
	})
	log.Fatal(http.ListenAndServe("localhost:8080", nil))

}
