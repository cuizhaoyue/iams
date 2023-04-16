package store

import "context"

// PolicyAuditStore 定义了policy_audit存储接口.
type PolicyAuditStore interface {
	ClearOutdated(ctx context.Context, maxReserveDays int) (int64, error)
}
