package initialize

import (
	"time"
	"webServer/common/otgrpc"
	"webServer/oss/global"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/opentracing/opentracing-go"
)

func InitGrpcClient() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`))
	opts = append(opts, grpc.WithUnaryInterceptor(retry.UnaryClientInterceptor(
		retry.WithMax(3),
		retry.WithPerRetryTimeout(1*time.Second),
		retry.WithCodes(codes.Unknown, codes.DeadlineExceeded, codes.Unavailable),
	)))
	opts = append(opts, grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())))
	client, err := grpc.Dial(
		"consul:192.168.8.222:8500?wait=14s&tag=gpmall",
		opts...,
	)
	if err != nil {
		zap.S().Warnf("gprc err:", err.Error())
	}

	global.ClientConn = client
}
