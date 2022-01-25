package main

import (
	"Go-/ch05-rpc/02-grpc-test/proto"
	"context"
	"errors"
	"google.golang.org/grpc"
	"net"
	"strings"
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
func (s *StringService) Diff(ctx context.Context, req *proto.StringRequest) (*proto.StringResponse, error) {
	var ret proto.StringResponse
	if len(req.A) < 1 || len(req.B) < 1 {
		ret.Result = ""
		return &ret, nil
	}
	if len(req.A) >= len(req.B) {
		for _, char := range req.B {
			if strings.Contains(req.A, string(char)) {
				ret.Result = ret.Result + string(char)
			}
		}
	} else {
		for _, char := range req.A {
			if strings.Contains(req.B, string(char)) {
				ret.Result = ret.Result + string(char)
			}
		}
	}
	return &ret, nil
}
