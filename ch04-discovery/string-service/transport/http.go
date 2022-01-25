package transport

import (
	"Go-/ch04-discovery/string-service/endpoint"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var ErrorBadRequest = errors.New("invalid request parameter")

func decodesStringRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	requestType, ok := vars["type"]
	if !ok {
		return nil, ErrorBadRequest
	}
	pa, ok := vars["a"]
	if !ok {
		return nil, ErrorBadRequest
	}
	pb, ok := vars["b"]
	if !ok {
		return nil, ErrorBadRequest
	}
	return endpoint.StringRequest{
		RequestType: requestType,
		A:           pa,
		B:           pb,
	}, nil
}

func decodesHealthCheckRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return endpoint.HealthRequest{}, nil
}

// 自定义错误响应
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	switch err {
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

// 将结果转为json
func encodesStringResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func MakeHttpHandler(ctx context.Context, endpoints endpoint.StringEndpoints, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	options := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
	}
	r.Methods("POST").Path("/op/{type}/{a}/{b}").Handler(kithttp.NewServer(
		endpoints.StringEndpoint,
		decodesStringRequest,
		encodesStringResponse,
		options...,
	))
	r.Path("/metrics").Handler(promhttp.Handler())
	r.Methods("GET").Path("/health").Handler(kithttp.NewServer(
		endpoints.HealthCheckEndpoint,
		decodesHealthCheckRequest,
		encodesStringResponse,
		options...,
	))
	return r
}
