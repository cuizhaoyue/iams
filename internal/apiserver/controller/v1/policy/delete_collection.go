package policy

import (
	"github.com/cuizhaoyue/iams/internal/pkg/middleware"
	"github.com/cuizhaoyue/iams/pkg/log"
	"github.com/cuizhaoyue/toolkit/core"
	"github.com/gin-gonic/gin"
	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"
)

// DeleteCollection delete policies by policy names.
func (p *PolicyController) DeleteCollection(c *gin.Context) {
	log.L(c).Info("batch delete policy function called.")

	// 从query参数中获取要删除的数据列表，删除用户下指定的policy列表
	if err := p.srv.Policies().DeleteCollection(c, c.GetString(middleware.UsernameKey),
		c.QueryArray("name"), metav1.DeleteOptions{}); err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, nil)
}
