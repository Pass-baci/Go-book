package main

import (
	"fmt"
	"log"
	"net/rpc"
)

type StringRequest struct {
	A string
	B string
}

func main() {
	client, err := rpc.DialHTTP("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Fatal("dialing: ", err)
	}
	stringReq := &StringRequest{"A", "BA"}
	var reply string
	// 同步调用
	if err = client.Call("StringService.Concat", stringReq, &reply); err != nil {
		log.Fatal("Concat error: ", err)
	}
	fmt.Println(*stringReq, reply)
	// 异步调用
	call := client.Go("StringService.Diff", stringReq, &reply, nil)
	fmt.Println(<-call.Done)
	fmt.Println(*stringReq, reply)
}
