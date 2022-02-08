package stringService

// 负责请求体与响应体的数据转换

import (
	"Go-/ch05-rpc/04-go-kit/proto"
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"strings"
)

type StringEndpoints struct {
	StringEndpoint      endpoint.Endpoint
	HealthCheckEndpoint endpoint.Endpoint
}

func (ue StringEndpoints) Concat(ctx context.Context, a string, b string) (string, error) {
	//ctx := context.Background()
	resp, err := ue.StringEndpoint(ctx, &proto.StringRequest{
		A: a,
		B: b,
	})
	response := resp.(*proto.StringResponse)
	return response.Result, err
}

func (ue StringEndpoints) Diff(ctx context.Context, a string, b string) (string, error) {
	//ctx := context.Background()
	resp, err := ue.StringEndpoint(ctx, proto.StringRequest{
		A: a,
		B: b,
	})
	response := resp.(proto.StringResponse)
	return response.Result, err
}

var (
	ErrInvalidRequestType = errors.New("RequestType has only two type: Concat, Diff")
)

// StringRequest define request struct
type StringRequest struct {
	RequestType string `json:"request_type"`
	A           string `json:"a"`
	B           string `json:"b"`
}

// StringResponse define response struct
type StringResponse struct {
	Result string `json:"result"`
	Error  error  `json:"error"`
}

func MakeStringEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(StringRequest)

		var (
			res, a, b string
			opError   error
		)

		a = req.A
		b = req.B

		if strings.EqualFold(req.RequestType, "Concat") {
			res, _ = svc.Concat(ctx, a, b)
		} else if strings.EqualFold(req.RequestType, "Diff") {
			res, _ = svc.Diff(ctx, a, b)
		} else {
			return nil, ErrInvalidRequestType
		}

		return StringResponse{Result: res, Error: opError}, nil
	}
}

// HealthRequest 健康检查请求结构
type HealthRequest struct{}

// HealthResponse 健康检查响应结构
type HealthResponse struct {
	Status bool `json:"status"`
}

// MakeHealthCheckEndpoint 创建健康检查Endpoint
func MakeHealthCheckEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		status := true
		return HealthResponse{status}, nil
	}
}

func DecodeConcatStringRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(*proto.StringRequest)
	return StringRequest{
		RequestType: "Concat",
		A:           string(req.A),
		B:           string(req.B),
	}, nil
}

func DecodeDiffStringRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(*proto.StringRequest)
	return StringRequest{
		RequestType: "Diff",
		A:           string(req.A),
		B:           string(req.B),
	}, nil
}

func EncodeStringResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(StringResponse)

	if resp.Error != nil {
		return &proto.StringResponse{
			Result: resp.Result,
			Err:    resp.Error.Error(),
		}, nil
	}

	return &proto.StringResponse{
		Result: resp.Result,
		Err:    "",
	}, nil
}
