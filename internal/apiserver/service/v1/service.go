package v1

import "github.com/cuizhaoyue/iams/internal/apiserver/store"

// Service 定义返回资源数据的函数
type Service interface {
	Users() UserSrv
	Secret() SecretSrv
	Policies() PolicySrv
}

var _ Service = &service{}

type service struct {
	store store.Factory
}

func NewService(store store.Factory) Service {
	return &service{store: store}
}

func (s *service) Users() UserSrv {
	return newUser(s)
}

func (s *service) Secret() SecretSrv {
	return newSecrets(s)
}

func (s *service) Policies() PolicySrv {
	return newPolicies(s)
}
