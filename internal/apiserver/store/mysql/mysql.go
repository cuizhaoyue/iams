package mysql

import (
	"fmt"
	"sync"

	"github.com/cuizhaoyue/iams/internal/pkg/logger"
	"github.com/cuizhaoyue/iams/pkg/db"

	"github.com/cuizhaoyue/iams/internal/apiserver/store"
	genericoptions "github.com/cuizhaoyue/iams/internal/pkg/options"
	"github.com/marmotedu/errors"
	"gorm.io/gorm"
)

type datastore struct {
	db *gorm.DB
}

var _ store.Factory = &datastore{}

func (ds *datastore) Users() store.UserStore {
	return newUsers(ds)
}

func (ds *datastore) Secret() store.SecretStore {
	return newSecrets(ds)
}

func (ds *datastore) Polices() store.PolicyStore {
	return newPolicies(ds)
}

func (ds *datastore) PolicyAudit() store.PolicyAuditStore {
	return newPolicyAudit(ds)
}

func (ds *datastore) Close() error {
	db, err := ds.db.DB()
	if err != nil {
		return errors.Wrap(err, "get gorm db instance failed")
	}

	return db.Close()
}

var (
	mysqlFactory store.Factory
	once         sync.Once
)

// GetMySQLFactoryOr 通过给定的配置创建mysql工厂实例.
func GetMySQLFactoryOr(opts *genericoptions.MySQLOptions) (store.Factory, error) {
	if opts == nil && mysqlFactory == nil {
		return nil, fmt.Errorf("failed to get mysql store factory")
	}

	var (
		err   error
		dbIns *gorm.DB
	)

	once.Do(func() {
		options := &db.Options{
			Host:                  opts.Host,
			Username:              opts.Username,
			Password:              opts.Password,
			Database:              opts.Database,
			MaxIdleConnections:    opts.MaxIdleConnections,
			MaxOpenConnections:    opts.MaxOpenConnections,
			MaxConnectionLifeTime: opts.MaxConnectionLifeTime,
			LogLevel:              opts.LogLevel,
			Logger:                logger.New(opts.LogLevel),
		}

		dbIns, err = db.New(options)

		mysqlFactory = &datastore{db: dbIns}
	})

	if mysqlFactory == nil || err != nil {
		return nil, fmt.Errorf("failed to get mysql store fatory, mysqlFactory: %+v, error: %w", mysqlFactory, err)
	}

	return mysqlFactory, nil
}
