package config

import "github.com/cuizhaoyue/iams/internal/apiserver/options"

// Config 是iam apiserver服务的配置.
type Config struct {
	*options.Options
}

// CreateConfigFromOptions 根据提供的Options创建服务配置.
func CreateConfigFromOptions(opts *options.Options) (*Config, error) {
	return &Config{opts}, nil
}
