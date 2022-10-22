package initialize

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/consul/api"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"

	"webServer/common/otgrpc"
	"webServer/order/global"
	"webServer/order/proto"
)

func InitGrpcClient() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithDefaultServiceConfig(`{"loadBalancePocily":"round_robin"}`))
	// opts = append(opts, grpc.WithUnaryInterceptor(func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	// 	start := time.Now()
	// 	err := invoker(ctx, method, req, reply, cc, opts...)
	// 	fmt.Printf("-------耗时:  %s\n", time.Since(start))
	// 	return err
	// }))
	// 这个请求应该多少时间超时, 这个重试应该几次, 当服务器返回什么状态码的时候重试
	opts = append(opts, grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(
		grpc_retry.WithMax(3),
		grpc_retry.WithPerRetryTimeout(1*time.Second),
		grpc_retry.WithCodes(codes.Unknown, codes.DeadlineExceeded, codes.Unavailable),
	)))
	opts = append(opts, grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())))
	client, err := grpc.Dial(
		fmt.Sprintf(
			"consul://%s:%d/%s?wait=14s&tag=gpmall",
			global.ServerConfig.Consul.Host,
			global.ServerConfig.Consul.Port,
			global.ServerConfig.Grpc.Name,
		),
		opts...,
	)
	if err != nil {
		zap.S().DPanic(err.Error())
	}

	global.OrderClient = proto.NewOrderClient(client)
	global.ShopCartClient = proto.NewShopCartClient(client)
}

/*
  可删除此方法
  这个方法是基于http请求的，如果每次请求时，都调用这个方法，会损失一定性能。
  此方法没有负载均衡，己被 GrpcBalancer() 替代
*/
func InitGrpcClient2() {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.Consul.Host, global.ServerConfig.Consul.Port)
	client, err := api.NewClient(cfg)
	if err != nil {
		zap.S().Warn("connect consul fail")
	}

	services, err := client.Agent().ServicesWithFilter(fmt.Sprintf(`Service == "%s"`, global.ServerConfig.Grpc.Name))
	if err != nil {
		zap.S().Warn("get consul service fail")
		return
	}
	if len(services) == 0 {
		zap.S().Warn(fmt.Sprintf(`not found Service == "%s":%v`, global.ServerConfig.Grpc.Name, services))
		return
	}

	var address string
	for _, v := range services {
		address = fmt.Sprintf("%s:%d", v.Address, v.Port)
		break
	}
	if len(address) == 0 {
		zap.S().Warnf("获取注册中心服务失败：%s", err.Error())
	}

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithDefaultServiceConfig(`{"loadBalancePocily":"round_robin"}`))
	opts = append(opts, grpc.WithUnaryInterceptor(func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		start := time.Now()
		fmt.Printf("-------耗时:  %s\n", time.Since(start))
		err := invoker(ctx, method, req, reply, cc, opts...)
		fmt.Printf("-------耗时:  %s\n", time.Since(start))

		return err
	}))
	var retryOpts = []grpc_retry.CallOption{
		grpc_retry.WithMax(3),
		grpc_retry.WithPerRetryTimeout(2 * time.Second),
		grpc_retry.WithCodes(codes.Unknown, codes.DeadlineExceeded, codes.Unavailable),
	}
	// 这个请求应该多少时间超时, 这个重试应该几次, 当服务器返回什么状态码的时候重试
	opts = append(opts, grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(retryOpts...)))

	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		zap.S().Warnf("远程连接错误：%s", err.Error())
	}
	// global.UserClient = proto.NewUserClient(conn)
	global.OrderClient = proto.NewOrderClient(conn)
	global.ShopCartClient = proto.NewShopCartClient(conn)
}
