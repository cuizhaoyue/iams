package v1

import (
	"context"

	"github.com/cuizhaoyue/iams/internal/pkg/code"
	"github.com/marmotedu/errors"

	"github.com/cuizhaoyue/iams/internal/apiserver/store"

	v1 "github.com/marmotedu/api/apiserver/v1"
	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"
)

// PolicySrv 定义处理策略相关请求的函数
type PolicySrv interface {
	Create(ctx context.Context, policy *v1.Policy, opts metav1.CreateOptions) error
	Update(ctx context.Context, policy *v1.Policy, opts metav1.UpdateOptions) error
	Delete(ctx context.Context, username string, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, username string, names []string, opts metav1.DeleteOptions) error
	Get(ctx context.Context, username string, name string, opts metav1.GetOptions) (*v1.Policy, error)
	List(ctx context.Context, username string, opts metav1.ListOptions) (*v1.PolicyList, error)
}

var _ PolicySrv = &policyService{}

type policyService struct {
	store store.Factory
}

func newPolicies(srv *service) *policyService {
	return &policyService{srv.store}
}

func (s *policyService) Create(ctx context.Context, policy *v1.Policy, opts metav1.CreateOptions) error {
	if err := s.store.Polices().Create(ctx, policy, opts); err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}

	return nil
}

func (s *policyService) Update(ctx context.Context, policy *v1.Policy, opts metav1.UpdateOptions) error {
	if err := s.store.Polices().Update(ctx, policy, opts); err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}

	return nil
}

func (s *policyService) Delete(ctx context.Context, username string, name string, opts metav1.DeleteOptions) error {
	if err := s.store.Polices().Delete(ctx, username, name, opts); err != nil {
		return err
	}

	return nil
}

func (s *policyService) DeleteCollection(
	ctx context.Context,
	username string,
	names []string,
	opts metav1.DeleteOptions,
) error {
	if err := s.store.Polices().DeleteCollection(ctx, username, names, opts); err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}

	return nil
}

func (s *policyService) Get(ctx context.Context, username string, name string, opts metav1.GetOptions) (*v1.Policy, error) {
	policy, err := s.store.Polices().Get(ctx, username, name, opts)
	if err != nil {
		return nil, err
	}

	return policy, nil
}

func (s *policyService) List(ctx context.Context, username string, opts metav1.ListOptions) (*v1.PolicyList, error) {
	policies, err := s.store.Polices().List(ctx, username, opts)
	if err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}

	return policies, nil
}
