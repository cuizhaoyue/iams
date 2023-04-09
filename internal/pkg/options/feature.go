package options

import (
	"github.com/cuizhaoyue/iams/internal/pkg/server"
	"github.com/spf13/pflag"
)

// FeatureOptions 包含apiserver额外功能的配置选项.
type FeatureOptions struct {
	EnableProfiling bool `json:"profiling"      mapstructure:"enable-profiling"`
	EnableMetrics   bool `json:"enable-metrics" mapstructure:"enable-metrics"`
}

// NewFeatureOptions 创建一个默认参数的FeatureOptions对象。
func NewFeatureOptions() *FeatureOptions {
	defaults := server.NewConfig()

	return &FeatureOptions{
		EnableProfiling: defaults.EnableProfile,
		EnableMetrics:   defaults.EnableMetrics,
	}
}

// ApplyTo 把配置选项应用到服务配置.
func (o *FeatureOptions) ApplyTo(c *server.Config) error {
	c.EnableProfile = o.EnableProfiling
	c.EnableMetrics = o.EnableMetrics

	return nil
}

// Validate 校验传入参数是否合法.
func (o *FeatureOptions) Validate() []error {
	return []error{}
}

// AddFlags 添加feature相关的参数到指定的FlagSet中.
func (o *FeatureOptions) AddFlags(fs *pflag.FlagSet) {
	if fs == nil {
		return
	}

	fs.BoolVar(&o.EnableProfiling, "feature.profiling", o.EnableProfiling,
		"Enable profiling via web interface host:port/debug/pprof/")

	fs.BoolVar(&o.EnableMetrics, "feature.enable-metrics", o.EnableMetrics,
		"Enables metrics on the apiserver at /metrics")
}
