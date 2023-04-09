package options

import (
	"github.com/cuizhaoyue/iams/internal/pkg/server"
	"github.com/spf13/pflag"
)

// ServerRunOptions 包含apiserver配置的通用选项.
type ServerRunOptions struct {
	Mode        string   `json:"mode,omitempty"        mapstructure:"mode"`
	Healthz     bool     `json:"healthz,omitempty"     mapstructure:"healthz"`
	Middlewares []string `json:"middlewares,omitempty" mapstructure:"middlewares"`
}

// NewServerRunOptions 创建带有默认参数的ServerRunOptions对象.
func NewServerRunOptions() *ServerRunOptions {
	cfg := server.NewConfig()

	return &ServerRunOptions{
		Mode:        cfg.Mode,
		Healthz:     cfg.Healthz,
		Middlewares: cfg.Middlewares,
	}
}

// ApplyTo 把配置选项应用到服务配置.
func (o *ServerRunOptions) ApplyTo(c *server.Config) error {
	c.Mode = o.Mode
	c.Healthz = o.Healthz
	c.Middlewares = o.Middlewares

	return nil
}

// Validate 校验命令行传入的参数是否合法.
func (o *ServerRunOptions) Validate() []error {
	var err []error

	return err
}

// AddFlags 添加apiserver的通用配置选项到指定的FlagSet中.
func (o *ServerRunOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Mode, "server.mode", o.Mode, ""+
		"Start the server in a specified server mode. Supported server mode: debug, test, release.")

	fs.BoolVar(&o.Healthz, "server.healthz", o.Healthz, ""+
		"Add self readiness check and install /healthz router.")

	fs.StringSliceVar(&o.Middlewares, "server.middlewares", o.Middlewares, ""+
		"List of allowed middlewares for server, comma separated. If this list is empty default middlewares will be used.")

}
