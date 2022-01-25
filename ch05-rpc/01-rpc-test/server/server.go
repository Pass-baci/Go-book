package main

import (
	"errors"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"strings"
)

const STRMAXSIZE = 1024

// 定义Service错误类型
var (
	ErrMaxSize  = errors.New("maximum size of 1024 bytes exceeded")
	ErrStrValue = errors.New("maximum size of 1024 bytes exceeded")
)

type StringRequest struct {
	A string
	B string
}

type Service interface {
	// 连接字符串
	Concat(a, b string) (string, error)
	// 获取字符串公共字符
	Diff(a, b string) (string, error)
}

type StringService struct {
}

func (s *StringService) Concat(req StringRequest, ret *string) error {
	if len(req.A)+len(req.B) > STRMAXSIZE {
		return ErrMaxSize
	}
	*ret = req.A + req.B
	return nil
}
func (s *StringService) Diff(req StringRequest, ret *string) error {
	if len(req.A) < 1 || len(req.B) < 1 {
		return nil
	}
	if len(req.A) >= len(req.B) {
		for _, char := range req.B {
			if strings.Contains(req.A, string(char)) {
				*ret = *ret + string(char)
			}
		}
	} else {
		for _, char := range req.A {
			if strings.Contains(req.B, string(char)) {
				*ret = *ret + string(char)
			}
		}
	}
	return nil
}

func main() {
	stringService := new(StringService)
	registerError := rpc.Register(stringService)
	if registerError != nil {
		log.Fatal("Register error: ", registerError)
	}
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Fatal("listen error:", err)
	}
	http.Serve(l, nil)
}
