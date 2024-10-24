package common

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
)

func ResetPassword(c meta.Client, user, old, new string) error {
	req := meta.NewChangePasswordReq(
		user,
		old,
		new,
	)
	if err := c.ChangePassword(req); err != nil {
		return err
	}
	return nil
}

func GetPromptPassword() (string, string, error) {
	currentPassword := promptui.Prompt{
		Label:     "Current password:",
		AllowEdit: true,
		Mask:      rune(' '),
	}
	currentPasswordStr, err := currentPassword.Run()
	if err != nil {
		return "", "", err
	}
	newPassword := promptui.Prompt{
		Label:     "New password:",
		AllowEdit: true,
		Mask:      rune(' '),
	}
	newPasswordStr, err := newPassword.Run()
	if err != nil {
		return "", "", err
	}
	confirmPassword := promptui.Prompt{
		Label:     "Retype new password:",
		AllowEdit: true,
		Mask:      rune(' '),
	}
	confirmPasswordStr, err := confirmPassword.Run()
	if err != nil {
		return "", "", err
	}
	if newPasswordStr != confirmPasswordStr {
		return "", "", fmt.Errorf("Sorry, the passwords you entered do not match.")
	}

	return currentPasswordStr, newPasswordStr, nil
}
