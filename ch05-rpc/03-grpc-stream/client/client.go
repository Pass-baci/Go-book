package main

import (
	"Go-/ch05-rpc/03-grpc-stream/proto"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
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
	//stream, err := stringClient.LotsServerStream(context.Background(), stringReq)
	//if err != nil {
	//	panic(err)
	//}
	//// 服务端流RPC调用
	//for {
	//	response, err := stream.Recv()
	//	if err == io.EOF {
	//		break
	//	}
	//	if err != nil {
	//		panic(err)
	//	}
	//	fmt.Println(response.Result)
	//}
	//if err = stream.CloseSend(); err != nil {
	//	panic(err)
	//}
	//// 客户端流RPC调用
	//streamClient, err := stringClient.LotsClientStream(context.Background())
	//if err != nil {
	//	panic(err)
	//}
	//for i:=0; i<10; i++ {
	//	streamClient.Send(stringReq)
	//}
	//if streamClient.CloseAndRecv(); err != nil {
	//	panic(err)
	//}
	// 双向流RPC调用
	streamAllClient, err := stringClient.LotsAllStream(context.Background())
	if err != nil {
		panic(err)
	}

	for i := 0; i < 10; i++ {
		err := streamAllClient.Send(stringReq)
		if err != nil {
			panic(err)
		}
		response, err := streamAllClient.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		fmt.Println(response.Result)
	}
	if streamAllClient.CloseSend(); err != nil {
		panic(err)
	}
}
