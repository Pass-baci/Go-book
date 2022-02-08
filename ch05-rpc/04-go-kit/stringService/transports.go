package stringService

import (
	"Go-/ch05-rpc/04-go-kit/proto"
	"context"
	"github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	concat grpc.Handler
	diff   grpc.Handler
}

func (s *grpcServer) Concat(ctx context.Context, r *proto.StringRequest) (*proto.StringResponse, error) {
	_, resp, err := s.concat.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*proto.StringResponse), nil
}

func (s *grpcServer) Diff(ctx context.Context, r *proto.StringRequest) (*proto.StringResponse, error) {
	_, resp, err := s.diff.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*proto.StringResponse), nil
}

func NewStringServer(ctx context.Context, endpoints StringEndpoints) proto.StringServiceServer {
	return &grpcServer{
		concat: grpc.NewServer(
			endpoints.StringEndpoint,
			DecodeConcatStringRequest,
			EncodeStringResponse,
		),
		diff: grpc.NewServer(
			endpoints.StringEndpoint,
			DecodeDiffStringRequest,
			EncodeStringResponse,
		),
	}
}
