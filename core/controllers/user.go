package controllers

import (
	"aegis/assessment-test/core/constant"
	"aegis/assessment-test/core/entity"
	"aegis/assessment-test/core/repository"
	"aegis/assessment-test/core/repository/models"
	hash "aegis/assessment-test/utils/encrypt"
	"context"
	"net/http"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	userRepo repository.UserRepository
}

func NewUserController(
	userRepo repository.UserRepository,
) *UserController {
	return &UserController{
		userRepo: userRepo,
	}
}

func (u *UserController) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		registerReq := entity.RegisterRequest{}
		c.Bind(&registerReq)
		err := c.Validate(&registerReq)
		if err != nil {
			logrus.Errorf("err validate request=%s", err.Error())
			return c.JSON(
				http.StatusBadRequest,
				constant.BadRequest(constant.CodeErrBadRequest, "there is some problem from input", err))
		}

		hashPassword, err := hash.HashPassword(registerReq.Password)
		if err != nil {
			logrus.Errorf("err hash password=%s", err.Error())
			return c.JSON(
				http.StatusInternalServerError,
				constant.InternalServerError(constant.CodeErrInternalServer, "failed to hash password", err))
		}

		err = u.userRepo.InsertNewUser(context.Background(), &models.UserSchema{
			Name:     registerReq.Name,
			Username: registerReq.Username,
			Password: hashPassword,
			Role:     registerReq.Role,
		})
		if err != nil {
			logrus.Errorf("err insert user=%s", err.Error())
			return c.JSON(
				http.StatusInternalServerError,
				constant.InternalServerError(constant.CodeErrInternalServer, "failed to register new user", err))
		}

		resp := entity.RegisterResponse{
			Name:     registerReq.Name,
			Username: registerReq.Username,
			Role:     registerReq.Role,
		}

		return c.JSON(
			http.StatusOK,
			constant.Success(constant.CodeSuccess, "success register new user", resp))
	}
}
