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
	"github.com/shopspring/decimal"
)

type StateChangeListener struct{}

func (s StateChangeListener) OnTransformToClosed(prev circuitbreaker.State, rule circuitbreaker.Rule) {
	fmt.Printf("[Strategy.SlowRatio] rule.Strategy:%+v, From %s to Closed, time:%d\n", rule.Strategy, prev.String(), util.CurrentTimeMillis())
}

func (s StateChangeListener) OnTransformToOpen(prev circuitbreaker.State, rule circuitbreaker.Rule, snapshot interface{}) {
	fmt.Printf("[Strategy.SlowRatio] rule.Strategy:%+v, From %s to Open, time:%d\n", rule.Strategy, prev.String(), util.CurrentTimeMillis())
}

func (s StateChangeListener) OnTransformToHalfOpen(prev circuitbreaker.State, rule circuitbreaker.Rule) {
	fmt.Printf("[Strategy.SlowRatio] rule.Strategy:%+v, From %s to HalfOpen, time:%d\n", rule.Strategy, prev.String(), util.CurrentTimeMillis())
}

var resName = "abcde"

func main() {
	cfg := config.NewDefaultConfig()
	// cfg.Sentinel.Log.Logger = logging.NewConsoleLogger()
	if err := sentinel.InitWithConfig(cfg); err != nil {
		log.Fatalf("%+v", err)
	}

	circuitbreaker.RegisterStateChangeListeners(&StateChangeListener{})

	logging.Info("[CircuitBreaker SlowRtRatio] Sentinel Go circuit breaking demo is running. You may see the pass/block metric in the metric log.")

	if ok, err := circuitbreaker.LoadRules([]*circuitbreaker.Rule{
		{
			Resource:                     resName,
			Strategy:                     circuitbreaker.SlowRequestRatio,
			Threshold:                    0.4,
			MinRequestAmount:             10,
			RetryTimeoutMs:               3000,
			MaxAllowedRtMs:               50,
			StatIntervalMs:               5000,
			StatSlidingWindowBucketCount: 10,
		},
	}); !ok {
		log.Fatalf("%+v", err)
	}
	// docker run --rm pantsel/konga:latest -c kong -a postgres -u postgresql://kong:ecs123@192.
	go func() {
		passedTotalTime := time.Duration(0)
		closedTotalTime := time.Duration(0)
		count := 0
		countAll := 0
		ratio := float64(0)
		for {
			if e, b := sentinel.Entry(resName); b != nil {
				// blocked
				passedTotalTime = time.Duration(0)
				countAll = 0
				count = 0
				ratio = float64(0)

				n := time.Duration(rand.Uint64()%20 + 20)
				closedTotalTime += n

				// fmt.Println(closedTotalTime, "==", time.Duration(2))
				if closedTotalTime > time.Duration(5)*time.Microsecond/2 {
					fmt.Println("blocked...", n*time.Millisecond, " - closedTotalTime:", closedTotalTime)
				}

				time.Sleep(n * time.Millisecond)
			} else {
				// passed
				closedTotalTime = time.Duration(0)
				if rand.Uint64()%20 > 9 {
					sentinel.TraceError(e, errors.New("biz error"))
					fmt.Printf("trace error...")
				}
				nn := time.Duration(rand.Uint64() % 80)

				if nn > time.Duration(50) {
					count++
				}
				if passedTotalTime > closedTotalTime {
					countAll++
					ratio, _ = decimal.NewFromFloat(float64(count) / float64(countAll)).Round(2).Float64()
				}

				passedTotalTime += nn

				fmt.Println("passed...", nn*time.Millisecond, " - total time:", passedTotalTime*time.Millisecond, " - Slow_RT_Ratio:", ratio*100, "%")
				time.Sleep(nn * time.Millisecond)

				e.Exit()
			}
		}
	}()

	go func() {
		for {
			if e, b := sentinel.Entry(resName); b != nil {
				println("blocked...")
				n := time.Duration(rand.Uint64()%20) * time.Millisecond
				time.Sleep(n)
			} else {
				if rand.Uint64()%20 > 9 {
					fmt.Printf("trace error ...")
				}
				nn := time.Duration(rand.Uint64()%80 + 10)
				fmt.Println("passed  - ", nn)
				time.Sleep(nn * time.Millisecond)

				e.Exit()
			}
		}
	}()

	ch := make(chan int)
	go func() {
		println("begin...")
		time.Sleep(15 * time.Second)
		ch <- 1
		println("end...")
	}()
	<-ch
}
