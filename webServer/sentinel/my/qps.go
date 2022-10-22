package main

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/go-redis/redis/v8"
)

var (
	lastTime = time.Now().UnixNano() / 1e6
	// StatIntervalMs = uint64(500) //时间窗口总大小（毫秒）
	// threshold      = uint64(10)  //阈值
	counter = uint64(0)
)

func main() {
	// fmt.Println(time.Unix(currentTime, 10).Format("2006-01-02 15:04:05"))

	// 固定窗
	// FixedWindowRateLimiterTest()

	// 滑动窗
	SlideWindowRateLimiterTest()

	// 漏桶
	// LeakBucketRateLimiterTest()

	// 令牌桶
	// TokenBucketRateLimiterTest()

	// Test()
}

func Test() {
	var x int32
	// threads := runtime.GOMAXPROCS(0)
	var j int32
	for i := 0; i < 2; i++ {
		go func() {
			for {
				// atomic.AddInt32(&x, 1)
				x++
				j = 2
			}
		}()
	}
	time.Sleep(time.Second)
	fmt.Println("x =", atomic.LoadInt32(&x), "j=", j)
}

// 测试用例
func FixedWindowRateLimiterTest() {
	var wg sync.WaitGroup
	wg.Add(2)
	for i := 0; i < 2; i++ {
		go func(i int) {
			defer wg.Done()
			if i%2 == 0 {
				time.Sleep(time.Duration(300 * time.Millisecond))
			} else {
				time.Sleep(time.Duration(500 * time.Millisecond))
			}
			fmt.Printf("|")
			for i := 0; i < 25; i++ {
				ok := FixedWindowRateLimiter(500, 10)
				if ok {
					fmt.Printf(".")
				} else {
					fmt.Printf("x")
				}
				// ms := rand.Uint64()%10 + 20
				time.Sleep(time.Duration(20) * time.Millisecond)
			}
		}(i)

	}
	start := time.Now().UnixNano() / 1e6
	wg.Wait()

	fmt.Printf("\n%d ms\n", time.Now().UnixNano()/1e6-start)
}

func SlideWindowRateLimiterTest() {
	var storage = "redis"
	var ms uint64
	for y := 0; y < 2; y++ {
		for i := 0; i < 150; i++ {
			switch storage {
			case "redis":
				if SlideWindowRateLimiterWithRedis("sorted-set", 10, 1000) {
					fmt.Printf(".")
				} else {
					fmt.Printf("x")
				}
			default:
				if SlideWindowRateLimiter("sorted-test", 10, 1000) {
					fmt.Printf(".")
				} else {
					fmt.Printf("x")
				}
			}

			ms = rand.Uint64() % 100
			// counter += ms
			// if counter%1000 > 900 {
			// 	println("--", counter)
			// }
			time.Sleep(time.Duration(ms) * time.Millisecond)
		}
		fmt.Printf("\n")
	}

	fmt.Printf("\n")
}

// 最多每秒只能过多少(无损),超时不候
func LeakBucketRateLimiterTest() {
	sum := make(chan string, 1)
	var counter int64 = 0
	LeakBucketInit()
	for i := 0; i < 10; i++ {
		go func(i int) {
			var sumCut uint64 = 0
			var countCut uint64 = 0
			for k := 0; k < 50; k++ {
				if LeakBucketRateLimiter(k) {
					// queue producer
					client, ctx := RedisConnect()
					client.LPush(ctx, "list-test", k)
					fmt.Printf(". ")
					atomic.AddInt64(&counter, 1)
					countCut++
				} else {
					fmt.Printf("x ")
				}
				ms := rand.Uint64()%50 + 10
				// ms := uint64(10)
				sumCut += ms
				time.Sleep(time.Duration(ms) * time.Millisecond)
			}
			sum <- fmt.Sprintf("%d:%dms-%d", i, sumCut, countCut)
		}(i)
	}

	for i := 0; i < 10; i++ {
		go LeakBucketConsumer()
	}
	startTime := time.Now().UnixNano() / 1e6
	for i := 0; i < 10; i++ {
		// fmt.Println(<-sum)
		<-sum
	}
	fmt.Printf("\n>> counter:%v | time:%v ms\n", atomic.LoadInt64(&counter), (time.Now().UnixNano()/1e6 - startTime))
}

