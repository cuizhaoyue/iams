package policy

import (
	"github.com/cuizhaoyue/iams/internal/pkg/middleware"
	"github.com/cuizhaoyue/iams/pkg/log"
	"github.com/cuizhaoyue/toolkit/core"
	"github.com/gin-gonic/gin"
	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"
)

// Delete deletes the policy by the policy identifier.
func (p *PolicyController) Delete(c *gin.Context) {
	log.L(c).Info("delete policy function called.")
	// 通过username和policy的name获取policy数据
	if err := p.srv.Policies().Delete(c, c.GetString(middleware.UsernameKey), c.Param("name"),
		metav1.DeleteOptions{}); err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, nil)
}
