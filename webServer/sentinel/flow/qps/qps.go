package main

import (
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/core/config"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/alibaba/sentinel-golang/logging"
)

const resName = "example-flow-qps-resource"

func main() {
	// We should initialize Sentinel first.
	conf := config.NewDefaultConfig()
	// for testing, logging output to console
	conf.Sentinel.Log.Logger = logging.NewConsoleLogger()

	// 初始化
	err := sentinel.InitWithConfig(conf)
	if err != nil {
		log.Fatalf("Unexpected error: %+v", err)
	}

	// 设置限流规则
	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               resName,
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              10,
			StatIntervalInMs:       1000,
		}, {},
	})
	if err != nil {
		log.Fatalf("加载配置异常 Unexpecked %+v", err)
	}

	//规则应用
	total := 0
	timeout := 0
	for i := 0; i < 5; i++ {
		go func(i int) {
			for {
				e, b := sentinel.Entry(resName, sentinel.WithTrafficType(base.Inbound))
				n := time.Duration(rand.Uint64() % 10 * 10)
				time.Sleep(n * time.Millisecond)
				if i == 1 {
					timeout += int(n)
				}
				if b != nil {
					println("阻塞住了 ", i, strings.Repeat("- ", i), n)
				} else {
					println("通了 ", i, strings.Repeat("- ", i), n)
					total++

					e.Exit()
				}
			}
		}(i)
	}

	var WG = sync.WaitGroup{}
	WG.Add(1)
	go func(WG *sync.WaitGroup) {
		// 1 秒后热更新配置 通过80个
		time.Sleep(1 * time.Second)
		println("耗时:", timeout, " 通过数:", total)
		total = 0
		_, err := flow.LoadRules([]*flow.Rule{
			{
				Resource:               resName,
				ControlBehavior:        flow.Reject,
				TokenCalculateStrategy: flow.Direct,
				Threshold:              80,
				StatIntervalInMs:       1000,
			},
		})
		if err != nil {
			log.Fatalf("加载新配置失败 %+v", err)
		} else {
			println("更新配置成功...")
		}
		time.Sleep(1 * time.Second)
		println("耗时:", timeout, " 通过数:", total)
		WG.Done()
	}(&WG)

	WG.Wait()
}
