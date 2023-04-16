package v1

import (
	"context"

	"github.com/cuizhaoyue/iams/internal/apiserver/store"

	v1 "github.com/marmotedu/api/apiserver/v1"
	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"
)

// UserSrv 定义处理用户请求的函数
type UserSrv interface {
	Create(ctx context.Context, user *v1.User, opts metav1.CreateOptions) error
	Update(ctx context.Context, user *v1.User, opts metav1.UpdateOptions) error
	Delete(ctx context.Context, username string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, usernames []string, opts metav1.DeleteOptions) error
	Get(ctx context.Context, username string, opts metav1.GetOptions) (*v1.User, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.UserList, error)
	ListWithBadPerformance(ctx context.Context, opts metav1.ListOptions) (*v1.UserList, error)
	ChangePassword(ctx context.Context, user *v1.User) error
}

var _ UserSrv = &userService{}

type userService struct {
	store store.Factory
}

func newUser(srv *service) *userService {
	return &userService{srv.store}
}

func (u *userService) Create(ctx context.Context, user *v1.User, opts metav1.CreateOptions) error {
	// TODO implement me
	panic("implement me")
}

func (u *userService) Update(ctx context.Context, user *v1.User, opts metav1.UpdateOptions) error {
	// TODO implement me
	panic("implement me")
}

func (u *userService) Delete(ctx context.Context, username string, opts metav1.DeleteOptions) error {
	// TODO implement me
	panic("implement me")
}

func (u *userService) DeleteCollection(ctx context.Context, usernames []string, opts metav1.DeleteOptions) error {
	// TODO implement me
	panic("implement me")
}

func (u *userService) Get(ctx context.Context, username string, opts metav1.GetOptions) (*v1.User, error) {
	// TODO implement me
	panic("implement me")
}

func (u *userService) List(ctx context.Context, opts metav1.ListOptions) (*v1.UserList, error) {
	// TODO implement me
	panic("implement me")
}

func (u *userService) ListWithBadPerformance(ctx context.Context, opts metav1.ListOptions) (*v1.UserList, error) {
	// TODO implement me
	panic("implement me")
}

func (u *userService) ChangePassword(ctx context.Context, user *v1.User) error {
	// TODO implement me
	panic("implement me")
}
