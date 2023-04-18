package user

import (
	"github.com/cuizhaoyue/iams/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/marmotedu/component-base/pkg/core"
	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"
)

// Delete 通过用户标识符删除用户，只允许admin用户调用这个函数
func (u *UserController) Delete(c *gin.Context) {
	log.L(c).Info("delete user function called.")

	// 从path参数中获取用户名，从数据库中删除对应的用户
	if err := u.srv.Users().Delete(c, c.Param("name"), metav1.DeleteOptions{Unscoped: true}); err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, nil)
}
