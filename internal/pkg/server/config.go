package server

import (
	"net"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	// RecommendedHomeDir 定义默认存放iam服务配置的目录.
	RecommendedHomeDir = ".iam"

	// RecomendedEnvPrefix 定义iam服务使用的环境变量的前缀.
	RecomendedEnvPrefix = "IAM"
)

// Config 用于配置GenericAPIServer的通用配置.它的成员大致按照重要性排序.
type Config struct {
	SecureServing   *SecureServingInfo
	InsecureServing *InsecureServingInfo
	Jwt             *JwtInfo
	Mode            string
	Middlewares     []string
	Healthz         bool
	EnableProfile   bool
	EnableMetrics   bool
}

// SecureServingInfo 保存tls服务的配置.
type SecureServingInfo struct {
	BindAddress string
	BindPort    int
	CertKey     CertKey
}

// CertKey 包含证书文件配置相关的选项.
type CertKey struct {
	// CertFile is a file containing a PEM-encoded certificate, and possibly the complete certificate chain.
	CertFile string
	// KeyFile is a file containing a PEM-encoded private key for the certificate specified by CertFile
	KeyFile string
}

// Address 连接主机ip和端口, 例如: 0.0.0.0:8443.
func (s *SecureServingInfo) Address() string {
	return net.JoinHostPort(s.BindAddress, strconv.Itoa(s.BindPort))
}

// InsecureServingInfo 保存http服务的配置.
type InsecureServingInfo struct {
	Address string
}

// JwtInfo 定义了用来创建jwt认证中间件的字段.
type JwtInfo struct {
	// defaults to "iam jwt"
	Realm string
	// defaults to empty
	Key string
	// defaults to one hour
	Timeout time.Duration
	// defaults to zero
	MaxRefresh time.Duration
}

// NewConfig 创建默认配置.
func NewConfig() *Config {
	return &Config{
		Jwt: &JwtInfo{
			Realm:      "iam jwt",
			Timeout:    time.Hour * 1,
			MaxRefresh: time.Hour * 1,
		},
		Mode:          gin.ReleaseMode,
		Middlewares:   []string{},
		Healthz:       true,
		EnableProfile: true,
		EnableMetrics: true,
	}
}

// Complete 对通用配置Config执行补全操作.
func (c *Config) Complete() *CompletedConfig {
	return &CompletedConfig{c}
}

// CompletedConfig 是GenericAPIServer的完整配置.
type CompletedConfig struct {
	*Config
}

func (c *CompletedConfig) New() (*GenericAPIServer, error) {
	// 执行gin.New前设置服务启动模式
	gin.SetMode(c.Mode)

	s := &GenericAPIServer{
		SecureServingInfo:   c.SecureServing,
		InsecureServingInfo: c.InsecureServing,
		middlewares:         c.Middlewares,
		enableMetrics:       c.EnableMetrics,
		enableProfile:       c.EnableProfile,
		healthz:             c.Healthz,
		Engine:              gin.New(),
	}

	initGenericAPIServer(s)

	return s, nil
}
