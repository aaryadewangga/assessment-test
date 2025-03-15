package controllers

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
)

func (u *UserController) isExistingUsername(ctx context.Context, username string) (bool, error) {
	user, err := u.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		logrus.Errorf("failed get user=%s", err.Error())
		return false, nil
	}

	if user != nil {
		logrus.Errorf("username is same with existing user username=%s existingUsername=%s", username, user.Username)
		return true, fmt.Errorf("existing username detected, please use another username")
	}

	return false, nil
}
