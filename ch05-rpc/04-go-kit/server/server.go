package main

import (
	"Go-/ch05-rpc/04-go-kit/proto"
	"Go-/ch05-rpc/04-go-kit/stringService"
	"context"
	"flag"
	"github.com/go-kit/kit/log"
	"google.golang.org/grpc"
	"net"
	"os"
)

func main() {
	flag.Parse()

	ctx := context.Background()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var svc stringService.Service
	svc = stringService.StringService{}

	// add logging middleware
	svc = stringService.LoggingMiddleware(logger)(svc)

	endpoint := stringService.MakeStringEndpoint(svc)

	//创建健康检查的Endpoint
	healthEndpoint := stringService.MakeHealthCheckEndpoint(svc)

	//把算术运算Endpoint和健康检查Endpoint封装至StringEndpoints
	endpts := stringService.StringEndpoints{
		StringEndpoint:      endpoint,
		HealthCheckEndpoint: healthEndpoint,
	}

	handler := stringService.NewStringServer(ctx, endpts)

	ls, _ := net.Listen("tcp", "127.0.0.1:8080")
	gRPCServer := grpc.NewServer()
	proto.RegisterStringServiceServer(gRPCServer, handler)
	gRPCServer.Serve(ls)
}
