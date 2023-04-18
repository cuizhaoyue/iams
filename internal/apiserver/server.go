package apiserver

import (
	"context"
	"fmt"

	pb "github.com/marmotedu/api/proto/apiserver/v1"
	"google.golang.org/grpc/reflection"

	"github.com/cuizhaoyue/iams/internal/apiserver/controller/v1/cache"

	"github.com/cuizhaoyue/iams/internal/apiserver/store"
	"github.com/cuizhaoyue/iams/internal/apiserver/store/mysql"

	"github.com/cuizhaoyue/iams/pkg/storage"

	"google.golang.org/grpc"

	"github.com/cuizhaoyue/iams/pkg/log"

	"google.golang.org/grpc/credentials"

	"github.com/cuizhaoyue/iams/internal/apiserver/config"
	genericoptions "github.com/cuizhaoyue/iams/internal/pkg/options"
	genericapiserver "github.com/cuizhaoyue/iams/internal/pkg/server"
	"github.com/cuizhaoyue/iams/pkg/shutdown"
	"github.com/cuizhaoyue/iams/pkg/shutdown/shutdownmanagers/posixsignal"
)

// apiserver应用配置
type apiServer struct {
	genericAPIServer *genericapiserver.GenericAPIServer // 通用api服务
	gRPCAPIServer    *grpcAPIServer                     // grpc服务
	gs               *shutdown.GracefuleShutdown        // 负责服务优雅关闭
	redisOptions     *genericoptions.RedisOptions       // redis配置选项
}

// 准备好的apiserver服务
type preparedAPIServer struct {
	*apiServer
}

func createAPIServer(cfg *config.Config) (*apiServer, error) {
	// 设置优雅关停
	gs := shutdown.New()
	gs.AddShutdownManager(posixsignal.NewPosixSignalManager())

	genericConfig, err := buildGenericConfig(cfg)
	if err != nil {
		return nil, err
	}

	extraConfig, err := buildExtraConfig(cfg)
	if err != nil {
		return nil, err
	}

	// 根据通用配置生成通用api服务.
	genericServer, err := genericConfig.Complete().New()
	if err != nil {
		return nil, err
	}

	extraServer, err := extraConfig.complete().New()
	if err != nil {
		return nil, err
	}

	server := &apiServer{
		genericAPIServer: genericServer,
		gRPCAPIServer:    extraServer,
		gs:               gs,
		redisOptions:     cfg.RedisOptions,
	}

	return server, nil
}

// PrepareRun 执行准备工作，包含初始化操作，如数据库初始化、安装业务相关的gin中间件、安装restful路由.
func (s *apiServer) PrepareRun() preparedAPIServer {
	// 初始化路由
	initRouter(s.genericAPIServer.Engine)

	// 初始化redis服务
	s.initRedisStore()

	// 添加服务结束时的关闭操作.
	s.gs.AddShutdownCallback(shutdown.ShutdownFunc(func(string) error {
		mysqlStore, _ := mysql.GetMySQLFactoryOr(nil)
		if mysqlStore != nil {
			// 关闭mysql连接池
			_ = mysqlStore.Close()
		}
		// 关闭grpc服务和通用apiserver服务.
		s.gRPCAPIServer.Close()
		s.genericAPIServer.Close()

		return nil
	}))

	return preparedAPIServer{s}
}

// Run 运行通用服务
func (s preparedAPIServer) Run() error {
	// 运行grpc服务
	go s.gRPCAPIServer.Run()

	// 启动shutdown manager
	if err := s.gs.Start(); err != nil {
		log.Fatalf("start shutdown manager failed: %s", err.Error())
	}

	// 启动通用api服务
	return s.genericAPIServer.Run()
}

// 初始化redis
func (s *apiServer) initRedisStore() {
	ctx, cancel := context.WithCancel(context.Background())
	s.gs.AddShutdownCallback(shutdown.ShutdownFunc(func(s string) error {
		cancel()

		return nil
	}))

	cfg := &storage.Options{
		Host:                  s.redisOptions.Host,
		Port:                  s.redisOptions.Port,
		Addrs:                 s.redisOptions.Addrs,
		MasterName:            s.redisOptions.MasterName,
		Username:              s.redisOptions.Username,
		Password:              s.redisOptions.Password,
		Database:              s.redisOptions.Database,
		MaxIdle:               s.redisOptions.MaxIdle,
		MaxActive:             s.redisOptions.MaxActive,
		Timeout:               s.redisOptions.Timeout,
		EnableCluster:         s.redisOptions.EnableCluster,
		UseSSL:                s.redisOptions.UseSSL,
		SSLInsecureSkipVerify: s.redisOptions.SSLInsecureSkipVerify,
	}

	// 启动一个goroutine执行redis连接操作.
	go storage.ConnectToRedis(ctx, cfg)
}

// 根据apiserver应用配置生成通用配置.
func buildGenericConfig(cfg *config.Config) (genericConfig *genericapiserver.Config, lastErr error) {
	// 创建默认的通用配置
	genericConfig = genericapiserver.NewConfig()

	if lastErr = cfg.GenericServerRunOptions.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	if lastErr = cfg.FeatureOptions.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	if lastErr = cfg.SecureServing.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	if lastErr = cfg.InsecureServing.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	return genericConfig, nil
}

// ExtraConfig 定义了iam-apiserver的额外配置.
type ExtraConfig struct {
	Addr         string
	MaxMsgSize   int
	ServerCert   genericoptions.GeneratableKeyCert
	mysqlOptions *genericoptions.MySQLOptions
}

// 完整的ExtraConfig
type completedExtraConfig struct {
	*ExtraConfig
}

// nolint: unparam
func buildExtraConfig(cfg *config.Config) (*ExtraConfig, error) {
	return &ExtraConfig{
		Addr:         fmt.Sprintf("%s:%d", cfg.GRPCOptions.BindAddress, cfg.GRPCOptions.BindPort),
		MaxMsgSize:   cfg.GRPCOptions.MaxMsgSize,
		ServerCert:   cfg.SecureServing.ServerCert,
		mysqlOptions: cfg.MySQLOptions,
	}, nil
}

// 填充必要的字段，生成完整的配置.
func (c *ExtraConfig) complete() *completedExtraConfig {
	if c.Addr == "" {
		c.Addr = "127.0.0.1:8081"
	}

	return &completedExtraConfig{c}
}

// New 创建一个grpcAPIServer实例, 并向其中注册需要的服务.
func (c *completedExtraConfig) New() (*grpcAPIServer, error) {
	creds, err := credentials.NewServerTLSFromFile(c.ServerCert.CertKey.CertFile, c.ServerCert.CertKey.KeyFile)

	if err != nil {
		log.Fatalf("Failed to generate credentials %s", err.Error())
	}

	opts := []grpc.ServerOption{grpc.MaxRecvMsgSize(c.MaxMsgSize), grpc.Creds(creds)}
	grpcServer := grpc.NewServer(opts...)

	// 注册缓存服务到grpc服务
	mysqlStore, _ := mysql.GetMySQLFactoryOr(c.mysqlOptions)
	store.SetClient(mysqlStore)
	cacheIns, err := cache.GetCacheInsOr(mysqlStore)
	if err != nil {
		log.Fatalf("Failed to get cache instance: %s", err.Error())
	}

	pb.RegisterCacheServer(grpcServer, cacheIns)

	reflection.Register(grpcServer)

	return &grpcAPIServer{grpcServer, c.Addr}, nil
}
