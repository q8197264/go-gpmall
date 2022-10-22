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
	"github.com/alibaba/sentinel-golang/util"
	"github.com/nacos-group/nacos-sdk-go/common/logger"
)

type StateChangeListener struct{}

func (s StateChangeListener) OnTransformToClosed(prev circuitbreaker.State, rule circuitbreaker.Rule) {
	fmt.Printf("[ErrorRatio] rule.Strategy: %+v, From %s to Closed, time:%d\n", rule.Strategy, prev.String(), util.CurrentTimeMillis())
}

func (s StateChangeListener) OnTransformToOpen(prev circuitbreaker.State, rule circuitbreaker.Rule, snapshot interface{}) {
	fmt.Printf("[ErrorRatio] rule.Strategy: %+v, From %s to Open, time:%d\n", rule.Strategy, prev.String(), util.CurrentTimeMillis())
}

func (s StateChangeListener) OnTransformToHalfOpen(prev circuitbreaker.State, rule circuitbreaker.Rule) {
	fmt.Printf("[ErrorRatio] rule.Strategy: %+v, From %s to HalfOpen, time:%d\n", rule.Strategy, prev.String(), util.CurrentTimeMillis())
}

var resName = "abcd"

func main() {
	cfg := config.NewDefaultConfig()
	// cfg.Sentinel.Log.Logger = logging.NewConsoleLogger()
	if err := sentinel.InitWithConfig(cfg); err != nil {
		log.Fatalf("%+v", err)
	}

	circuitbreaker.RegisterStateChangeListeners(&StateChangeListener{})

	if ok, err := circuitbreaker.LoadRules([]*circuitbreaker.Rule{
		// Statistic time span=3s, recoveryTimeout=2s, maxErrorRatio=40%
		{
			Resource:                     resName,
			Strategy:                     circuitbreaker.ErrorRatio,
			Threshold:                    0.4,
			RetryTimeoutMs:               2000, //halfOpen time
			MinRequestAmount:             10,   //静默数
			StatIntervalMs:               3000, //计数周期
			StatSlidingWindowBucketCount: 1,    //越大越准[能被 StatIntervalMs 整除] bucket
		},
	}); !ok {
		log.Fatalf("load rules err: %+v", err)
	}

	logger.Info("hello world...")
	go func() {
		for {
			e, b := sentinel.Entry(resName)
			if b != nil {
				// blocked
				n := time.Duration(rand.Uint64()%20) * time.Millisecond
				fmt.Println("blocked...g1 - ", n)
				time.Sleep(n)
			} else {
				// passed
				m := rand.Uint64() % 20
				if m > 6 {
					// Record current invocation as error.
					fmt.Printf("trace error count.... %d ", m)
					sentinel.TraceError(e, errors.New("biz error"))
				}
				n := time.Duration(rand.Uint64()%80+20) * time.Millisecond
				fmt.Println("passed...g1 - ", n)
				time.Sleep(n)

				e.Exit()
			}
		}
	}()

	// go func() {
	// 	for {
	// 		if e, b := sentinel.Entry(resName); b != nil {
	// 			// blocked
	// 			time.Sleep(time.Duration(rand.Uint64()%20) * time.Millisecond)
	// 			println("blocked...g2")
	// 		} else {
	// 			// passed
	// 			time.Sleep(time.Duration(rand.Uint64()%80+40) * time.Millisecond)
	// 			println("passed..g2")
	// 			e.Exit()
	// 		}
	// 	}
	// }()

	ch := make(chan int)
	go func() {
		println("Error ratio begin...")
		time.Sleep(5 * time.Second)
		ch <- 1
		println("Error ratio end...")
	}()
	<-ch
}
