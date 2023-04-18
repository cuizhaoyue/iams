package user

import (
	"github.com/cuizhaoyue/iams/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/marmotedu/component-base/pkg/core"
	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"
)

// Get get an user by the user identifier.
func (u *UserController) Get(c *gin.Context) {
	log.L(c).Info("get user function called.")

	// 从path参数中获取用户名，根据用户名从数据库中获取user数据
	user, err := u.srv.Users().Get(c, c.Param("name"), metav1.GetOptions{})
	if err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, user)
}
