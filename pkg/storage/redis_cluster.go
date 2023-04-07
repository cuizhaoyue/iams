package storage

import (
	"context"
	"crypto/tls"
	"github.com/cuizhaoyue/toolkit/log"
	"github.com/marmotedu/errors"
	"github.com/redis/go-redis/v9"
	uuid "github.com/satori/go.uuid"
	"strconv"
	"sync/atomic"
	"time"
)

// ErrRedisDown 在无法和redis进行通信时返回.
var ErrRedisDown = errors.New("storage: Redis is either down or not configured")

// Options 定义了redis集群选项.
type Options struct {
	Host                  string   // redis服务器ip地址，默认127.0.0.1
	Port                  int      // redis端口，默认6379
	Addrs                 []string // redis服务器地址,ip:port格式
	MasterName            string   // redis集群 master 名称(哨兵模式中master模式，只有创建failover客户端时使用)
	Username              string   // redis登录用户名
	Password              string   // redis密码
	Database              int      // redis数据库
	MaxIdle               int      // redis连接池中最大空闲连接数
	MaxActive             int      // 最大活跃连接数
	Timeout               int      // 连接redis的超时时间
	EnableCluster         bool     // 是否开启集群模式
	UseSSL                bool     // 是否启用TLS
	SSLInsecureSkipVerify bool     // 是否跳过安全验证，当连接redis时允许使用自签名证书
}

// NewRedisClusterPool 创建redis集群连接池.
func NewRedisClusterPool(opts *Options) redis.UniversalClient {
	log.Debug("Creating new Redis connection pool")

	// poolSize 适用于每个群集节点，而不是整个群集
	poolSize := 500
	if opts.MaxActive > 0 {
		poolSize = opts.MaxActive
	}

	timeout := time.Second * 5

	if opts.Timeout > 0 {
		timeout = time.Duration(opts.Timeout) * time.Second
	}

	var tlsConfig *tls.Config

	if opts.UseSSL {
		tlsConfig = &tls.Config{
			InsecureSkipVerify: opts.SSLInsecureSkipVerify,
		}
	}

	redisOpts := RedisOpts{
		Addrs:        getRedisAddrs(opts),
		MasterName:   opts.MasterName,
		Password:     opts.Password,
		DB:           opts.Database,
		DialTimeout:  timeout,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
		PoolSize:     poolSize,
		TLSConfig:    tlsConfig,
	}

	var client redis.UniversalClient
	if redisOpts.MasterName != "" {
		// 哨兵模式
		log.Info("---> [REDIS] Creating sentinel-backed failover client")
		client = redis.NewFailoverClient(redisOpts.failoverOpts())
	} else if opts.EnableCluster {
		// 集群模式
		log.Info("---> [REDIS] Creating cluster client")
		client = redis.NewClusterClient(redisOpts.clusterOpts())
	} else {
		// 单节点模式
		log.Info("---> [REDIS] Creating single-node client")
		client = redis.NewClient(redisOpts.simpleOpts())
	}

	return client
}

// 获取redis地址列表
func getRedisAddrs(opts *Options) (addrs []string) {
	// 如果已经配置了addrs则使用配置的地址
	if len(opts.Addrs) != 0 {
		addrs = opts.Addrs
	}

	// 未配置addrs且已配置端口，则向addrs中添加addr
	if len(addrs) == 0 && opts.Port != 0 {
		addr := opts.Host + ":" + strconv.Itoa(opts.Port)
		addrs = append(addrs, addr)
	}

	return addrs
}

type RedisOpts redis.UniversalOptions

// 创建集群模式客户端的选项
func (o *RedisOpts) clusterOpts() *redis.ClusterOptions {
	if len(o.Addrs) == 0 {
		o.Addrs = []string{"127.0.0.1:6379"}
	}

	return &redis.ClusterOptions{
		Addrs:     o.Addrs,
		OnConnect: o.OnConnect,

		Password: o.Password,

		MaxRedirects:   o.MaxRedirects,
		ReadOnly:       o.ReadOnly,
		RouteByLatency: o.RouteByLatency,
		RouteRandomly:  o.RouteRandomly,

		MaxRetries:      o.MaxRetries,
		MinRetryBackoff: o.MinRetryBackoff,
		MaxRetryBackoff: o.MaxRetryBackoff,

		DialTimeout:  o.DialTimeout,
		ReadTimeout:  o.ReadTimeout,
		WriteTimeout: o.WriteTimeout,
		PoolSize:     o.PoolSize,
		MinIdleConns: o.MinIdleConns,
		PoolTimeout:  o.PoolTimeout,

		// v9中没有这三个字段，暂注释
		// MaxConnAge:         o.MaxConnAge,
		// IdleTimeout:        o.IdleTimeout,
		// IdleCheckFrequency: o.IdleCheckFrequency,

		TLSConfig: o.TLSConfig,
	}
}

