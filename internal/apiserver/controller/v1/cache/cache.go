package cache

import "github.com/cuizhaoyue/iams/internal/apiserver/store"

// Cache 定义了一个用来获取所有secrets和policies列表的缓存服务.
type Cache struct {
	store store.Factory
}
