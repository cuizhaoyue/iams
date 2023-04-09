package options

import (
	"fmt"
	"net"
	"strconv"

	"github.com/cuizhaoyue/iams/internal/pkg/server"

	"github.com/spf13/pflag"
)

// InsecureServingOptions 助于创建未认证、未授权、不安全的端口.
type InsecureServingOptions struct {
	BindAddress string `json:"bind-address,omitempty" mapstructure:"bind-address"`
	BindPort    int    `json:"bind-port,omitempty"    mapstructure:"bind-port"`
}

// NewInsecureServingOptions 创建InsecureServingOptions的默认实例.
func NewInsecureServingOptions() *InsecureServingOptions {
	return &InsecureServingOptions{
		BindAddress: "127.0.0.1",
		BindPort:    8080,
	}
}

// ApplyTo 把配置选项应用到服务配置.
func (o *InsecureServingOptions) ApplyTo(c *server.Config) error {
	c.InsecureServing = &server.InsecureServingInfo{
		Address: net.JoinHostPort(o.BindAddress, strconv.Itoa(o.BindPort)),
	}

	return nil
}

// Validate 校验命令行传入参数的合法性.
func (o *InsecureServingOptions) Validate() []error {
	var errs []error

	if o.BindPort < 0 || o.BindPort > 65535 {
		errs = append(errs, fmt.Errorf(
			"--insecure.bind-port %v must be between 0 and 65535, inclusive. 0 for turning off insecure (HTTP) port",
			o.BindPort,
		))
	}

	return errs
}

// AddFlags 添加和insecure服务相关的flag到指定的FlagSet中.
func (o *InsecureServingOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.BindAddress, "insecure.bind-address", o.BindAddress, ""+
		"The IP address on which to serve the --insecure.bind-port "+
		"(set to 0.0.0.0 for all IPv4 interfaces and :: for all IPv6 interfaces).")
	fs.IntVar(&o.BindPort, "insecure.bind-port", o.BindPort, ""+
		"The port on which to serve unsecured, unauthenticated access. It is assumed "+
		"that firewall rules are set up such that this port is not reachable from outside of "+
		"the deployed machine and that port 443 on the iam public address is proxied to this "+
		"port. This is performed by nginx in the default setup. Set to zero to disable.")
}
