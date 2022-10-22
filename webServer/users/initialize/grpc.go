package initialize

import (
	"fmt"
	"time"

	retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/hashicorp/consul/api"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"webServer/common/otgrpc"
	"webServer/users/global"
	"webServer/users/proto"
)

/*
  grpc-consul-resolver 这个包导入即用，不需要关心任何实现过程
  是基于grpc通信解析开发的解析器
  每次进行函数请求时，都会通过consul来转发到服务端，而不是直接和服务端通信
*/
func GrpcBalancer() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`))
	opts = append(opts, grpc.WithUnaryInterceptor(retry.UnaryClientInterceptor(
		retry.WithMax(3),
		retry.WithPerRetryTimeout(1*time.Second),
		retry.WithCodes(codes.Unknown, codes.Unavailable, codes.DeadlineExceeded),
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
	// defer conn.Close()

	if err != nil {
		zap.S().DPanic(err.Error())
		panic(err.Error())
	}
	global.UserClient = proto.NewUserClient(conn)
}

/*
  可删除此方法
  这个方法是基于http请求的，如果每次请求时，都调用这个方法，会损失一定性能。
  此方法没有负载均衡，己被 GrpcBalancer() 替代
*/
func GrpcClient() {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.Consul.Host, global.ServerConfig.Consul.Port)
	client, err := api.NewClient(cfg)
	if err != nil {
		zap.S().Warn("connect consul fail")
	}

	services, err := client.Agent().ServicesWithFilter(fmt.Sprintf(`Service == "%s"`, global.ServerConfig.GrpcAddr.Name))
	if err != nil {
		zap.S().Warn("get consul service fail")
		return
	}
	if len(services) == 0 {
		zap.S().Warn(fmt.Sprintf(`not found Service == "%s":%v`, global.ServerConfig.GrpcAddr.Name, services))
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

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		zap.S().Warnf("远程连接错误：%s", err.Error())
	}
	global.UserClient = proto.NewUserClient(conn)
}
