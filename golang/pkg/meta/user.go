package meta

import (
	"context"
	"fmt"

	admin "github.com/vesoft-inc/nebula-ng-tools/golang/pkg/internel/generated_code/v5.0.0/proto/admin"
)

type (
	ChangePasswordReq struct {
		user            string
		currentPassword string
		newPassword     string
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
		return responseHeader.GetError()
	}
	return nil
}
