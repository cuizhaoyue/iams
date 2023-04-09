package options

import (
	"time"

	"github.com/spf13/pflag"
)

// MySQLOptions 定义mysql数据库的配置选项
type MySQLOptions struct {
	Host                  string        `json:"host,omitempty"                     mapstructure:"host"`
	Username              string        `json:"username,omitempty"                 mapstructure:"username"`
	Password              string        `json:"password,omitempty"                 mapstructure:"password"`
	Database              string        `json:"database,omitempty"                 mapstructure:"database"`
	MaxIdleConnections    int           `json:"max-idle-connections,omitempty"     mapstructure:"max-idle-connections"`
	MaxOpenConnections    int           `json:"max-open-connections,omitempty"     mapstructure:"max-open-connections"`
	MaxConnectionLifeTime time.Duration `json:"max-connection-life-time,omitempty" mapstructure:"max-connection-life-time"`
	LogLevel              int           `json:"log-level,omitempty"                mapstructure:"log-level"`
}

func NewMySQLOptions() *MySQLOptions {
	return &MySQLOptions{
		Host:                  "127.0.0.1:3306",
		Username:              "",
		Password:              "",
		Database:              "",
		MaxIdleConnections:    100,
		MaxOpenConnections:    100,
		MaxConnectionLifeTime: time.Duration(10) * time.Second,
		LogLevel:              1, // Silent
	}
}

// Validate 验证传给MySQLOptions的flag.
func (o *MySQLOptions) Validate() []error {
	var errs []error

	return errs
}

// AddFlags 添加指定的和mysql存储相关的flag到指定的FlagSet中.
func (o *MySQLOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Host, "mysql.host", o.Host, ""+
		"MySQL service host address. If left blank, the following related mysql options will be ignored.")

	fs.StringVar(&o.Username, "mysql.username", o.Username, ""+
		"Username for access to mysql service.")

	fs.StringVar(&o.Password, "mysql.password", o.Password, ""+
		"Password for access to mysql, should be used pair with password.")

	fs.StringVar(&o.Database, "mysql.database", o.Database, ""+
		"Database name for the server to use.")

	fs.IntVar(&o.MaxIdleConnections, "mysql.max-idle-connections", o.MaxOpenConnections, ""+
		"Maximum idle connections allowed to connect to mysql.")

	fs.IntVar(&o.MaxOpenConnections, "mysql.max-open-connections", o.MaxOpenConnections, ""+
		"Maximum open connections allowed to connect to mysql.")

	fs.DurationVar(
		&o.MaxConnectionLifeTime,
		"mysql.max-connection-life-time",
		o.MaxConnectionLifeTime,
		""+
			"Maximum connection life time allowed to connect to mysql.",
	)

	fs.IntVar(&o.LogLevel, "mysql.log-mode", o.LogLevel, ""+
		"Specify gorm log level.")
}
