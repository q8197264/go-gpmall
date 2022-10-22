package initialize

import (
	"fmt"
	"time"

	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"

	"webServer/common/otgrpc"
	"webServer/goods/global"
	"webServer/goods/proto"
)

// grpc连接池 [负载均衡]
func GrpcBalancer() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`))
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
			global.ServerConfig.GrpcAddr.Name,
		),
		opts...,
	)
	if err != nil {
		zap.S().Fatalf("grpc连接consul失败%s", err.Error())
	}
	// defer conn.Close()

	global.GoodsClient = proto.NewGoodsClient(conn)
}
