package initialize

import (
	"fmt"
	"time"

	retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	_ "github.com/mbobakov/grpc-consul-resolver" // consul协议解析器
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"webServer/common/otgrpc"
	"webServer/userop/global"
	"webServer/userop/proto"
)

func InitGrpcClient() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithDefaultServiceConfig(`{"loadBalancePocily":"round_robin"}`))
	opts = append(opts, grpc.WithUnaryInterceptor(retry.UnaryClientInterceptor(
		retry.WithMax(3),
		retry.WithCodes(codes.Unknown, codes.DeadlineExceeded, codes.Unavailable),
		retry.WithPerRetryTimeout(1*time.Second),
	)))
	opts = append(opts, grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())))
	conn, err := grpc.Dial(
		fmt.Sprintf(
			"consul://%s:%d/%s?wait=14s&tag=gpmall",
			global.ServerConfig.Consul.Host,
			global.ServerConfig.Consul.Port,
			global.ServerConfig.GrpcClient.Name,
		),
		opts...,
	)
	if err != nil {
		zap.S().DPanic(err.Error())
		panic(err.Error())
	}

	global.GrpcFavClient = proto.NewFavoritesClient(conn)
	global.GrpcAddressClient = proto.NewAddressClient(conn)
	global.GrpcPostClient = proto.NewPostClient(conn)
}
