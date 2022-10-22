package main

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/core/flow"
)

var resName = "example-flow-warmup-resource"

func main() {
	if err := sentinel.InitDefault(); err != nil {
		log.Fatalf("%+v", err)
	}

	if ok, err := flow.LoadRules([]*flow.Rule{
		{
			Resource:               resName,
			TokenCalculateStrategy: flow.WarmUp,
			ControlBehavior:        flow.Reject,
			Threshold:              100,
			WarmUpPeriodSec:        10, //预热时间长度
			WarmUpColdFactor:       3,  // 预热的因子，默认是3，该值的设置会影响预热的速度
			StatIntervalInMs:       1000,
		},
	}); err != nil {
		println(ok, err.Error())
	} else {
		println(ok, "配置不一样")
	}
	println("开始:", time.Now().String())
	total := 0
	limit := 0
	for i := 0; i < 10; i++ {
		go func(i int) {
			for {
				e, b := sentinel.Entry(resName, sentinel.WithTrafficType(base.Inbound))
				time.Sleep(time.Duration(rand.Uint64()%10*10) * time.Millisecond)
				if b != nil {
					fmt.Printf("%s", strings.Repeat(".", 1))
					total++
				} else {
					// println("限制了")
					fmt.Printf("%s", strings.Repeat("-", 1))
					limit++
					e.Exit()
				}
			}
		}(i)
	}

	// var wg sync.WaitGroup
	// wg.Add(1)
	// go func(wg *sync.WaitGroup) {
	// 	time.Sleep(1 * time.Second)
	// 	println("通过:", total, " 限制:", limit)
	// 	wg.Done()
	// }(&wg)

	// wg.Wait()

	time.Sleep(1000 * time.Millisecond)
	println("耗时:", time.Now().String(), "通过:", total, " 限制:", limit)
}
