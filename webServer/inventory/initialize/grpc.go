package initialize

import (
	"fmt"
	"time"

	retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"webServer/common/otgrpc"
	"webServer/inventory/global"
	"webServer/inventory/proto"
)

func InitGrpcClient() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`))
	opts = append(opts, grpc.WithUnaryInterceptor(retry.UnaryClientInterceptor(
		retry.WithMax(3),
		retry.WithCodes(codes.Unknown, codes.Unavailable, codes.DeadlineExceeded),
		retry.WithPerRetryTimeout(1*time.Second),
	)))
	opts = append(opts, grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())))
	client, err := grpc.Dial(
		fmt.Sprintf(
			"consul://%s:%d/%s?wait=14s&tag=gpmall",
			global.ServerConfig.Consul.Host,
			global.ServerConfig.Consul.Port,
			global.ServerConfig.GrpcClient.Name,
		),
		opts...,
	)
	if err != nil {
		zap.S().Fatalf("grpc 连接失败: %s", err.Error())
	}

	global.GrpcClient = proto.NewInventoryClient(client)
}
