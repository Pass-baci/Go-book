package main

import (
	"Go-/ch05-rpc/02-grpc-test/proto"
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"net"
	"strings"
)

func main() {
	interceptor := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		fmt.Println("拦截到一个请求")
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return resp, status.Error(codes.Unauthenticated, "token无效")
		}
		var (
			appid string
			password string
		)
		if va1, ok := md["appid"]; ok{
			appid = va1[0]
		}
		if va1, ok := md["password"]; ok{
			password = va1[0]
		}
		if appid != "10086" || password != "123456" {
			return resp, status.Error(codes.Unauthenticated, "token无效")
		}
		fmt.Println("请求成功")
		return handler(ctx, req)}
	opt := []grpc.ServerOption{grpc.UnaryInterceptor(interceptor)}
	l, err := net.Listen("tcp", "127.0.0.1:1234")
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer(opt...)
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
