package options

import (
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/cuizhaoyue/iams/internal/pkg/server"

	"github.com/spf13/pflag"
)

// JWTOptions 包含jwt功能相关的配置选项.
type JWTOptions struct {
	Realm      string        `json:"realm,omitempty"       mapstructure:"realm"`
	Key        string        `json:"key,omitempty"         mapstructure:"key"`
	Timeout    time.Duration `json:"timeout,omitempty"     mapstructure:"timeout"`
	MaxRefresh time.Duration `json:"max-refresh,omitempty" mapstructure:"max-refresh"`
}

// NewJWTOptions 创建带有默认参数的配置选项.
func NewJWTOptions() *JWTOptions {
	defaults := server.NewConfig()

	return &JWTOptions{
		Realm:      defaults.Jwt.Realm,
		Key:        defaults.Jwt.Key,
		Timeout:    defaults.Jwt.Timeout,
		MaxRefresh: defaults.Jwt.MaxRefresh,
	}
}

// ApplyTo 把配置选项应用到服务配置
func (o *JWTOptions) ApplyTo(c *server.Config) error {
	c.Jwt = &server.JwtInfo{
		Realm:      o.Realm,
		Key:        o.Key,
		Timeout:    o.Timeout,
		MaxRefresh: o.MaxRefresh,
	}

	return nil
}

// Validate 校验参数是否合法
func (o *JWTOptions) Validate() []error {
	var errs []error

	strLength := utf8.RuneCountInString(o.Key)
	if !(strLength >= 6 && strLength <= 32) {
		errs = append(errs, fmt.Errorf("--jwt.key must greater than 5 and less than 33"))
	}

	return errs
}

// AddFlags 添加jwt功能相关的flag到指定的FlagSet中.
func (o *JWTOptions) AddFlags(fs *pflag.FlagSet) {
	if fs == nil {
		return
	}

	fs.StringVar(&o.Realm, "jwt.realm", o.Realm, "Realm name to display to the user.")
	fs.StringVar(&o.Key, "jwt.key", o.Key, "Private key used to sign jwt token.")
	fs.DurationVar(&o.Timeout, "jwt.timeout", o.Timeout, "JWT token timeout.")

	fs.DurationVar(&o.MaxRefresh, "jwt.max-refresh", o.MaxRefresh, ""+
		"This field allows clients to refresh their token until MaxRefresh has passed.")
}
