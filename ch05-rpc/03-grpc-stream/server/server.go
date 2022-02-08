package main

import (
	"Go-/ch05-rpc/03-grpc-stream/proto"
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"net"
	"strings"
	"time"
)

func main() {
	l, err := net.Listen("tcp", "127.0.0.1:1234")
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()
	stringService := new(StringService)
	proto.RegisterStringServiceServer(grpcServer, stringService)
	if err = grpcServer.Serve(l); err != nil {
		panic(err)
	}
}

type StringService struct {
}

func (s *StringService) Concat(ctx context.Context, req *proto.StringRequest) (*proto.StringResponse, error) {
	var ret string
	if len(req.A)+len(req.B) > 1024 {
		return &proto.StringResponse{
			Result: ret,
			Err:    errors.New("max size is 1024").Error(),
		}, nil
	}
	ret = req.A + req.B
	return &proto.StringResponse{
		Result: ret,
	}, nil
}

func (s *StringService) LotsServerStream(req *proto.StringRequest, qs proto.StringService_LotsServerStreamServer) error {
	response := proto.StringResponse{Result: req.A + req.B}
	for i := 0; i < 10; i++ {
		qs.Send(&response)
		time.Sleep(time.Second)
	}
	return nil
}
func (s *StringService) LotsClientStream(qs proto.StringService_LotsClientStreamServer) error {
	var params []string
	for {
		in, err := qs.Recv()
		if err == io.EOF {
			qs.SendAndClose(&proto.StringResponse{Result: strings.Join(params, "")})
			return nil
		}
		if err != nil {
			return err
		}
		params = append(params, in.A, in.B)
		fmt.Println(params)
		time.Sleep(time.Second)
	}
}
func (s *StringService) LotsAllStream(qs proto.StringService_LotsAllStreamServer) error {
	var params []string
	for {
		time.Sleep(time.Second)
		in, err := qs.Recv()
		if err == io.EOF {
			return nil
		}
		qs.Send(&proto.StringResponse{Result: in.A + in.B})
		params = append(params, in.A, in.B)
		fmt.Println(params)
	}
}
