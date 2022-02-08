package main

import (
	"Go-/ch05-rpc/02-grpc-test/proto"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"time"
)

func main() {
	interceptor := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		start := time.Now()
		md := metadata.New(map[string]string{
			"appid": "10086",
			"password": "123456",
		})
		ctx = metadata.NewOutgoingContext(context.Background(), md)
		err := invoker(ctx, method, req, reply, cc, opts...)
		fmt.Printf("耗时: %s\n", time.Since(start))
		return err
	}
	serviceAddress := "127.0.0.1:1234"
	conn, err := grpc.Dial(serviceAddress, grpc.WithInsecure(), grpc.WithUnaryInterceptor(interceptor))
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
