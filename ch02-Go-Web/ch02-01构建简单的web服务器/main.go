package main

import (
	"fmt"
	"log"
	"net/http"
)

func sayHello(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()                 // 解析参数
	fmt.Println(r.Form)               // 传入的参数
	fmt.Println("Path: ", r.URL.Path) // 访问路径
	fmt.Println("Host: ", r.Host)     // 访问地址
	for k, v := range r.Form {
		fmt.Println("k: ", k)
		fmt.Println("k: ", v)
	}
	fmt.Fprintf(w, "hello world, %s", r.Form.Get("name"))
}

func main() {
	http.HandleFunc("/", sayHello)                            // 访问路由
	if err := http.ListenAndServe(":8080", nil); err != nil { // 监听服务端口
		log.Fatalln(err)
	}
}