// 创建单实例客户端的选项
func (o *RedisOpts) simpleOpts() *redis.Options {
	addr := "127.0.0.1:6379"

	if len(o.Addrs) > 0 {
		addr = o.Addrs[0]
	}

	return &redis.Options{
		Addr:      addr,
		OnConnect: o.OnConnect,

		DB:       o.DB,
		Password: o.Password,

		MaxRetries:      o.MaxRetries,
		MinRetryBackoff: o.MinRetryBackoff,
		MaxRetryBackoff: o.MaxRetryBackoff,

		DialTimeout:  o.DialTimeout,
		ReadTimeout:  o.ReadTimeout,
		WriteTimeout: o.WriteTimeout,

		PoolSize:     o.PoolSize,
		MinIdleConns: o.MinIdleConns,
		PoolTimeout:  o.PoolTimeout,

		// MaxConnAge:         o.MaxConnAge,
		// IdleTimeout:        o.IdleTimeout,
		// IdleCheckFrequency: o.IdleCheckFrequency,

		TLSConfig: o.TLSConfig,
	}
}

// 创建哨兵模式客户端的选项
func (o *RedisOpts) failoverOpts() *redis.FailoverOptions {
	if len(o.Addrs) == 0 {
		o.Addrs = []string{"127.0.0.1:26379"}
	}

	return &redis.FailoverOptions{
		SentinelAddrs: o.Addrs,
		MasterName:    o.MasterName,
		OnConnect:     o.OnConnect,

		DB:       o.DB,
		Password: o.Password,

		MaxRetries:      o.MaxRetries,
		MinRetryBackoff: o.MinRetryBackoff,
		MaxRetryBackoff: o.MaxRetryBackoff,

		DialTimeout:  o.DialTimeout,
		ReadTimeout:  o.ReadTimeout,
		WriteTimeout: o.WriteTimeout,

		PoolSize:     o.PoolSize,
		MinIdleConns: o.MinIdleConns,
		PoolTimeout:  o.PoolTimeout,

		// MaxConnAge:         o.MaxConnAge,
		// IdleTimeout:        o.IdleTimeout,
		// IdleCheckFrequency: o.IdleCheckFrequency,

		TLSConfig: o.TLSConfig,
	}
}

var (
	singlePool      atomic.Value
	singleCachePool atomic.Value
	redisUp         atomic.Value

	// true: 禁用redis， false: 不禁用redis
	disableRedis atomic.Value
)

// DisableRedis very handy when testsing it allows to dynamically enable/disable talking with redisW.
func DisableRedis(ok bool) {
	if ok {
		disableRedis.Store(true)
		redisUp.Store(false)

		return
	}
	disableRedis.Store(false)
	redisUp.Store(true)
}

// 判断是否禁用redis.
func shouldConnect() bool {
	if v := disableRedis.Load(); v != nil {
		return !v.(bool)
	}

	return true
}

// 获取redis客户端连接池
func singleton(cache bool) redis.UniversalClient {
	// 如果cache为true，则先从缓存池中取客户端
	if cache {
		if v := singleCachePool.Load(); v != nil {
			return v.(redis.UniversalClient)
		}
		return nil
	}

	// 获取客户端
	if v := singlePool.Load(); v != nil {
		return v.(redis.UniversalClient)
	}

	return nil
}

// nolint: unparam
// 创建redis客户端连接池
func connectSingleton(cache bool, opts *Options) {
	if singleton(cache) == nil {
		log.Debug("Connecting to redis cluster")

		if cache {
			// 缓存redis客户端连接池
			singleCachePool.Store(NewRedisClusterPool(opts))
		}

		singlePool.Store(NewRedisClusterPool(opts))
	}
}

// RedisCluster 是一个使用redis数据库的存储管理器
type RedisCluster struct {
	KeyPrefix string // key的前缀
	HashKey   bool   // 是否对key做hash运算
	IsCache   bool   // 是否对redis连接池做缓存
}

// 添加一个测试key检测连接是否正常
func checkClusterConnectionIsOpen(cluster RedisCluster) bool {
	client := singleton(cluster.IsCache)
	testKey := "redis-test-" + uuid.NewV4().String()

	if err := client.Set(context.Background(), testKey, "test", time.Second).Err(); err != nil {
		log.Warnf("Error trying to set test key: %s", err.Error())

		return false
	}

	if _, err := client.Get(context.Background(), testKey).Result(); err != nil {
		log.Warnf("Error trying to get test key: %s", err.Error())

		return false
	}

	return true
}

func ConnectToRedis(ctx context.Context, opts *Options) {
	tick := time.NewTicker(time.Second)
	defer tick.Stop()

	clusters := []RedisCluster{
		{},
		{IsCache: true},
	}

	var ok bool // 集群连接状态
	for _, cluster := range clusters {
		connectSingleton(cluster.IsCache, opts)

		// 检测redis集群连接状态
		if !checkClusterConnectionIsOpen(cluster) {
			// 设置redis集群不可连接状态
			redisUp.Store(false)

			break
		}
		ok = true
	}

	redisUp.Store(ok)

again:
	for {
		select {
		case <-ctx.Done():
			// 正常退出
			return
		case <-tick.C:
			// 定时检测集群状态，自动连接集群
			if !shouldConnect() {
				continue
			}

			for _, cluster := range clusters {
				connectSingleton(cluster.IsCache, opts)

				if !checkClusterConnectionIsOpen(cluster) {
					redisUp.Store(false)

					goto again
				}
			}

			redisUp.Store(true)
		}
	}
}
