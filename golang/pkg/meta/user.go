package meta

import (
	"context"
	"fmt"
	"time"

	nebula "github.com/vesoft-inc/nebula-ng-tools/golang"
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
		Name            string     `json:"name"`
		Active          bool       `json:"active"`
		AuthType        string     `json:"auth_type"`
		CreatedTime     *time.Time `json:"created_time"`
		LastLoginTime   *time.Time `json:"last_login_time"`
		LastUpdatedTime *time.Time `json:"last_updated_time"`
		DisabledTime    *time.Time `json:"disabled_time"`
		AuthInfo        string     `json:"auth_info"`
	}

	ListUsersResp struct {
		Users []*UserInfo `json:"users"`
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
	ctx, cancel := context.WithTimeout(context.Background(), c.requestTimeout)
	defer cancel()
	in := &admin.ChangePasswordRequest{
		Username:    []byte(req.user),
		OldPassword: []byte(req.currentPassword),
		NewPassword: []byte(req.newPassword),
	}
	resp, err := c.client.ChangePassword(ctx, in)
	if err != nil {
		return err
	}
	// retry for change password
	if nebula.ErrorCode(resp.Header.GetStatus().GetCode()) != nebula.ERROR_LEADER_CHANGED {
		return responseIsErr(resp)
	}
	leader := resp.Header.GetLeader()
	if leader == nil {
		return fmt.Errorf("invalid leader")
	}
	c.Close()
	c.address = fmt.Sprintf("%s:%d", leader.GetHost(), leader.GetPort())
	err = c.open(string(leader.GetHost()), int(leader.GetPort()), c.connectTimeout, nil)
	if err != nil {
		return err
	}
	resp, err = c.client.ChangePassword(ctx, in)
	if err != nil {
		return err
	}
	return responseIsErr(resp)
}

func NewCreateUserReq(user string, authType string, authInfo string) *CreateUserReq {
	return &CreateUserReq{
		user:     user,
		authType: authType,
		authInfo: authInfo,
	}
}

func (c *metaClient) CreateUser(req *CreateUserReq) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.requestTimeout)
	defer cancel()
	in := &admin.CreateUserRequest{
		Header:   &admin.RequestHeader{Token: c.token},
		Username: []byte(req.user),
		AuthType: []byte(req.authType),
		AuthInfo: []byte(req.authInfo),
	}
	resp, err := c.execute(func() (responseHeader, error) {
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
	ctx, cancel := context.WithTimeout(context.Background(), c.requestTimeout)
	defer cancel()
	in := &admin.DropUserRequest{
		Header:   &admin.RequestHeader{Token: c.token},
		Username: []byte(req.user),
	}
	resp, err := c.execute(func() (responseHeader, error) {
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
	ctx, cancel := context.WithTimeout(context.Background(), c.requestTimeout)
	defer cancel()
	in := &admin.AlterUserRequest{
		Header:   &admin.RequestHeader{Token: c.token},
		Username: []byte(req.user),
		AuthInfo: []byte(req.authInfo),
		Active:   req.active,
	}
	resp, err := c.execute(func() (responseHeader, error) {
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
	ctx, cancel := context.WithTimeout(context.Background(), c.requestTimeout)
	defer cancel()
	users := make([][]byte, 0)
	for _, u := range req.user {
		users = append(users, []byte(u))
	}
	in := &admin.ListUserRequest{
		Header:    &admin.RequestHeader{Token: c.token},
		Usernames: users,
	}
	resp, err := c.execute(func() (responseHeader, error) {
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
