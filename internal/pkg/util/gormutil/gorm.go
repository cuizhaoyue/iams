package gormutil

// DefaultLimit 定义检索记录的默认数量
const DefaultLimit = 1000

// LimitAndOffset 包含offset和limit字段
type LimitAndOffset struct {
	Offset int
	Limit  int
}

// Unpointer 如果offset/limit没有被设置，使用默认值填充LimitAndOffset，否则使用传入值.
func Unpointer(offset *int64, limit *int64) *LimitAndOffset {
	// 设置默认值
	var o, l int = 0, DefaultLimit

	if offset != nil {
		o = int(*offset)
	}
	if limit != nil {
		l = int(*limit)
	}

	return &LimitAndOffset{
		Offset: o,
		Limit:  l,
	}
}
