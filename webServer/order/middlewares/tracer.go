package middlewares

import (
	"fmt"
	"webServer/order/global"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

func Trace() gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg := jaegercfg.Configuration{
			Sampler: &jaegercfg.SamplerConfig{
				Type:  jaeger.SamplerTypeConst,
				Param: 1,
			},
			Reporter: &jaegercfg.ReporterConfig{
				LocalAgentHostPort: fmt.Sprintf("%s:%d", global.ServerConfig.Jaeger.Host, global.ServerConfig.Jaeger.Port),
				LogSpans:           true,
			},
			ServiceName: global.ServerConfig.Jaeger.Name,
		}
		tracer, closer, err := cfg.NewTracer()
		if err != nil {
			panic(err)
		}
		defer closer.Close()
		opentracing.SetGlobalTracer(tracer)

		span := tracer.StartSpan(c.Request.URL.Path)
		defer span.Finish()

		// 不用opentracing.GetGlobalTracer() 的目的是为了不同 api网关group 下 tracer的独立性
		c.Set("tracer", tracer)
		c.Set("parentSpan", span)

		c.Next()
	}
}
