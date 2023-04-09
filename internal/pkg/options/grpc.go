package options

import (
	"fmt"

	"github.com/spf13/pflag"
)

type GRPCOptions struct {
	BindAddress string `json:"bind-address,omitempty" mapstructure:"bind-address"`
	BindPort    int    `json:"bind-port,omitempty"    mapstructure:"bind-port"`
	MaxMsgSize  int    `json:"max-msg-size,omitempty" mapstructure:"max-msg-size"`
}

// NewGRPCOptions 创建grpc配置的默认选项.
func NewGRPCOptions() *GRPCOptions {
	return &GRPCOptions{
		BindAddress: "0.0.0.0",
		BindPort:    8081,
		MaxMsgSize:  4 * 1024 * 1024,
	}
}

// Validate 难和解析命令行传入的参数是否合法.
func (o *GRPCOptions) Validate() []error {
	var errs []error

	if o.BindPort < 0 || o.BindPort > 65535 {
		errs = append(errs, fmt.Errorf(
			"--insecure-port %v must be between 0 and 65535, inclusive. 0 for turning off insecure (HTTP) port",
			o.BindPort,
		))
	}

	return errs
}

// AddFlags 添加和grpc服务相关的flag到指定的FlagSet中.
func (o *GRPCOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.BindAddress, "grpc.bind-address", o.BindAddress, ""+
		"The IP address on which to serve the --grpc.bind-port(set to 0.0.0.0 for all IPv4 interfaces and :: for all IPv6 interfaces).")

	fs.IntVar(&o.BindPort, "grpc.bind-port", o.BindPort, ""+
		"The port on which to serve unsecured, unauthenticated grpc access. It is assumed "+
		"that firewall rules are set up such that this port is not reachable from outside of "+
		"the deployed machine and that port 443 on the iam public address is proxied to this "+
		"port. This is performed by nginx in the default setup. Set to zero to disable.")

	fs.IntVar(&o.MaxMsgSize, "grpc.max-msg-size", o.MaxMsgSize, "gRPC max message size.")
}
