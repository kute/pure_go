package main

import (
	"context"
	"fmt"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8" // redis6用 v8，redis7用 v9
	"github.com/golang/groupcache/consistenthash"
	"hash/crc32"
	"time"
)

func main() {

	var ctx = context.Background()

	// 另一种通过解析 url 初始化 client 的方式，opt 通过查看源码可以获取到
	//opt, err := redis.ParseURL("redis://:kuteredis@localhost:6379/0?max_retries=3")
	//if err != nil {
	//	panic(err)
	//}
	//rdb := redis.NewClient(opt)

	rdb := redis.NewClient(&redis.Options{
		Network: "tcp",            // tcp or unix,default tcp
		Addr:    "localhost:6379", // host:port
		Dialer:  nil,
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			fmt.Println("redis 连接已建立", cn.String())
			return nil
		},
		Password:   "kuteredis",
		DB:         0,
		MaxRetries: 3,
		//MinRetryBackoff: 0,
		//MaxRetryBackoff: 0,
		//DialTimeout:        0,
		ReadTimeout:  time.Second * 3,
		WriteTimeout: time.Second,
		PoolSize:     10, // default 10
		MinIdleConns: 2,
		//MaxConnAge:         0,
		PoolTimeout: 0, // 从池中获取连接的超时时间（可能连接不足，或者 全部阻塞）
		//IdleTimeout:        0,
		//IdleCheckFrequency: 0,
		//TLSConfig:          nil,
	})

	defer CloseRedis(rdb)

	var key = "test:kute:110"
	var value = "kute"

	// 执行 set
	var statusCmd = rdb.Set(ctx, key, value, time.Minute*10)
	fmt.Println("当前执行的命令的名称是：", statusCmd.Name())
	fmt.Println("当前命令的参数：", statusCmd.Args())
	if err := statusCmd.Err(); err != nil { // 执行是否出错
		panic(err)
	}
	defer rdb.Del(ctx, key)

	// 执行 get
	var v, err = rdb.Get(ctx, key).Result()
	if err != nil {
		panic(err)
	} else if err == redis.Nil { // redis.Nil属于特殊的 nil 值，表示 redis-server 返回的 nil
		fmt.Println("当前 key 不存在")
	} else {
		fmt.Println("key 对应的值为：", v)
	}

	// 执行任意命令，包括 lua
	if v, err := rdb.Do(ctx, "get", key).Result(); err == nil {
		fmt.Println("key 对应的值为：", v)
	}

	// 执行 pipeline，执行多条命令
	pipeline := rdb.Pipeline()
	pipeline.Set(ctx, key, value, 0)
	pipeline.Get(ctx, key)
	if cmds, err := pipeline.Exec(ctx); err == nil {
		for _, v := range cmds {
			fmt.Println(v.String())
		}
	}

	// 函数形式 pipeline
	cmds, err := rdb.Pipelined(ctx, func(pip redis.Pipeliner) error {
		pip.Set(ctx, key, value, 0)
		pip.Get(ctx, key)
		return nil
	})
	if err == nil {
		for _, v := range cmds {
			fmt.Println(v.String())
		}
	}

	// 执行lua
	var incrBy = redis.NewScript(`
local key = KEYS[1]
local change = ARGV[1]

local value = redis.call("GET", key)
if not value then
  value = 0
end

value = value + change
redis.call("SET", key, value)

return value
`)
	var keys = []string{key}
	var values = []interface{}{+1}
	incrBy.Run(context.Background(), rdb, keys, values...)

	// 测试 cache, go-redis 中的 cache 模块，使用 redis 实现的缓存系统
	TestCache(rdb)

}

/**
测试缓存系统
*/
func TestCache(rdb *redis.Client) {
	var c = cache.New(&cache.Options{
		Redis:        rdb,
		LocalCache:   nil,
		StatsEnabled: false,
		Marshal:      nil,
		Unmarshal:    nil,
	})
}

// 构建线程安全的 redis-client，即 redis ring，使用一致性哈希将 key 分布在所有节点上
func TestRedisRing() {
	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			// 指定节点名称
			"node-1": "localhost:70000",
			"node-2": "localhost:70000",
		},
		Password: "xxxx",
		DB:       0,
		NewConsistentHash: func(shards []string) redis.ConsistentHash {
			// 返回一致性哈希，默认使用  Rendezvous 哈希算法
			return consistenthash.New(len(shards), crc32.ChecksumIEEE)
		},
	})
	defer ring.Close()

}

// redis 集群
func TestRedisCluster() {
	var ctx = context.Background()

	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{"localhost:70000", "localhost:70001", "localhost:70002", "localhost:70003", "localhost:70004", "localhost:70005"},
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			return nil
		},
		Password: "xxxx",
	})
	defer rdb.Close()
	// 遍历 master 节点
	var err = rdb.ForEachMaster(ctx, func(ctx context.Context, client *redis.Client) error {
		fmt.Println("当前 master 是：", client.ClientGetName(ctx))
		return nil
	})

	// 遍历 slave 节点
	err = rdb.ForEachSlave(ctx, func(ctx context.Context, client *redis.Client) error {
		fmt.Println("当前 slave 是：", client.ClientGetName(ctx))
		return nil
	})

	// 遍历所有节点
	err = rdb.ForEachShard(ctx, func(ctx context.Context, client *redis.Client) error {
		fmt.Println("当前 节点 是：", client.ClientGetName(ctx))
		return nil
	})
	if err != nil {
		panic(err)
	}

}

// RedisUniversalClient 是一个通用型 client，可以表示 cluster, sentinel, single 类型，具体实际的类型取决于 传参选项
// 当 MasterName 不为空时，表示 是 基于 sentinel 的FailoverClient
// 当 Addrs 中节点个数 >= 2时，表示是 cluster
// 否则就是 single client
func TestUniversalClient() {
	client := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:      nil,
		MasterName: "",
	})
	defer client.Close()
}

func CloseRedis(client *redis.Client) {
	if err := client.Close(); err != nil {
		panic(err)
	} else {
		fmt.Println("redis closed")
	}
}
