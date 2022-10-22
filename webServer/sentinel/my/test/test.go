package main

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

// 令牌桶限流
/*
	threshold / capacity 阈值(容量)
	rate 出水速率  x/s
	leftwater 剩余量
	lastTime

		permits
			|
			|
		|-------|
		|		|
		|~~~~~~~|
		|_______|
			|
	------->|-------->pass
			|
			reject
需结合队列
*/

// redis connect ...
var ctx = context.Background()
var rdb *redis.Client

func RedisConnect() (client *redis.Client, c context.Context) {
	if rdb == nil {
		rdb = redis.NewClient(&redis.Options{
			Addr:     ":6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		})
		// 心跳
		_, err := rdb.Ping(ctx).Result()
		if err != nil {
			println(err)
		}
	}

	return rdb, ctx
}

func main() {
	opts := "redis"
	ratelimiter := NewRatelimiter()
	ch := make(chan int)
	for i := 0; i < 10; i++ {
		go func() {
			for i := 0; i < 25; i++ {
				switch opts {
				case "local":
					if ratelimiter.TokenBucketRateLimiter() {
						fmt.Printf(" .")
					} else {
						fmt.Printf(" x")
					}
				case "redis":
					if ratelimiter.TokenBucketRateLimiterRedis("currenttoken") {
						fmt.Printf(" .")
					} else {
						fmt.Printf(" x")
					}
				}
				time.Sleep(time.Duration(rand.Uint64()%100) * time.Millisecond)
			}
			ch <- 1
		}()
	}
	start := time.Now().UnixNano() / 1e6
	for i := 0; i < 5; i++ {
		<-ch
	}
	end := time.Now().UnixNano() / 1e6

	fmt.Printf("\ntime: %v ms\n", end-start)
}

/**
  令牌桶限流
**/

// 令牌发放速率
var once sync.Once

type ratelimiter struct {
	rdb               *redis.Client
	ctx               context.Context
	lastOutTime       int64
	generateTokenRate int64
	currentTokens     float64
	capacity          int64
	luaHash           string
}

func NewRatelimiter() *ratelimiter {
	rdb, ctx := RedisConnect()
	return &ratelimiter{
		rdb:               rdb,
		ctx:               ctx,
		generateTokenRate: 15,
		currentTokens:     0,
		capacity:          25,
	}
}

func (r *ratelimiter) TokenBucketRateLimiter() bool {
	if r.currentTokens == 0 {
		// 初始请求时间
		once.Do(func() {
			r.lastOutTime = time.Now().UnixNano() / 1e6
		})
	}

	time.Sleep(50 * time.Millisecond)

	now := time.Now().UnixNano() / 1e6
	// 生成令牌数
	//每秒发放令牌速率, 不足1s则不发放. 大于1s时间差
	if now-r.lastOutTime >= 1000 {
		generateTokens := (time.Now().UnixNano()/1e6 - r.lastOutTime) / 1e3 * r.generateTokenRate
		r.currentTokens = math.Min(r.currentTokens+float64(generateTokens), float64(r.capacity))
		fmt.Printf(" - current token total:%v\n", r.currentTokens)
		r.lastOutTime = now
	}

	// 令牌数量大于0
	if r.currentTokens >= 1 {
		r.currentTokens--
		return true
	}

	return false
}

func (r *ratelimiter) TokenBucketRateLimiterRedis(key string) bool {
	if r.currentTokens == 0 {
		once.Do(func() {
			r.lastOutTime = time.Now().UnixNano() / 1e6
			r.rdb.Del(r.ctx, key)

			script := `
				-- 当前总令牌数
				local currenttoken = redis.call("Get", KEYS[1])
				if currenttoken == false then
					-- 新增令牌数
					redis.call("Set", KEYS[1], ARGV[1])
					return ARGV[1]
				else
					-- 当前总令牌数
					redis.call("Set", KEYS[1], math.min(currenttoken+ARGV[1], ARGV[2]))
					return math.min(currenttoken+ARGV[1], ARGV[2])
				end
			`
			// rdb, ctx := RedisConnect()
			r.luaHash, _ = r.rdb.ScriptLoad(r.ctx, script).Result()
			// if err != nil {
			// 	fmt.Println(err.Error())
			// }
		})
	}
	// 生成令牌
	now := time.Now().UnixNano() / 1e6
	if now-r.lastOutTime >= 1000 {
		generateTokens := (now - r.lastOutTime) / 1e3 * r.generateTokenRate
		v, err := rdb.EvalSha(ctx, r.luaHash, []string{key}, generateTokens, r.capacity).Result()
		if err != nil {
			fmt.Println(err.Error())
			return false
		}
		fmt.Printf(" -- %v\n", v)
		r.lastOutTime = now
	}

	// 消耗令牌
	decrScript := redis.NewScript(`
		local v = redis.call("Get", KEYS[1])
		if v ~= false and (v*100) >= 100 then
			redis.call("Decr", KEYS[1])
			return 1
		end
		return 0
	`)
	m, err := decrScript.Run(ctx, rdb, []string{key}).Result()
	if err != nil {
		fmt.Println(err.Error())
	}
	if n, ok := m.(int64); ok {
		if n >= 1 {
			return true
		}
	}

	return false
}
