package service

import (
	"errors"
	"strings"
)

// 固定字符串长度
const STRMAXSIZE = 1024

// 定义Service错误类型
var (
	ErrMaxSize  = errors.New("maximum size of 1024 bytes exceeded")
	ErrStrValue = errors.New("maximum size of 1024 bytes exceeded")
)

// 定义服务接口
type Service interface {
	// 连接字符串
	Concat(a, b string) (string, error)
	// 获取字符串公共字符
	Diff(a, b string) (string, error)
	// 健康检查
	HealthCheck() bool
}

type StringService struct {
}

func (s StringService) Concat(a, b string) (string, error) {
	if len(a)+len(b) > STRMAXSIZE {
		return "", ErrMaxSize
	}
	return a + b, nil
}
func (s StringService) Diff(a, b string) (string, error) {
	if len(a) < 1 || len(b) < 1 {
		return "", nil
	}
	res := ""
	if len(a) >= len(b) {
		for _, char := range b {
			if strings.Contains(a, string(char)) {
				res = res + string(char)
			}
		}
	} else {
		for _, char := range a {
			if strings.Contains(b, string(char)) {
				res = res + string(char)
			}
		}
	}
	return res, nil
}

func (s StringService) HealthCheck() bool {
	return true
}

type ServiceMiddleware func(Service) Service
