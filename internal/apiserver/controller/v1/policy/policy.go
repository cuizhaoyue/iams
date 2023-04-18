package policy

import (
	srvv1 "github.com/cuizhaoyue/iams/internal/apiserver/service/v1"
	"github.com/cuizhaoyue/iams/internal/apiserver/store"
)

// PolicyController 创建了关于策略资源请求的处理器
type PolicyController struct {
	srv srvv1.Service
}

func NewPolicyController(store store.Factory) *PolicyController {
	return &PolicyController{
		srv: srvv1.NewService(store),
	}
}
