package secret

import (
	"github.com/cuizhaoyue/iams/internal/pkg/middleware"
	"github.com/cuizhaoyue/iams/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/marmotedu/component-base/pkg/core"
	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"
)

// Get get an policy by the secret identifier.
func (s *SecretController) Get(c *gin.Context) {
	log.L(c).Info("get secret function called.")

	// 多path参数中获取secret名称，通过用户名和secret名称获取secret信息
	secret, err := s.srv.Secrets().Get(c, c.GetString(middleware.UsernameKey), c.Param("name"), metav1.GetOptions{})
	if err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, secret)
}
