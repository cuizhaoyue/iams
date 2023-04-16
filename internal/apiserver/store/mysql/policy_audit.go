package mysql

import (
	"context"
	"time"

	"github.com/cuizhaoyue/iams/internal/apiserver/store"
	"gorm.io/gorm"
)

type policyAudit struct {
	db *gorm.DB
}

var _ store.PolicyAuditStore = &policyAudit{}

func newPolicyAudit(ds *datastore) *policyAudit {
	return &policyAudit{ds.db}
}

// ClearOutdated 清理
func (p *policyAudit) ClearOutdated(ctx context.Context, maxReserveDays int) (int64, error) {
	date := time.Now().AddDate(0, 0, -maxReserveDays).Format(time.DateTime)

	d := p.db.Exec("delete from policy_audit where deletedAt < ?", date)

	return d.RowsAffected, d.Error
}
