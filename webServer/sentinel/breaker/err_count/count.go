package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
	"github.com/alibaba/sentinel-golang/core/config"
	"github.com/alibaba/sentinel-golang/logging"
	"github.com/alibaba/sentinel-golang/util"
)

type StateChangeListener struct{}

func (s StateChangeListener) OnTransformToClosed(prev circuitbreaker.State, rule circuitbreaker.Rule) {
	fmt.Printf("rule.strategy: %+v, From %s to Closed, time: %d\n", rule.Strategy, prev.String(), util.CurrentTimeMillis())
}

func (s StateChangeListener) OnTransformToOpen(prev circuitbreaker.State, rule circuitbreaker.Rule, snapshot interface{}) {
	fmt.Printf("rule.strategy: %+v, From %s to Open, time: %d\n", rule.Strategy, prev.String(), util.CurrentTimeMillis())
}

func (s StateChangeListener) OnTransformToHalfOpen(prev circuitbreaker.State, rule circuitbreaker.Rule) {
	fmt.Printf("rule.strategy: %+v, From %s to HalOpen, time: %d\n", rule.Strategy, prev.String(), util.CurrentTimeMillis())
}

var resName = "abc"

func main() {
	cfg := config.NewDefaultConfig()
	// for testing, logging output to console
	// cfg.Sentinel.Log.Logger = logging.NewConsoleLogger() // 打印日志
	if err := sentinel.InitWithConfig(cfg); err != nil {
		log.Fatalf("%+v", err)
	}

	circuitbreaker.RegisterStateChangeListeners(&StateChangeListener{})

	// 异常数大于设定的阈值 触发熔断
	_, err := circuitbreaker.LoadRules([]*circuitbreaker.Rule{
		{
			Resource:                     resName,
			Strategy:                     circuitbreaker.ErrorCount,
			RetryTimeoutMs:               2000, //即熔断触发后持续的时间（单位为 ms）HalfOpen 时间
			MinRequestAmount:             5,    //最小静默数
			Threshold:                    10,   // 阈值
			StatIntervalMs:               3000, //  统计的时间窗口长度（单位为 ms）
			StatSlidingWindowBucketCount: 1,    //
		},
	})
	if err != nil {
		log.Fatalf("%v", err)
	}

	logging.Info("[CircuitBreaker ErrorCount] Sentinel to circuit breaker demo is running. You may see the passsed/closed metric in the metric log.")

	var ch = make(chan int)
	go func() {
		for i := 0; ; i++ {
			if e, b := sentinel.Entry(resName); b != nil {
				// blocked
				time.Sleep(time.Duration(rand.Uint64()%20) * time.Millisecond)
				println("circuit breaker: blocked 1~20ms ")
			} else {
				// passed
				if rand.Uint64()%20 > 9 {
					fmt.Printf("随机创造错误: 错误数 大于9...")
					sentinel.TraceError(e, errors.New("biz error"))
				}
				n := time.Duration(rand.Uint64()%80 + 10)
				time.Sleep(n * time.Millisecond)
				fmt.Println("passed - ", n*time.Millisecond)
				e.Exit()
			}
		}
	}()
	go func() {
		for {
			e, b := sentinel.Entry(resName)
			if b != nil {
				println("blocked")
				time.Sleep(time.Duration(rand.Uint64()%20) * time.Millisecond)
			} else {
				println("passed")
				time.Sleep(time.Duration(rand.Uint64()%80) * time.Millisecond)
				e.Exit()
			}
		}
	}()
	go func() {
		println("debug start...")
		time.Sleep(15 * time.Second)
		println("debug end...")
		ch <- 1
	}()

	<-ch
}
