package meta

import (
	"context"
	"fmt"
	"time"

	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/internel/generated_code/v5.0.0/proto"
	admin "github.com/vesoft-inc/nebula-ng-tools/golang/pkg/internel/generated_code/v5.0.0/proto/admin"
)

type (
	ChangePasswordReq struct {
		user            string
		currentPassword string
		newPassword     string
	}

	CreateUserReq struct {
		user     string
		authType string
		authInfo string
	}

	AlterUserReq struct {
		user     string
		authInfo string
		active   bool
	}

	DropUserReq struct {
		user string
	}

	ListUsersReq struct {
		user []string
	}
	UserInfo struct {
		Name            string
		Active          bool
		AuthType        string
		CreatedTime     *time.Time
		LastLoginTime   *time.Time
		LastUpdatedTime *time.Time
		DisabledTime    *time.Time
		AuthInfo        string
	}

	ListUsersResp struct {
		Users []*UserInfo
	}
)

func NewChangePasswordReq(user, currentPassword, newPassword string) *ChangePasswordReq {
	return &ChangePasswordReq{
		user:            user,
		currentPassword: currentPassword,
		newPassword:     newPassword,
	}
}

func (c *metaClient) ChangePassword(req *ChangePasswordReq) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	in := &admin.ChangePasswordRequest{
		Username:    []byte(req.user),
		OldPassword: []byte(req.currentPassword),
		NewPassword: []byte(req.newPassword),
	}
	resp, err := c.retry(func() (responseHeader, error) {
		return c.client.ChangePassword(ctx, in)
	})
	if err != nil {
		return err
	}
	response, ok := resp.(*admin.ChangePasswordResponse)
	if !ok {
		return fmt.Errorf("invalid response")
	}
	responseHeader, err := getResponseHeader(response)
	if err != nil {
		return err
	}
	if !responseHeader.IsSucceeded() {
		return responseHeader.GetStatus()
	}
	return nil
}

func NewCreateUserReq(user string, authType string, authInfo string) *CreateUserReq {
	return &CreateUserReq{
		user:     user,
		authType: authType,
		authInfo: authInfo,
	}
}

func (c *metaClient) CreateUser(req *CreateUserReq) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	in := &admin.CreateUserRequest{
		Header:   &admin.AdminRequestHeader{Token: c.token},
		Username: []byte(req.user),
		AuthType: []byte(req.authType),
		AuthInfo: []byte(req.authInfo),
	}
	resp, err := c.retry(func() (responseHeader, error) {
		return c.client.CreateUser(ctx, in)
	})
	if err != nil {
		return err
	}
	return responseIsErr(resp)
}

func NewDropUserReq(user string) *DropUserReq {
	return &DropUserReq{
		user: user,
	}
}

func (c *metaClient) DropUser(req *DropUserReq) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	in := &admin.DropUserRequest{
		Header:   &admin.AdminRequestHeader{Token: c.token},
		Username: []byte(req.user),
	}
	resp, err := c.retry(func() (responseHeader, error) {
		return c.client.DropUser(ctx, in)
	})
	if err != nil {
		return err
	}
	return responseIsErr(resp)
}

func NewAlterUserReq(user string, authInfo string, active bool) *AlterUserReq {
	return &AlterUserReq{
		user:     user,
		authInfo: authInfo,
		active:   active,
	}
}

func (c *metaClient) AlterUser(req *AlterUserReq) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	in := &admin.AlterUserRequest{
		Header:   &admin.AdminRequestHeader{Token: c.token},
		Username: []byte(req.user),
		AuthInfo: []byte(req.authInfo),
		Active:   req.active,
	}
	resp, err := c.retry(func() (responseHeader, error) {
		return c.client.AlterUser(ctx, in)
	})
	if err != nil {
		return err
	}
	return responseIsErr(resp)
}

func NewListUsersReq(user []string) *ListUsersReq {
	return &ListUsersReq{
		user: user,
	}
}

func (c *metaClient) ListUsers(req *ListUsersReq) (*ListUsersResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	users := make([][]byte, 0)
	for _, u := range req.user {
		users = append(users, []byte(u))
	}
	in := &admin.ListUserRequest{
		Header:    &admin.AdminRequestHeader{Token: c.token},
		Usernames: users,
	}
	resp, err := c.retry(func() (responseHeader, error) {
		return c.client.ListUser(ctx, in)
	})
	if err != nil {
		return nil, err
	}
	if err := responseIsErr(resp); err != nil {
		return nil, err
	}
	response, ok := resp.(*admin.ListUserResponse)
	if !ok {
		return nil, fmt.Errorf("invalid response")
	}
	usersInfo := make([]*UserInfo, 0)
	for _, u := range response.Users {
		userInfo := &UserInfo{
			Name:            string(u.Username),
			Active:          u.Active,
			AuthType:        string(u.AuthType),
			AuthInfo:        string(u.AuthInfo),
			CreatedTime:     proto.ConvertZonedTime(u.CreatedTime),
			LastLoginTime:   proto.ConvertZonedTime(u.LastLoginTime),
			LastUpdatedTime: proto.ConvertZonedTime(u.LastUpdatedTime),
			DisabledTime:    proto.ConvertZonedTime(u.DisabledTime),
		}
		usersInfo = append(usersInfo, userInfo)
	}
	return &ListUsersResp{Users: usersInfo}, nil
}
