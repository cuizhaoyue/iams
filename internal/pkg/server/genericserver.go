package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/marmotedu/errors"

	"github.com/cuizhaoyue/iams/internal/pkg/middleware"
	"golang.org/x/sync/errgroup"

	"github.com/cuizhaoyue/iams/pkg/log"

	"github.com/marmotedu/component-base/pkg/version"

	"github.com/cuizhaoyue/toolkit/core"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"

	ginprometheus "github.com/zsais/go-gin-prometheus"
)

// GenericAPIServer 是一个通用服务.
type GenericAPIServer struct {
	// SecureServingInfo 保存TLS服务的配置
	SecureServingInfo *SecureServingInfo
	// InsecureServingInfo 保存http服务的配置
	InsecureServingInfo *InsecureServingInfo
	// 需要加载的中间件
	middlewares []string
	// 服务关闭超时时间
	ShutdownTimeout time.Duration

	healthz       bool
	enableProfile bool
	enableMetrics bool

	*gin.Engine
	insecureServer *http.Server
	secureServer   *http.Server
}

// 对GenericAPIServer执行初始化操作.
func initGenericAPIServer(server *GenericAPIServer) {
	// 执行初始化操作, 包括设置debug日志格式、安装中间件、安装通用api.
	server.Setup()
	server.InstallMiddlewares()
	server.InstallAPIs()
}

// Setup 为gin.Engine做一些设置工作.
func (s *GenericAPIServer) Setup() {
	// 设置debug日志输出格式.
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Infof(
			"%-6s %-s ---> %s (%d handlers)",
			httpMethod,
			absolutePath,
			handlerName,
			nuHandlers,
		)
	}
}

// InstallMiddlewares 安装通用中间件.
func (s *GenericAPIServer) InstallMiddlewares() {
	// 安装两个必要的中间件
	s.Use(middleware.RequestID())
	s.Use(middleware.Context())

	// 安装自定义的中间件
	for _, m := range s.middlewares {
		mw, ok := middleware.Middlewares[m]
		if !ok {
			log.Warnf("can not find middleware: %s", m)
			continue
		}

		log.Infof("install middleware: %s", m)
		s.Use(mw)
	}
}

// InstallAPIs 安装通用的api.
// 包括健康检查api、暴露Metrics的api、启用性能分析的api.
func (s *GenericAPIServer) InstallAPIs() {
	// install healthz handler
	if s.healthz {
		s.GET("/healthz", func(c *gin.Context) {
			core.WriteResponse(c, nil, map[string]string{"status": "ok"})
		})
	}

	// install metric handler
	if s.enableMetrics {
		prometheus := ginprometheus.NewPrometheus("gin")
		prometheus.Use(s.Engine)
	}

	// install pprof handler
	if s.enableProfile {
		pprof.Register(s.Engine)
	}

	s.GET("/version", func(c *gin.Context) {
		core.WriteResponse(c, nil, version.Get())
	})
}

// Run 会生成 http 服务器。仅当初始化时无法侦听端口时，它才会返回.
func (s *GenericAPIServer) Run() error {
	// 对于可伸缩性，请在此处使用自定义 HTTP 配置模式
	s.insecureServer = &http.Server{
		Addr:    s.InsecureServingInfo.Address,
		Handler: s,
	}
	s.secureServer = &http.Server{
		Addr:    s.SecureServingInfo.Address(),
		Handler: s,
	}

	var eg errgroup.Group

	// 在goroutine中初始化服务，所以它不会阻塞下面的优雅关闭服务.
	eg.Go(func() error {
		log.Infof(
			"Start to listening the incoming requests on http address: %s",
			s.InsecureServingInfo.Address,
		)

		if err := s.insecureServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err.Error())

			return err
		}

		log.Infof("Server on %s stopped", s.InsecureServingInfo.Address)

		return nil
	})

	eg.Go(func() error {
		cert, key := s.SecureServingInfo.CertKey.CertFile, s.SecureServingInfo.CertKey.KeyFile
		if err := s.secureServer.ListenAndServeTLS(cert, key); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err.Error())

			return err
		}

		log.Infof("Server on %s stopped", s.SecureServingInfo.Address())

		return nil
	})

	// 执行ping服务，做健康检查，确保10s内路由可以正常工作.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if s.healthz {
		if err := s.ping(ctx); err != nil {
			return err
		}
	}

	// 等待启动http/https服务的goroutine执行完成
	if err := eg.Wait(); err != nil {
		log.Fatal(err.Error())
	}

	return nil
}

// Close 优雅关闭服务
func (s *GenericAPIServer) Close() {
	// context用来通知服务它有10s的时间完成正在处理的请求.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := s.secureServer.Shutdown(ctx); err != nil {
		log.Warnf("Shutdown secure server failed: %s", err.Error())
	}
	if err := s.insecureServer.Shutdown(ctx); err != nil {
		log.Warnf("Shutdown insecure server failed: %s", err.Error())
	}
}

// 通过ping http服务来确认路由正在工作.
func (s *GenericAPIServer) ping(ctx context.Context) error {
	// 当服务监听在所有网卡时，设置请求ip为127.0.0.1
	// 当HTTP服务监听在指定网卡时，请求指定的网卡ip.
	url := fmt.Sprintf("http://%s/healthz", s.InsecureServingInfo.Address)
	if strings.Contains(s.InsecureServingInfo.Address, "0.0.0.0") {
		url = fmt.Sprintf(
			"http://127.0.0.1:%s/healthz",
			strings.Split(s.InsecureServingInfo.Address, ":")[1],
		)
	}

	for {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return err
		}

		// Ping the server by sending a GET request to `/healthz`

		resp, err := http.DefaultClient.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			log.Infof("The router has been deployed successfully.")

			resp.Body.Close()

			return nil
		}

		// 检查失败后暂停1s再次检查
		log.Info("Waiting for the router, retry in 1 second。")
		time.Sleep(time.Second * 1)

		select {
		case <-ctx.Done():
			log.Fatal("can not ping http server within the speciied the time interval.")
		default:
		}
	}
}
