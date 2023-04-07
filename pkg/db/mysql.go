package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

// Options 定义了mysql数据库的选项.
type Options struct {
	Host                  string        // mysql host地址，ip[:port]形式
	Username              string        // 访问mysql的username
	Password              string        // 访问mysql的password
	Database              string        // 要访问的数据库
	MaxIdleConnections    int           // mysql最大空闲连接数,推荐100
	MaxOpenConnections    int           // mysql最大连接数，推荐100
	MaxConnectionLifeTime time.Duration // mysql的空闲连接最大存活时间，推荐10s
	LogLevel              int           // 日志等级
	Logger                logger.Interface
}

// New 根据给出的Options创建*gorm.DB实例
func New(opts *Options) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=%t&loc=%s",
		opts.Username,
		opts.Password,
		opts.Host,
		opts.Database,
		true,
		"Local")

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: opts.Logger,
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 设置MySQL最大连接数
	sqlDB.SetMaxOpenConns(opts.MaxOpenConnections)

	// 设置最大空闲连接数
	sqlDB.SetMaxIdleConns(opts.MaxIdleConnections)

	// 设置MySQL空闲连接最大存活时间
	sqlDB.SetConnMaxLifetime(opts.MaxConnectionLifeTime)

	return db, nil
}
