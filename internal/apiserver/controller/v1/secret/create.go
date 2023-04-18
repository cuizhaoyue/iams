package secret

import (
	"github.com/AlekSi/pointer"
	"github.com/cuizhaoyue/iams/internal/pkg/code"
	"github.com/cuizhaoyue/iams/internal/pkg/middleware"
	"github.com/cuizhaoyue/iams/pkg/log"
	"github.com/cuizhaoyue/toolkit/core"
	"github.com/gin-gonic/gin"
	v1 "github.com/marmotedu/api/apiserver/v1"
	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"
	"github.com/marmotedu/component-base/pkg/util/idutil"
	"github.com/marmotedu/errors"
)

const maxSecretCount = 10

// Create add new secret key pairs to the storage.
func (s *SecretController) Create(c *gin.Context) {
	log.L(c).Info("create secret function called.")

	var r v1.Secret
	// 获取请求中的数据
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, errors.WithCode(code.ErrBind, err.Error()), nil)

		return
	}

	// 校验请求数据
	if errs := r.Validate(); len(errs) != 0 {
		core.WriteResponse(c, errors.WithCode(code.ErrValidation, errs.ToAggregate().Error()), nil)

		return
	}

	// 从gin.Context中获取username
	username := c.GetString(middleware.UsernameKey)
	// 查询用户的secret列表
	secrets, err := s.srv.Secrets().List(c, username, metav1.ListOptions{
		Offset: pointer.ToInt64(0),
		Limit:  pointer.ToInt64(-1),
	})
	if err != nil {
		core.WriteResponse(c, err, nil)

		return
	}
	// 校验secret数量是否超过了最大的secret数量
	if secrets.TotalCount >= maxSecretCount {
		core.WriteResponse(c, errors.WithCode(code.ErrReachMaxCount, "secret count: %d", secrets.TotalCount), nil)

		return
	}

	// 必须设置secret中的用户名
	r.Username = username

	// 生成secretId和secret key
	r.SecretID = idutil.NewSecretID()
	r.SecretKey = idutil.NewSecretKey()

	if err := s.srv.Secrets().Create(c, &r, metav1.CreateOptions{}); err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, r)
}
