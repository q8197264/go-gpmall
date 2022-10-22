package middlewares

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"

	"webServer/userop/global"
)

func Trace() gin.HandlerFunc {
	return func(ctx *gin.Context) {
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
		tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
		if err != nil {
			panic(err)
		}
		defer closer.Close()
		opentracing.SetGlobalTracer(tracer)

		span := tracer.StartSpan(ctx.Request.URL.Path)
		defer span.Finish()

		ctx.Set("tracer", tracer)
		ctx.Set("parentSpan", span)

		ctx.Next()
	}
}
