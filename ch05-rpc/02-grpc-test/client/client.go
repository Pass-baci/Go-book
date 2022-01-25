package main

import (
	"Go-/ch05-rpc/02-grpc-test/proto"
	"context"
	"fmt"
	"google.golang.org/grpc"
)

func main() {
	serviceAddress := "127.0.0.1:1234"
	conn, err := grpc.Dial(serviceAddress, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	// 获取一个gRPC客户端
	stringClient := proto.NewStringServiceClient(conn)
	stringReq := &proto.StringRequest{
		A: "AB",
		B: "AB",
	}
	reply, err := stringClient.Diff(context.Background(), stringReq)
	if err != nil {
		panic(err)
	}
	fmt.Println(reply.Result, reply.Err)
}
