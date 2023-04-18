package secret

import (
	"github.com/cuizhaoyue/iams/internal/pkg/middleware"
	"github.com/cuizhaoyue/iams/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/marmotedu/component-base/pkg/core"
	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"
)

// DeleteCollection delete secrets by secret names.
func (s *SecretController) DeleteCollection(c *gin.Context) {
	log.L(c).Info("batch delete policy function called.")

	// 从Query参数中获取要删除的secret名称列表，删除username下这些指定的secret数据
	if err := s.srv.Secrets().DeleteCollection(
		c,
		c.GetString(middleware.UsernameKey),
		c.QueryArray("name"),
		metav1.DeleteOptions{},
	); err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, nil)
}