func TokenBucketRateLimiterTest() {
	opts := "redis"
	var luaHash string
	var err error
	if opts == "redis" {
		luaHash, err = initTokenBucket()
		if err != nil {
			println(err.Error())
		}
	}
	ch := make(chan int)
	for i := 0; i < 5; i++ {
		go func() {
			for i := 0; i < 50; i++ {
				switch opts {
				case "local":
					if TokenBucketRateLimiter() {
						fmt.Printf(" .")
					} else {
						fmt.Printf(" x")
					}
				case "redis":
					if TokenBucketRateLimiterRedis(luaHash) {
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

// 固定窗口
// 10/500ms
// 出现问题: 时间边界,  瞬时流量2n
func FixedWindowRateLimiter(statIntervalMs int64, threshold uint64) bool {
	var currentTime = time.Now().UnixNano() / 1e6
	if currentTime-lastTime < int64(statIntervalMs) {
		if atomic.LoadUint64(&counter) < threshold {
			atomic.AddUint64(&counter, 1)
			return true
		}
		return false
	}
	lastTime = currentTime
	counter = uint64(0)

	return true

}

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

// 滑动窗口限流
/*
	每分钟限制请求数
	计数器, k-为当前窗口的开始时间值秒，value为当前窗口的计数
*/
// storage local map
var windowList map[string][]int64

func SlideWindowRateLimiter(key string, threshold int, StatIntervalMs int64) bool {
	// 阈值 threshold = 200/s
	// 时间周期 statIntervalMs 1000ms

	var currentTime = time.Now().UnixNano() / 1e6 //ms
	if windowList == nil {
		windowList = make(map[string][]int64)
	}

	if _, ok := windowList[key]; !ok {
		windowList[key] = make([]int64, 0)
	}

	if len(windowList[key]) < threshold {
		// fmt.Printf(".")
		windowList[key] = append(windowList[key], currentTime)
		return true
	}

	// 队列满了, 取出最早访问的时间
	if currentTime-windowList[key][0] <= int64(StatIntervalMs) {
		// 时间周期内, 数量已满, 拒绝 reject
		// fmt.Printf("x")
		return false
	} else {
		// 时间周期外,更新
		windowList[key] = windowList[key][1:]
		windowList[key] = append(windowList[key], currentTime)
		// fmt.Printf("*")
		return true
	}
}

// storage redis
func SlideWindowRateLimiterWithRedis(key string, threshold uint64, statIntervalMs int64) bool {
	var currentTime = time.Now().UnixNano() / 1e6

	// 数量不足, 通过
	client, ctx := RedisConnect()
	count, err := client.ZCard(ctx, key).Result()
	if err != nil {
		println(err)
	}
	if uint64(count) < threshold {
		client.ZAdd(ctx, key, &redis.Z{
			Score:  float64(currentTime),
			Member: currentTime,
		}).Result()

		return true
	}

	// 时间周期内, 拒绝
	res, err := client.ZRangeByScoreWithScores(ctx, key, &redis.ZRangeBy{
		Min:    "-inf",
		Max:    "+inf",
		Offset: 0,
		Count:  2,
	}).Result()
	if err != nil {
		println(err)
	}
	if currentTime-int64(res[0].Score) <= statIntervalMs {
		return false
	} else {
		// 删除所有周期外过期成员
		if _, err = client.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", currentTime-statIntervalMs)).Result(); err != nil {
			println(err)
		}

		client.ZAdd(ctx, key, &redis.Z{
			Score:  float64(currentTime),
			Member: currentTime,
		}).Result()

		return true
	}
}

// 漏桶限流
/*
	threshold / capacity 阈值(容量)
	rate 出水速率  x/s
	leftwater 剩余量
	lastTime
	|
	|
|-------|
|		|
|~~~~~~~|
|_______|
	|
	|
需结合队列
*/
var capacity float64 = 25
var leakRate int64 = 25 // 流出速率越大, 漏水量越大
var lastOutTime int64
var leftWater float64 = 0 // 剩余水量

var once sync.Once

var minLuaHash string
var compareLuaHash string

func LeakBucketInit() {
	rdb, ctx := RedisConnect()
	var script string = `
		local value = redis.call("Get", KEYS[1])
		if (value*100) <= (ARGV[1]*100) then
			redis.call("Set", KEYS[1], 0)
			return 0
		else
			redis.call("Set", KEYS[1], string.format("%.3f", value-ARGV[1]))
			return 1
		end
	`
	minLuaHash, _ = rdb.ScriptLoad(ctx, script).Result()
	script = `
		local v = redis.call("Get", KEYS[1])
		if v==false or (v*100) <= (ARGV[1]*100) then
			redis.call("Incrbyfloat", KEYS[1], 1)
			return 1
		end
		return 0
	`
	compareLuaHash, _ = rdb.ScriptLoad(ctx, script).Result()
}

func LeakBucketRateLimiter(k int) bool {
	rdb, ctx := RedisConnect()

	// 初始请求时间
	if leftWater == 0 {
		once.Do(func() {
			rdb.Del(ctx, "leftwater")
		})

		resetincr := redis.NewScript(`
			local val = redis.call("Get", KEYS[1])
			if (val == false or val*100 < ARGV[1]*100) then
				return redis.call("Incrbyfloat", KEYS[1], 1)
			end
			return "0"
		`)
		ss, err := resetincr.Run(ctx, rdb, []string{"leftwater"}, capacity).Result()
		if err != nil {
			println("once err <<", err.Error())
		}

		lastOutTime = time.Now().UnixNano() / 1e6 //ms
		if n, ok := ss.(string); ok {
			leftWater, _ = strconv.ParseFloat(n, 64)
			if leftWater >= float64(1) {
				return true
			}
		}

		return false
	}

	// 漏水量 = (当前时间 - 最近一次时间)*漏水速度
	waterLeaked := float64((time.Now().UnixNano()/1e6-lastOutTime)*leakRate) / 1e3

	// 剩余水量
	_, err := rdb.EvalSha(ctx, minLuaHash, []string{"leftwater"}, waterLeaked).Result()
	if err != nil {
		fmt.Println("lua err: ", err.Error())
		return false
	}

	// fmt.Printf("漏水:%.2f, 剩余水量:%.2f\n", waterLeaked, leftWater)
	// fmt.Printf("时间差:%.2fs, 桶内水占比: %.2f%%\n", float64(time.Now().UnixNano()/1e6-lastOutTime)/1e3, leftWater/capacity*100)
	lastOutTime = time.Now().UnixNano() / 1e6

	b, errs := rdb.EvalSha(ctx, compareLuaHash, []string{"leftwater"}, capacity).Result()
	if errs != nil {
		fmt.Println("==> ", errs)
		return false
	}
	if n, ok := b.(int64); ok {
		if n == 1 {
			return true
		}
	}

	return false
}

func AddFloat32(val *float32, delta float32) (new float32) {
	for {
		old := *val
		new = old + delta
		if atomic.CompareAndSwapUint32(
			(*uint32)(unsafe.Pointer(val)),
			math.Float32bits(old),
			math.Float32bits(new),
		) {
			break
		}
	}
	return
}

// 恒定速率处理请求
func LeakBucketConsumer() {
	client, ctx := RedisConnect()
	i := 0
	for {
		if v, err := client.RPop(ctx, "list-test").Result(); err == nil {
			i++
			println("loading consumer[", goID(), "]: value=", v, "counter:", i)
		}
		time.Sleep(time.Millisecond * 10)
	}
}

// 获取协程id
func goID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

/**
  令牌桶限流
**/
// 令牌发放速率
var generateTokenRate int64 = 15
var currentTokens float64

func TokenBucketRateLimiter() bool {
	if currentTokens == 0 {
		// 初始请求时间
		once.Do(func() {
			lastOutTime = time.Now().UnixNano() / 1e6
		})
	}

	time.Sleep(50 * time.Millisecond)

	now := time.Now().UnixNano() / 1e6
	// 生成令牌数
	//每秒发放令牌速率, 不足1s则不发放. 大于1s时间差
	if now-lastOutTime >= 1000 {
		generateTokens := (time.Now().UnixNano()/1e6 - lastOutTime) / 1e3 * generateTokenRate
		currentTokens = math.Min(currentTokens+float64(generateTokens), float64(capacity))
		fmt.Printf(" - current token total:%v\n", currentTokens)
		lastOutTime = now
	}

	// 令牌数量大于0
	if currentTokens >= 1 {
		currentTokens--
		return true
	}

	return false
}

func initTokenBucket() (string, error) {
	script := `
		local currenttoken = redis.call("Get", KEYS[1])
		if currenttoken == false then
			redis.call("Set", KEYS[1], ARGV[1])
			return ARGV[1]
		else
			redis.call("Set", KEYS[1], math.min(currenttoken+ARGV[1], ARGV[2]))
			return math.min(currenttoken+ARGV[1], ARGV[2])
		end
	`
	rdb, ctx := RedisConnect()
	luaHash, err := rdb.ScriptLoad(ctx, script).Result()

	return luaHash, err
}

func TokenBucketRateLimiterRedis(luaHash string) bool {
	rdb, ctx := RedisConnect()

	if currentTokens == 0 {
		once.Do(func() {
			lastOutTime = time.Now().UnixNano() / 1e6
			rdb.Del(ctx, "currenttoken")
		})
	}
	// 生成令牌
	now := time.Now().UnixNano() / 1e6
	if now-lastOutTime >= 1000 {
		generateTokens := (now - lastOutTime) / 1e3 * generateTokenRate
		v, err := rdb.EvalSha(ctx, luaHash, []string{"currenttoken"}, generateTokens, capacity).Result()
		if err != nil {
			fmt.Println(err.Error())
			return false
		}
		fmt.Printf(" -- %v\n", v)
		lastOutTime = now
	}

	decrScript := redis.NewScript(`
		local v = redis.call("Get", KEYS[1])
		if v ~= false and (v*100) >= 100 then
			redis.call("Decr", KEYS[1])
			return 1
		end
		return 0
	`)
	m, err := decrScript.Run(ctx, rdb, []string{"currenttoken"}).Result()
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
