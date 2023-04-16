package options

import (
	"fmt"
	"path"

	"github.com/cuizhaoyue/iams/internal/pkg/server"

	"github.com/spf13/pflag"
)

// SecureServingOptions 包含和HTTPS服务启动配置相关的选项.
type SecureServingOptions struct {
	BindAddress string `json:"bind-address,omitempty" mapstructure:"bind-address"`
	// 当Listener被设置时，BindPort会被忽略.
	BindPort int `json:"bind-port,omitempty"    mapstructure:"bind-port"`
	// Required设置为true时，BindPort不能为0.
	Required bool
	// ServerCert 是提供安全流量的TLS证书信息.
	ServerCert GeneratableKeyCert `json:"tls"                    mapstructure:"tls"`
}

// GeneratableKeyCert 包含和证书相关的配置选项.
type GeneratableKeyCert struct {
	CertKey       CertKey `json:"cert-key"                 mapstructure:"cert-key"`
	CertDirectory string  `json:"cert-directory,omitempty" mapstructure:"cert-directory"`
	PairName      string  `json:"pair-name,omitempty"      mapstructure:"pair-name"`
}

// CertKey 包含证书文件配置相关的选项.
type CertKey struct {
	CertFile string `json:"cert-file,omitempty"        mapstructure:"cert-file"`
	KeyFile  string `json:"private-key-file,omitempty" mapstructure:"private-key-file"`
}

// NewSecureServingOptions 创建带有默认参数的配置选项.
func NewSecureServingOptions() *SecureServingOptions {
	return &SecureServingOptions{
		BindAddress: "0.0.0.0",
		BindPort:    8443,
		Required:    true,
		ServerCert: GeneratableKeyCert{
			PairName:      "iam",
			CertDirectory: "/var/run/iam",
		},
	}
}

// ApplyTo 应用配置选项到服务配置.
func (o *SecureServingOptions) ApplyTo(c *server.Config) error {
	c.SecureServing = &server.SecureServingInfo{
		BindAddress: o.BindAddress,
		BindPort:    o.BindPort,
		CertKey: server.CertKey{
			CertFile: o.ServerCert.CertKey.CertFile,
			KeyFile:  o.ServerCert.CertKey.KeyFile,
		},
	}

	return nil
}

// Validate 校验传入的参数是否合法.
func (o *SecureServingOptions) Validate() []error {
	if o == nil {
		return nil
	}

	var errs []error
	if o.Required && o.BindPort < 1 || o.BindPort > 65535 {
		errs = append(errs, fmt.Errorf(
			"--secure.bind-port %v must be between 1 and 65535, inclusive. It cannot be turned off with 0",
			o.BindPort,
		))
	}

	return errs
}

// AddFlags 添加https服务配置相关的选项到指定的FlagSet中.
func (o *SecureServingOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.BindAddress, "secure.bind-address", o.BindAddress, ""+
		"The IP address on which to listen for the --secure.bind-port port. The "+
		"associated interface(s) must be reachable by the rest of the engine, and by CLI/web "+
		"clients. If blank, all interfaces will be used (0.0.0.0 for all IPv4 interfaces and :: for all IPv6 interfaces).")

	desc := "The port on which to serve HTTPS with authentication and authorization."
	if o.Required {
		desc += " It cannot be switched off with 0."
	} else {
		desc += " If 0, don't serve HTTPS at all."
	}
	fs.IntVar(&o.BindPort, "secure.bind-port", o.BindPort, desc)
	fs.StringVar(&o.ServerCert.CertDirectory, "secure.tls.cert-dir", o.ServerCert.CertDirectory, ""+
		"The directory where the TLS certs are located. "+
		"If --secure.tls.cert-key.cert-file and --secure.tls.cert-key.private-key-file are provided, "+
		"this flag will be ignored.")
	fs.StringVar(&o.ServerCert.PairName, "secure.tls.pair-name", o.ServerCert.PairName, ""+
		"The name which will be used with --secure.tls.cert-dir to make a cert and key filenames. "+
		"It becomes <cert-dir>/<pair-name>.crt and <cert-dir>/<pair-name>.key")
	fs.StringVar(
		&o.ServerCert.CertKey.CertFile,
		"secure.tls.cert-key.cert-file",
		o.ServerCert.CertKey.CertFile,
		""+
			"File containing the default x509 Certificate for HTTPS. (CA cert, if any, concatenated "+
			"after server cert).",
	)
	fs.StringVar(
		&o.ServerCert.CertKey.KeyFile,
		"secure.tls.cert-key.private-key-file",
		o.ServerCert.CertKey.KeyFile,
		""+
			"File containing the default x509 private key matching --secure.tls.cert-key.cert-file.",
	)
}

// Complete 填充任何必要但是没有设置的字段.
func (o *SecureServingOptions) Complete() error {
	if o == nil || o.BindPort == 0 {
		return nil
	}

	keyCert := &o.ServerCert.CertKey
	if keyCert.CertFile != "" || keyCert.KeyFile != "" {
		return nil
	}

	if o.ServerCert.CertDirectory != "" {
		if o.ServerCert.PairName == "" {
			return fmt.Errorf("--secure.tls.pair-name is required if --secure.tls.cert-dir is set")
		}
		keyCert.CertFile = path.Join(o.ServerCert.CertDirectory, o.ServerCert.PairName+".crt")
		keyCert.KeyFile = path.Join(o.ServerCert.CertDirectory, o.ServerCert.PairName+".key")
	}

	return nil
}
