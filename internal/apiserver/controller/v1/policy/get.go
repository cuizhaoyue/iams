package policy

import (
	"github.com/cuizhaoyue/iams/internal/pkg/middleware"
	"github.com/cuizhaoyue/iams/pkg/log"
	"github.com/cuizhaoyue/toolkit/core"
	"github.com/gin-gonic/gin"
	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"
)

// Get return policy by the policy identifier.
func (p *PolicyController) Get(c *gin.Context) {
	log.L(c).Info("get policy function called.")
	// 从path参数中获取policy的name，通过username和name获取policy数据
	pol, err := p.srv.Policies().Get(c, c.GetString(middleware.UsernameKey), c.Param("name"), metav1.GetOptions{})
	if err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, pol)
}
