package endpoint

import (
	"Go-/ch04-discovery/string-service/service"
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
)

// 定义服务接口端点
type StringEndpoints struct {
	StringEndpoint      endpoint.Endpoint
	HealthCheckEndpoint endpoint.Endpoint
}

var ErrInvalidRequestType = errors.New("RequestType has only two type: Concat, Diff")

type StringRequest struct {
	RequestType string `json:"request_type"`
	A           string `json:"a"`
	B           string `json:"b"`
}

type StringResponse struct {
	Result string `json:"result"`
	Error  error  `json:"error"`
}

func MakeStringEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(StringRequest)
		var (
			res, a, b string
			opError   error
		)
		a = req.A
		b = req.B
		switch req.RequestType {
		case "Concat":
			res, opError = svc.Concat(a, b)
		case "Diff":
			res, opError = svc.Diff(a, b)
		default:
			return nil, ErrInvalidRequestType
		}
		return StringResponse{
			Result: res,
			Error:  opError,
		}, nil
	}
}

type HealthRequest struct {
}

type HealthResponse struct {
	Status bool `json:"status"`
}

func MakeHealthCheckEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		status := svc.HealthCheck()
		return HealthResponse{
			Status: status,
		}, nil
	}
}
