package initialize

import (
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/flow"
	"go.uber.org/zap"
)

func InitSentinel() {
	resName := "order-flow-warmup-resource"

	err := sentinel.InitDefault()
	if err != nil {
		zap.S().Fatalf("初始化sentinel 异常: %v\n", err)
	}

	// 这种配置应该从nacos读取
	if ok, err := flow.LoadRules([]*flow.Rule{
		{
			Resource:               resName,
			TokenCalculateStrategy: flow.WarmUp,
			ControlBehavior:        flow.Reject,
			Threshold:              5,
			WarmUpPeriodSec:        10, //预热时间长度
			WarmUpColdFactor:       3,  // 预热的因子，默认是3，该值的设置会影响预热的速度
			StatIntervalInMs:       1000,
		},
	}); !ok {
		zap.S().Fatalf("加载规则失败: %v", err)
	}

}

func breaking() {}
